/* eslint-disable @typescript-eslint/camelcase */
import {Injectable} from '@nestjs/common';
import {Connection, StartSession} from 'mongoose';
import {InjectConnection} from '@nestjs/mongoose';
import {IbcChainConfigSchema} from '../schema/ibc_chain_config.schema';
import {IbcDenomSchema} from '../schema/ibc_denom.schema';
import {IbcTxSchema} from '../schema/ibc_tx.schema';
import {TxSchema} from '../schema/tx.schema';
import {IbcBlockSchema} from '../schema/ibc_block.schema';
import {IbcTaskRecordSchema} from '../schema/ibc_task_record.schema';
import {IbcTxType} from '../types/schemaTypes/ibc_tx.interface';
import {JSONparse} from '../util/util';
import {getDcDenom} from '../helper/denom.helper';
import {SubState, TaskTime, IbcTxTable} from '../constant';

import {
    TaskEnum,
    TxType,
    TxStatus,
    IbcTxStatus,
    RecordLimit,
    IbcTaskRecordStatus,
} from '../constant';
import {dateNow} from "../helper/date.helper";
import {getTaskStatus} from "../helper/task.helper";
import {SyncTaskSchema} from "../schema/sync.task.schema";
import {Logger} from "../logger";

@Injectable()
export class IbcTxHandler {
    private ibcTaskRecordModel;
    private chainConfigModel;
    private ibcTxLatestModel;
    private ibcTxModel;
    private ibcDenomModel;

    constructor(@InjectConnection() private readonly connection: Connection) {
        this.getModels();
    }

    // getModels
    async getModels(): Promise<void> {
        // ibcTaskRecordModel
        this.ibcTaskRecordModel = await this.connection.model(
            'ibcTaskRecordModel',
            IbcTaskRecordSchema,
            'ibc_task_record',
        );

        // chainConfigModel
        this.chainConfigModel = await this.connection.model(
            'chainConfigModel',
            IbcChainConfigSchema,
            'chain_config',
        );


        // ibcTxModel
        this.ibcTxModel = await this.connection.model(
            'ibcTxModel',
            IbcTxSchema,
            IbcTxTable.IbcTxTableName,
        );

        this.ibcTxLatestModel = await this.connection.model(
            'ibcTxLatestModel',
            IbcTxSchema,
            IbcTxTable.IbcTxLatestTableName,
        );

        // ibcDenomModel
        this.ibcDenomModel = await this.connection.model(
            'ibcDenomModel',
            IbcDenomSchema,
            'ibc_denom',
        );

    }

    getIbcTxModel() {
        return this.ibcTxModel
    }

    getIbcTxLatestModel() {
        return this.ibcTxLatestModel
    }

    // ibcTx first（transfer）
    async getAllChainsMap() {
        const allChains = await this.chainConfigModel.findAll();
        const allChainsMap = new Map, allChainsDenomPathsMap = new Map
        if (allChains?.length) {
            allChains.forEach(item => {
                if (item?.chain_id) {
                    allChainsMap.set(item.chain_id, item)
                    if (item?.ibc_info?.length) {
                        item.ibc_info.forEach(info => {
                            if (info?.paths?.length) {
                                info.paths.forEach(path => {
                                    if (path?.channel_id && path?.port_id)
                                        allChainsDenomPathsMap.set(`${info.chain_id}${path.counterparty.channel_id}${path.counterparty.port_id}`, `${path.counterparty.channel_id}${path.counterparty.port_id}`)
                                })
                            }
                        })
                    }
                }
            })
        }
        return {
            allChainsMap,
            allChainsDenomPathsMap
        }
    }

    async getDenomRecordByChainId(chainId) {
        const ibcDenomRecordMap = new Map
        const ibcDenomRecord = await this.ibcDenomModel.findAllDenomRecord(chainId);
        if (ibcDenomRecord?.length) {
            ibcDenomRecord.forEach(ibcDenomRecordItem => {
                if (ibcDenomRecordItem?.denom) {
                    ibcDenomRecordMap.set(ibcDenomRecordItem?.denom, ibcDenomRecordItem)
                }
            })
        }
        return ibcDenomRecordMap
    }

    async getRecordLimitTx(chainId, height, limit): Promise<Array<any>> {
        const txModel = await this.connection.model(
            'txModel',
            TxSchema,
            `sync_${chainId}_tx`,
        );
        let txs = [];
        //根据块高排序 查询最后限制条数的交易
        const txsByLimit = await txModel.queryTxListSortHeight({
            type: TxType.transfer,
            height: height,
            limit: limit,
        });
        // 根据块高查询限制条数的最后一条交易
        const txsByHeight = txsByLimit.length
            ? await txModel.queryTxListByHeight(
                TxType.transfer,
                txsByLimit[txsByLimit.length - 1].height,
            )
            : [];
        //去重
        const hash = {};
        txs = [...txsByLimit, ...txsByHeight].reduce((txsResult, next) => {
            hash[next.tx_hash]
                ? ''
                : (hash[next.tx_hash] = true) && txsResult.push(next);
            return txsResult;
        }, []);
        return txs
    }

    async checkTaskFollowingStatus(chainId): Promise<boolean> {
        const taskModel = await this.connection.model(
            'txModel',
            SyncTaskSchema,
            `sync_${chainId}_task`,
        );
        return await getTaskStatus(chainId, taskModel, TaskEnum.tx)
    }

    async parseIbcTx(ibcTxModel,dateNow): Promise<void> {
        const allChains = await this.chainConfigModel.findAll();
        const {allChainsMap, allChainsDenomPathsMap} = await this.getAllChainsMap()
        let ibcDenoms = []
        for (const {chain_id} of allChains) {
            // get taskRecord by chain_id
            let taskRecord = await this.ibcTaskRecordModel.findTaskRecord(chain_id);
            if (!taskRecord) {
                await this.ibcTaskRecordModel.insertManyTaskRecord({
                    task_name: `sync_${chain_id}_transfer`,
                    status: IbcTaskRecordStatus.OPEN,
                    height: 0,
                    create_at: dateNow,
                    update_at: dateNow,
                });
                taskRecord = await this.ibcTaskRecordModel.findTaskRecord(chain_id);
            } else {
                if (taskRecord.status === IbcTaskRecordStatus.CLOSE) continue;
            }
            const taskCount = await this.checkTaskFollowingStatus(chain_id)
            if (!taskCount) continue

            const denomMap =  await this.getDenomRecordByChainId(chain_id)

            const txs = await this.getRecordLimitTx(chain_id, taskRecord.height, RecordLimit)
            let {
                handledTx,
                denoms
            } = await this.handlerSourcesTx(txs, chain_id, dateNow, allChainsMap, allChainsDenomPathsMap,denomMap)
            ibcDenoms = [...denoms]
            if (ibcDenoms?.length) {
                await this.ibcDenomModel.insertManyDenom(ibcDenoms);
            }
            if (handledTx?.length) {
                const session = await this.connection.startSession()
                session.startTransaction()
                try {

                    await ibcTxModel.insertManyIbcTx(handledTx)
                    taskRecord.height = handledTx[handledTx.length - 1]?.sc_tx_info?.height;
                    taskRecord.update_at = dateNow
                    await this.ibcTaskRecordModel.updateTaskRecord(taskRecord);
                    await session.commitTransaction();
                    session.endSession();
                } catch (e) {
                    Logger.log(e, 'transaction is error')
                    await session.abortTransaction()
                    session.endSession();
                }
            }

        }
        Logger.debug(`end parseIbcTx time ${dateNow}`)
    }

    async getProcessingTxs(ibcTxModel,substate) {
        if (substate?.length == 1 && substate[0] === 0) {
            const ibcTxs = await ibcTxModel.queryTxList({
                status: IbcTxStatus.PROCESSING,
                substate: substate,
                limit: RecordLimit,
            });
            return ibcTxs
        } else {
            const substateTxs = await ibcTxModel.queryTxListBySubstate({
                status: IbcTxStatus.PROCESSING,
                substate: substate,
                limit: RecordLimit,
            });
            return substateTxs
        }

    }

    async getPacketIds(txs) {
        const packetIds = []
        if (txs?.length) {
            for (const tx of txs) {
                if (tx?.sc_tx_info?.msg?.msg?.packet_id) {
                    // ibcTx.sc_tx_info.msg.msg.timeout_height.revision_height;
                    if (tx?.dc_chain_id
                        && tx?.sc_tx_info?.msg?.msg?.packet_id
                        && tx?.sc_tx_info?.msg?.msg?.timeout_height?.revision_height >= 0
                        // && tx?.sc_tx_info?.msg?.msg?.timeout_timestamp >= 0
                    ) {
                        packetIds.push(
                            {
                                chainId: tx.dc_chain_id,
                                height: tx.sc_tx_info.msg.msg.timeout_height.revision_height,
                                packetId: tx.sc_tx_info.msg.msg.packet_id,
                                timeOutTime: tx.sc_tx_info.msg.msg.timeout_timestamp
                            }
                        )
                    }
                }

            }
        }
        return packetIds
    }

    async changeIbcTxState(ibcTxModel,dateNow, substate: number[]): Promise<void> {
        const ibcTxs = await this.getProcessingTxs(ibcTxModel,substate)
        // const ibcTxs = await ibcTxModel.queryTxByRecordId("transferchannel-28transferchannel-12163541irishub_qaF29DEFE6DE7C6355489F5D2CCC96EF2D630E351F843852050D1A29317C21FBDB0")


        let packetIdArr = ibcTxs?.length ? await this.getPacketIds(ibcTxs) : [];
        let recvPacketTxMap = new Map, chainHeightMap = new Map, refundedTxTxMap = new Map,
            acknowledgeTxsMap = new Map, needUpdateTxs = [],
            denoms = []

        let dcChains = [],scChains = [];
        ibcTxs.forEach(item => {
            if (item?.dc_chain_id ) {
                dcChains.push(item.dc_chain_id)
            }
            if(item?.sc_chain_id){
                scChains.push(item.sc_chain_id)
            }
        })

        const allChains = [...dcChains,...scChains]
        const currentDcChains = Array.from(new Set(allChains))
        // 根据链的配置信息，查询出每条链的 recv_packet 成功的交易的那条记录
        if (currentDcChains?.length) {

            for (const chain of currentDcChains) {
                if (chain) {
                    const blockModel = await this.connection.model(
                        'blockModel',
                        IbcBlockSchema,
                        `sync_${chain}_block`,
                    );
                    const chainHeightObj = await blockModel.findLatestBlock();
                    if (chainHeightObj && JSON.stringify(chainHeightObj) !== '{}') {
                        let {height, time} = await blockModel.findLatestBlock();
                        chainHeightMap.set(chain, {height, time})
                    }
                }
            }
            for (const chain of currentDcChains) {
                if (chain) {
                    const txModel = await this.connection.model(
                        'txModel',
                        TxSchema,
                        `sync_${chain}_tx`,
                    );


                    const taskCount = await this.checkTaskFollowingStatus(chain)
                    if (!taskCount) continue
                    //每条链最新的高度
                    let refundedTxPacketIdsMap = new Map
                    const refundedTxPacketIds = [],acknowledgeTxPacketIds = [];
                    packetIdArr.forEach(item => {
                        if (item?.chainId && item?.height || item?.timeOutTime) {
                            const currentChainLatestObj = chainHeightMap.get(item.chainId)
                            const dateLen = String(item?.timeOutTime).length
                            if (dateLen > 10) {
                                refundedTxPacketIdsMap.set(item.packetId,'')
                                refundedTxPacketIds.push(item.packetId)
                            }else if (item.height < currentChainLatestObj?.height || item.timeOutTime < currentChainLatestObj?.time) {
                                if (item?.packetId) {
                                    refundedTxPacketIdsMap.set(item.packetId, '')
                                    refundedTxPacketIds.push(item.packetId)
                                }
                            }
                        }
                    })
                    const recvPacketIds = packetIdArr.map(item => {
                        if (item?.packetId) {
                            return item.packetId
                        }
                    })


                    // txs  数组
                    if (recvPacketIds?.length) {
                        const txs = await txModel.queryTxsByPacketId({
                            type: TxType.recv_packet,
                            limit: packetIdArr.length,
                            // status: TxStatus.SUCCESS,
                            packet_id: recvPacketIds,
                        });

                        if (txs?.length) {
                            for (let tx of txs) {
                                if (tx?.msgs?.length) {
                                    for (let msg of tx.msgs) {
                                        if (msg?.type === TxType.recv_packet && msg.msg.packet_id) {
                                            if (tx?.status) {
                                                // recv_packet tx is success
                                                recvPacketTxMap.set(`${chain}${msg.msg.packet_id}`, tx)
                                            }else{
                                                // recv_packet tx is fail
                                                acknowledgeTxPacketIds.push(msg.msg.packet_id)
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                    if (refundedTxPacketIds?.length) {
                        // const notIncludeNullPacketIds = refundedTxPacketIds.filter( item =>  item)
                        const refundedTxs = await txModel.queryTxListByPacketId({
                            type: TxType.timeout_packet,
                            limit: refundedTxPacketIds.length,
                            status: TxStatus.SUCCESS,
                            packet_id: refundedTxPacketIds,
                        });
                        if (refundedTxs?.length) {
                            for (let refundedTx of refundedTxs) {
                                if (refundedTx?.msgs?.length) {
                                    for (let msg of refundedTx.msgs) {
                                        if (msg?.type === TxType.timeout_packet && msg.msg.packet_id) {
                                            refundedTxTxMap.set(`${chain}${msg.msg.packet_id}`, refundedTx)
                                        }
                                    }
                                }
                            }
                        }
                    }

                    if (acknowledgeTxPacketIds?.length) {
                        const acknowledgeTxs = await txModel.queryTxListByPacketId({
                            type: TxType.acknowledge_packet,
                            limit: acknowledgeTxPacketIds.length,
                            status: TxStatus.SUCCESS,
                            packet_id: acknowledgeTxPacketIds,
                        });

                        if (acknowledgeTxs?.length) {
                            for (let acknowledgeTx of acknowledgeTxs) {
                                if (acknowledgeTx?.msgs?.length) {
                                    for (let msg of acknowledgeTx.msgs) {
                                        if (msg?.type === TxType.acknowledge_packet && msg.msg.packet_id) {
                                            acknowledgeTxsMap.set(`${chain}${msg.msg.packet_id}`,acknowledgeTx)
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
        for (let [index, ibcTx] of ibcTxs.entries()) {
            if (!ibcTx.dc_chain_id) continue
            if (!recvPacketTxMap.size && !refundedTxTxMap.size && !acknowledgeTxsMap.size) {
                ibcTx.substate = SubState.SuccessRecvPacketNotFound;
                ibcTx = this.setNextTryTime(ibcTx, index)
                needUpdateTxs.push(ibcTx)
            } else if (recvPacketTxMap?.has(`${ibcTx.dc_chain_id}${ibcTx.sc_tx_info.msg.msg.packet_id}`)) {
                const recvPacketTx = recvPacketTxMap?.get(`${ibcTx.dc_chain_id}${ibcTx.sc_tx_info.msg.msg.packet_id}`)
                // let counter_party_tx = null
                recvPacketTx && recvPacketTx.msgs.length && recvPacketTx.msgs.forEach(async (msg, msgIndex) => {
                    if (msg.type === TxType.recv_packet && ibcTx.sc_tx_info.msg.msg.packet_id === msg.msg.packet_id) {
                        const {dc_denom, dc_denom_path} = getDcDenom(msg);

                        // add write_acknowledgement solution， value type is string;
                        let result = '';
                        const tx_events = recvPacketTx.events_new.find(
                            event_new => {
                                return event_new.msg_index === msgIndex;
                            },
                        );
                        tx_events &&
                        tx_events.events.forEach(event => {
                            if (event.type === 'write_acknowledgement') {
                                event.attributes.forEach(attribute => {
                                    if (attribute.key === 'packet_ack') {
                                        result = attribute.value === '' ? 'false' : attribute.value && attribute.value.includes('error') ? 'false' : 'true';
                                    }
                                });
                            }
                        });
                        switch (result) {
                            case "true":
                                ibcTx.status = IbcTxStatus.SUCCESS
                                ibcTx.substate = 0
                                ibcTx.retry_times = 0
                                ibcTx.next_try_time = 0
                                break;
                            case "false":
                                ibcTx.substate = SubState.RecvPacketAckFailed;
                                ibcTx = this.setNextTryTime(ibcTx, index)
                                break;
                        }
                        ibcTx.status =
                            result === 'false' ? IbcTxStatus.FAILED : IbcTxStatus.SUCCESS;
                        ibcTx.dc_tx_info = {
                            hash: recvPacketTx.tx_hash,
                            status: recvPacketTx.status,
                            time: recvPacketTx.time,
                            height: recvPacketTx.height,
                            fee: recvPacketTx.fee,
                            msg_amount: msg.msg.token,
                            msg,
                        };
                        ibcTx.update_at = dateNow;
                        // ibcTx.tx_time = counter_party_tx.time;
                        ibcTx.denoms['dc_denom'] = dc_denom;
                        const denom_path =
                            dc_denom_path === ibcTx.base_denom
                                ? ''
                                : dc_denom_path.replace(`/${ibcTx.base_denom}`, '');
                        needUpdateTxs.push(ibcTx)
                        if (ibcTx.status === IbcTxStatus.SUCCESS) {
                            const denomObj = {
                                chain_id: ibcTx.dc_chain_id,
                                denom: dc_denom,
                                base_denom: ibcTx.base_denom,
                                denom_path: denom_path,
                                is_source_chain: !Boolean(denom_path),
                                is_base_denom: Boolean(dc_denom === ibcTx.base_denom),
                                create_at: dateNow,
                                update_at: ''
                            }
                            denoms.push(denomObj)
                        }

                    }
                })

            } else if (acknowledgeTxsMap?.has(`${ibcTx.sc_chain_id}${ibcTx.sc_tx_info.msg.msg.packet_id}`)){
                const acknowledgeTx = acknowledgeTxsMap?.get(`${ibcTx.sc_chain_id}${ibcTx.sc_tx_info.msg.msg.packet_id}`);
                acknowledgeTx && acknowledgeTx.msgs.forEach(msg => {
                    if (
                        msg.type === TxType.acknowledge_packet &&
                        ibcTx.sc_tx_info.msg.msg.packet_id === msg.msg.packet_id
                    ){
                        ibcTx.status = IbcTxStatus.REFUNDED;
                        ibcTx.retry_times = 0;
                        ibcTx.next_try_time = 0
                        ibcTx.refunded_tx_info = {
                            hash: acknowledgeTx.tx_hash,
                            status: acknowledgeTx.status,
                            time: acknowledgeTx.time,
                            height: acknowledgeTx.height,
                            fee: acknowledgeTx.fee,
                            msg_amount: msg.msg.token,
                            msg,
                        };

                        ibcTx.update_at = dateNow;
                        ibcTx.substate = 0
                        needUpdateTxs.push(ibcTx)
                    }
                        })
            }else {
                /*
                * 没有找到的结果
                * */
                if (refundedTxTxMap.has(`${ibcTx.sc_chain_id}${ibcTx.sc_tx_info.msg.msg.packet_id}`)) {
                    const refunded_tx = refundedTxTxMap?.get(`${ibcTx.sc_chain_id}${ibcTx.sc_tx_info.msg.msg.packet_id}`);
                    refunded_tx &&
                    refunded_tx.msgs.forEach(msg => {
                        if (
                            msg.type === TxType.timeout_packet &&
                            ibcTx.sc_tx_info.msg.msg.packet_id === msg.msg.packet_id
                        ) {
                            ibcTx.status = IbcTxStatus.REFUNDED;
                            ibcTx.retry_times = 0;
                            ibcTx.next_try_time = 0
                            ibcTx.refunded_tx_info = {
                                hash: refunded_tx.tx_hash,
                                status: refunded_tx.status,
                                time: refunded_tx.time,
                                height: refunded_tx.height,
                                fee: refunded_tx.fee,
                                msg_amount: msg.msg.token,
                                msg,
                            };

                            ibcTx.update_at = dateNow;
                            // ibcTx.tx_time = refunded_tx.time;
                            ibcTx.substate = 0
                            needUpdateTxs.push(ibcTx)
                        }
                    });
                } else {
                    ibcTx.substate = SubState.SuccessTimeoutPacketNotFound;
                    ibcTx = this.setNextTryTime(ibcTx, index)
                    needUpdateTxs.push(ibcTx)
                }
            }
        }
        if (needUpdateTxs?.length) {
            for (let needUpdateTx of needUpdateTxs) {
                await ibcTxModel.updateIbcTx(needUpdateTx);
            }
        }
        if (denoms?.length) {
            await this.ibcDenomModel.insertManyDenom(denoms);
        }
    }

    setNextTryTime(ibcTx, index) {
        ibcTx.retry_times = ibcTx.retry_times ? Number(ibcTx.retry_times) + 1 : 1
        const taskDiffTime = Math.floor(Number(ibcTx.retry_times) * TaskTime)
        ibcTx.next_try_time = Math.floor(Number(taskDiffTime) + Number(dateNow) + index)
        return ibcTx
    }

    // get dc_port、dc_channel、sequence
    getIbcInfoFromEventsMsg(
        tx,
        msgIndex,
    ): {
        dc_port: string;
        dc_channel: string;
        sequence: string;
        base_denom: string;
        denom_path: string;
    } {
        const msg = {
            dc_port: '',
            dc_channel: '',
            sequence: '',
            base_denom: '',
            denom_path: '',
        };

        tx.events_new[msgIndex] &&
        tx.events_new[msgIndex].events.forEach(evt => {
            if (evt.type === 'send_packet') {
                evt.attributes.forEach(attr => {
                    switch (attr.key) {
                        case 'packet_dst_port':
                            msg.dc_port = attr.value;
                            break;
                        case 'packet_dst_channel':
                            msg.dc_channel = attr.value;
                            break;
                        case 'packet_sequence':
                            msg.sequence = attr.value;
                            break;
                        case 'packet_data':
                            const packet_data = JSONparse(attr.value);
                            const denomOrigin = packet_data.denom;
                            const denomOriginSplit = denomOrigin.split('/');
                            msg.base_denom = denomOriginSplit[denomOriginSplit.length - 1];
                            msg.denom_path = denomOriginSplit
                                .slice(0, denomOriginSplit.length - 1)
                                .join('/');
                            break;
                        default:
                            break;
                    }
                });
            }
        });
        return msg;
    }

    async handlerSourcesTx(sourcesTx, chain_id, currentTime, allChainsMap, allChainsDenomPathsMap,denomMap) {
        let handledTx = [], denoms = []
        if (sourcesTx && chain_id) {
            sourcesTx.forEach((tx, txIndex) => {
                const height = tx.height;
                const log = tx.log;
                const time = tx.time;
                const hash = tx.tx_hash;
                const status = tx.status;
                const fee = tx.fee;
                tx.msgs.forEach(async (msg, msgIndex) => {
                    if (msg.type === TxType.transfer) {
                        const ibcTx: IbcTxType = {
                            record_id: '',
                            sc_addr: '',
                            dc_addr: '',
                            sc_port: '',
                            sc_channel: '',
                            sc_chain_id: '',
                            dc_port: '',
                            dc_channel: '',
                            dc_chain_id: '',
                            sequence: '',
                            status: 0,
                            sc_tx_info: {},
                            dc_tx_info: {},
                            refunded_tx_info: {},
                            log: {},
                            substate: 0,
                            denoms: {
                                sc_denom: '',
                                dc_denom: '',
                            },
                            base_denom: '',
                            create_at: '',
                            update_at: '',
                            tx_time: '',
                        };
                        switch (tx.status) {
                            case TxStatus.SUCCESS:
                                ibcTx.status = IbcTxStatus.PROCESSING;
                                break;
                            case TxStatus.FAILED:
                                ibcTx.status = IbcTxStatus.FAILED;
                                break;
                            default:
                                break;
                        }

                        const sc_chain_id = chain_id;
                        const sc_port = msg.msg.source_port;
                        const sc_channel = msg.msg.source_channel;
                        let dc_chain_id = '';
                        let dc_channel = '';
                        let dc_port = '';
                        const sc_addr = msg.msg.sender;
                        const dc_addr = msg.msg.receiver;
                        const sc_denom = msg.msg.token.denom;
                        const msg_amount = msg.msg.token;
                        let base_denom = '';
                        let denom_path = '';
                        let sequence = '';
                        //根据 sc_chain_id  sc_port sc_channel 查询目标链的信息
                        let dcChainConfig: any = {}
                        if (sc_chain_id && allChainsMap) {
                            if (allChainsMap.has(sc_chain_id)) {
                                const currentChainInfo = allChainsMap.get(sc_chain_id)
                                if (currentChainInfo?.ibc_info?.length) {
                                    currentChainInfo?.ibc_info.forEach(item => {
                                        if (item.paths?.length) {
                                            item.paths.forEach(pathItem => {
                                                if (pathItem?.channel_id === sc_channel && pathItem?.port_id === sc_port) {
                                                    dcChainConfig = currentChainInfo
                                                    dc_chain_id = item.chain_id;
                                                    dc_channel = pathItem.counterparty.channel_id;
                                                    dc_port = pathItem.counterparty.port_id;
                                                }
                                            })
                                        }
                                    })
                                }
                            }
                        }
                        if (ibcTx.status === IbcTxStatus.FAILED) {
                            // get base_denom、denom_path from ibc_denom collection
                            if (sc_denom.indexOf('ibc') !== -1) {
                                if (denomMap?.has(sc_denom)) {
                                    const ibcDenom = denomMap.get(sc_denom)
                                    if (ibcDenom?.base_denom) {
                                        base_denom = ibcDenom.base_denom
                                        denom_path = ibcDenom?.denom_path
                                    }
                                }
                            } else {
                                base_denom = sc_denom;
                            }
                        } else {
                            // get base_denom、denom_path、sequence from events
                            const event_msg = this.getIbcInfoFromEventsMsg(tx, msgIndex);
                            dc_port = event_msg.dc_port;
                            dc_channel = event_msg.dc_channel;
                            base_denom = event_msg.base_denom;
                            denom_path = event_msg.denom_path;
                            sequence = event_msg.sequence;
                        }

                        if (!dc_chain_id && ibcTx.status !== IbcTxStatus.FAILED) {
                            ibcTx.status = IbcTxStatus.SETTING;
                        }

                        ibcTx.record_id = `${sc_port}${sc_channel}${dc_port}${dc_channel}${sequence}${sc_chain_id}${hash}${msgIndex}`;
                        ibcTx.sc_addr = sc_addr;
                        ibcTx.dc_addr = dc_addr;
                        ibcTx.sc_port = sc_port;
                        ibcTx.sc_channel = sc_channel;
                        ibcTx.sc_chain_id = sc_chain_id;
                        ibcTx.dc_port = dc_port;
                        ibcTx.dc_channel = dc_channel;
                        ibcTx.dc_chain_id = dc_chain_id;
                        ibcTx.sequence = sequence;
                        ibcTx.denoms['sc_denom'] = sc_denom;
                        ibcTx.base_denom = base_denom;
                        ibcTx.create_at = currentTime;
                        ibcTx.update_at = currentTime;
                        ibcTx.tx_time = tx.time;
                        ibcTx.sc_tx_info = {
                            hash,
                            status,
                            time,
                            height,
                            fee,
                            msg_amount,
                            msg,
                        };
                        ibcTx.log['sc_log'] = log;
                        let isBaseDenom = null
                        if (sc_denom === base_denom) {
                            isBaseDenom = true
                        } else {
                            if (denom_path && denom_path.split('/').length > 1) {
                                const dc_port = denom_path.split('/')[0];
                                const dc_channel = denom_path.split('/')[1];
                                if (allChainsDenomPathsMap.has(`${dc_chain_id}${dc_channel}${dc_port}`)) {
                                    isBaseDenom = false
                                } else {
                                    isBaseDenom = true
                                }
                            }
                        }


                        if (ibcTx.status === IbcTxStatus.PROCESSING) {
                            denoms.push({
                                chain_id: sc_chain_id,
                                denom: sc_denom,
                                base_denom: base_denom,
                                denom_path: denom_path,
                                is_source_chain: !Boolean(denom_path),
                                is_base_denom: isBaseDenom,
                                create_at: dateNow,
                                update_at: ''
                            })
                        }
                        handledTx.push(ibcTx)
                    }
                });
            });
            return {
                handledTx,
                denoms
            }
        }
    }
}
