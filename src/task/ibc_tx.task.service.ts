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
import {ChainHttp} from '../http/lcd/chain.http';
import {IbcTxType} from '../types/schemaTypes/ibc_tx.interface';
import {JSONparse} from '../util/util';
import {getDcDenom} from '../helper/denom.helper';

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

@Injectable()
export class IbcTxTaskService {
    private ibcTaskRecordModel;
    private chainConfigModel;
    private ibcTxModel;
    private ibcDenomModel;

    constructor(@InjectConnection() private readonly connection: Connection) {
        this.getModels();
        this.doTask = this.doTask.bind(this);
    }

    async doTask(taskName?: TaskEnum): Promise<void> {
        const dateNow = Math.floor(new Date().getTime() / 1000);
        this.parseIbcTx(dateNow);
        this.changeIbcTxState(dateNow);
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
            'ex_ibc_tx',
        );

        // ibcDenomModel
        this.ibcDenomModel = await this.connection.model(
            'ibcDenomModel',
            IbcDenomSchema,
            'ibc_denom',
        );

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
                                        allChainsDenomPathsMap.set(info.chain_id, `${path.channel_id}${path.port_id}`)
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

    async getDenomRecord() {
        const ibcDenomRecordMap = new Map
        const ibcDenomRecord = await this.ibcDenomModel.findAllDenomRecord();
        if (ibcDenomRecord?.length) {
            ibcDenomRecord.forEach(ibcDenomRecordItem => {
                if (ibcDenomRecordItem?.chain_id) {
                    ibcDenomRecordMap.set(ibcDenomRecordItem?.chain_id, ibcDenomRecordItem)
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

    async parseIbcTx(dateNow): Promise<void> {
        const allChains = await this.chainConfigModel.findAll();
        const {allChainsMap, allChainsDenomPathsMap} = await this.getAllChainsMap()
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
                if (taskRecord.status === IbcTaskRecordStatus.CLOSE) return;
            }
            const taskModel = await this.connection.model(
                'txModel',
                SyncTaskSchema,
                `sync_${chain_id}_task`,
            );
            const taskCount = await getTaskStatus(taskModel,TaskEnum.tx)
            if(!taskCount) continue

            const txs = await this.getRecordLimitTx(chain_id, taskRecord.height, RecordLimit)
            let {handledTx, denoms} = await this.handlerSourcesTx(txs, chain_id, dateNow, allChainsMap, allChainsDenomPathsMap)


            if (handledTx?.length) {
                await this.ibcTxModel.insertManyIbcTx(handledTx, async err => {

                    taskRecord.height = handledTx[handledTx.length - 1]?.sc_tx_info?.height;
                    await this.ibcTaskRecordModel.updateTaskRecord(taskRecord);
                })
                // todo Transaction
                // await session.commitTransaction();
                // session.endSession();
            }
            if (denoms?.length) {
                await this.ibcDenomModel.insertManyDenom(denoms);
            }
        }
    }

    async getProcessingTxs() {
        const ibcTxs = await this.ibcTxModel.queryTxList({
            status: IbcTxStatus.PROCESSING,
            limit: RecordLimit,
        });
        return ibcTxs
    }

    async getPacketIds(txs) {
        const packetIds = []
        if (txs?.length) {
            for (const tx of txs) {
                if (tx?.sc_tx_info?.msg?.msg?.packet_id) {
                    // ibcTx.sc_tx_info.msg.msg.timeout_height.revision_height;
                    if(tx?.dc_chain_id
                        && tx?.height
                        && tx?.sc_tx_info?.msg?.msg?.packet_id
                        && tx?.sc_tx_info?.msg?.msg?.timeout_height?.revision_height
                        && tx?.sc_tx_info?.msg?.msg?.timeout_timestamp
                    ){
                        packetIds.push(
                            {
                                chainId: tx.dc_chain_id,
                                height: tx.height,
                                packetId: tx.sc_tx_info.msg.msg.packet_id,
                                timeOutTime :tx.sc_tx_info.msg.msg.timeout_timestamp
                            }
                        )
                    }
                }

            }
        }
        return packetIds
    }
    async changeIbcTxState(dateNow): Promise<void> {
        const ibcTxs = await this.getProcessingTxs()
        let packetIdArr = ibcTxs?.length ? await this.getPacketIds(ibcTxs) : [];
        let recvPacketTxMap = new Map, chainHeightMap = new Map, refundedTxTxMap = new Map,needUpdateTxs = [],denoms = [] //packetIdArr= [],
        // const allChains = await this.chainConfigModel.findAll();
        const dcChains = ibcTxs.map(item => {
            if (item.dc_chain_id) {
                return item.dc_chain_id
            }
        })
        const currentDcChains = Array.from(new Set(dcChains))
        // 根据链的配置信息，查询出每条链的 recv_packet 成功的交易的那条记录
        if (currentDcChains?.length) {
            for (const chain of currentDcChains) {
                if (chain) {
                    const txModel = await this.connection.model(
                        'txModel',
                        TxSchema,
                        `sync_${chain}_tx`,
                    );
                    const blockModel = await this.connection.model(
                        'blockModel',
                        IbcBlockSchema,
                        `sync_${chain}_block`,
                    );

                    const taskModel = await this.connection.model(
                        'txModel',
                        SyncTaskSchema,
                        `sync_${chain}_task`,
                    );
                    const taskCount = await getTaskStatus(taskModel,TaskEnum.tx)
                    if(!taskCount) continue
                    //每条链最新的高度
                    const chainHeightObj = await blockModel.findLatestBlock();
                    if (chainHeightObj && JSON.stringify(chainHeightObj) !== '{}') {
                        let {height, time} = await blockModel.findLatestBlock();
                        chainHeightMap.set(chain, {height, time})
                    }
                    const recvPacketIds = packetIdArr.map(item => {
                        return item.packetId
                    })
                    const refundedTxPacketIds = packetIdArr.map( item => {
                        if(item?.chainId && item?.height || item?.timeOutTime){
                            const currentChainLatestHeight = chainHeightMap.get(item.chainId)
                            const currentChainLatestBlockTime = chainHeightMap.get(item.chainId)
                            if(item.height < currentChainLatestHeight || item.timeOutTime < currentChainLatestBlockTime){
                                return item.packetId
                            }
                        }
                    })
                    // txs  数组
                    if (recvPacketIds?.length) {
                        const txs = await txModel.queryTxListByPacketId({
                            type: TxType.recv_packet,
                            limit: packetIdArr.length,
                            status: TxStatus.SUCCESS,
                            packet_id: recvPacketIds,
                        });
                        if (txs?.length) {
                            for (let tx of txs) {
                                if (tx?.msgs?.length) {
                                    for (let msg of tx.msgs) {
                                        if (msg?.type === TxType.recv_packet && msg.msg.packet_id) {
                                            recvPacketTxMap.set(`${chain}${msg.msg.packet_id}`, tx)
                                        }
                                    }
                                }
                            }
                        }
                    }
                    if(refundedTxPacketIds?.length){
                        const refundedTxs = await txModel.queryTxListByPacketId({
                            type: TxType.timeout_packet,
                            limit: packetIdArr.length,
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
                }
            }
        }

        for (let ibcTx of ibcTxs) {
            if (!ibcTx.dc_chain_id) continue
            if (recvPacketTxMap?.has(`${ibcTx.dc_chain_id}${ibcTx.sc_tx_info.msg.msg.packet_id}`)) {
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
                                        const resultObj = JSONparse(attribute.value);
                                        result = resultObj.hasOwnProperty('error')
                                            ? 'false'
                                            : 'true';
                                    }
                                });
                            }
                        });
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
                        if(ibcTx.status === IbcTxStatus.SUCCESS){
                            const denomObj = {
                                chain_id:ibcTx.dc_chain_id,
                                denom: dc_denom,
                                base_denom: ibcTx.base_denom,
                                denom_path: denom_path,
                                is_source_chain: !Boolean(denom_path),
                                create_at: dateNow,
                                update_at: ''
                            }
                            denoms.push(denomObj)
                        }

                    }
                })

            } else {
                /*
                * 没有找到的结果
                * */
                if (refundedTxTxMap.has(`${ibcTx.dc_chain_id}${ibcTx.sc_tx_info.msg.msg.packet_id}`)) {
                        const refunded_tx = refundedTxTxMap?.get(`${ibcTx.dc_chain_id}${ibcTx.sc_tx_info.msg.msg.packet_id}`);
                        refunded_tx &&
                        refunded_tx.msgs.forEach(msg => {
                            if (
                                msg.type === TxType.timeout_packet &&
                                ibcTx.sc_tx_info.msg.msg.packet_id === msg.msg.packet_id
                            ) {
                                ibcTx.status = IbcTxStatus.REFUNDED;
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
                                needUpdateTxs.push(ibcTx)
                            }
                        });
                    }
                }
            }
        if (needUpdateTxs?.length) {
            for (let needUpdateTx of needUpdateTxs) {
                await this.ibcTxModel.updateIbcTx(needUpdateTx);
            }
        }
        if(denoms?.length){
            await this.ibcDenomModel.insertManyDenom(denoms);
        }
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
                        default:
                            break;
                    }
                });
            }
        });
        return msg;
    }
    // parse Denom
    async parseDenom(
        chain_id,
        denom,
        base_denom,
        denom_path,
        is_source_chain,
        create_at,
        update_at,
        tx_time,
        is_base_denom,
        real_denom,
    ): Promise<void> {
        const ibcDenomRecord = await this.ibcDenomModel.findDenomRecord(
            chain_id,
            denom,
        );

        if (!ibcDenomRecord) {
            const ibcDenom = {
                chain_id,
                denom,
                base_denom,
                denom_path,
                is_source_chain,
                create_at,
                update_at,
                tx_time,
                is_base_denom,
                real_denom,
            };
            await this.ibcDenomModel.insertManyDenom(ibcDenom);
        } else {
            ibcDenomRecord.update_at = update_at;
            ibcDenomRecord.real_denom = ibcDenomRecord.real_denom || real_denom
            const currentTime = ibcDenomRecord.tx_time;
            if (tx_time > currentTime) {
                ibcDenomRecord.tx_time = tx_time;
            }
            await this.ibcDenomModel.updateDenomRecord(ibcDenomRecord);
        }
    }

    // get configed channels
    async getChannelsConfig() {
        const channels_all_record = [];

        const allChains = await this.chainConfigModel.findAll();
        allChains.forEach(chain => {
            chain.ibc_info.forEach(ibc_info_item => {
                ibc_info_item.paths.forEach(path_item => {
                    channels_all_record.push({
                        channel_id: path_item.channel_id,
                        record_id: `${chain.chain_id}${path_item.channel_id}`,
                        state: path_item.state,
                    });
                });
            });
        });

        return channels_all_record;
    }

    async handlerSourcesTx(sourcesTx, chain_id, currentTime, allChainsMap, allChainsDenomPathsMap) {
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
                        let lcd = '';
                        //根据 sc_chain_id  sc_port sc_channel 查询目标链的信息
                        /*
                        * 能否通过管道直接查询出目标链的数据
                        * 不要在循环里查询数据库 单独拎出来
                        * todo 失败状态处理
                        * */
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
                        /* const dcChainConfig = await this.chainConfigModel.findDcChain({
                             sc_chain_id,
                             sc_port,
                             sc_channel,
                         });*/

                        if (
                            dcChainConfig &&
                            dcChainConfig.ibc_info &&
                            dcChainConfig.ibc_info.length
                        ) {
                            lcd = dcChainConfig.lcd;
                            // dcChainConfig.ibc_info.forEach(info_item => {
                            //     info_item.paths.forEach(path_item => {
                            //         if (
                            //             path_item.channel_id === sc_channel &&
                            //             path_item.port_id === sc_port
                            //         ) {
                            //             dc_chain_id = info_item.chain_id;
                            //             dc_channel = path_item.counterparty.channel_id;
                            //             dc_port = path_item.counterparty.port_id;
                            //         }
                            //     });
                            // });
                        }
                        /*
                        * 通过lcd 查询的话可以不可以在起个定时任务去查询
                        * */

                        if (ibcTx.status === IbcTxStatus.FAILED) {
                            // get base_denom、denom_path from lcd API
                            if (sc_denom.indexOf('ibc') !== -1) {
                                const result = await ChainHttp.getDenomByLcdAndHash(
                                    lcd,
                                    sc_denom.replace('ibc/', ''),
                                );
                                base_denom = result?.base_denom;
                                denom_path = result?.denom_path;
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
                        let isBaseDenom = true
                        if (Boolean(denom_path) && denom_path.split('/').length > 1) {
                            const dc_port = denom_path.split('/')[0];
                            const dc_channel = denom_path.split('/')[1];
                            if (allChainsDenomPathsMap.has(sc_chain_id)) {
                                const dcDenomPath = allChainsDenomPathsMap.get(sc_chain_id)
                                const currentIbcTxDenomPath = `${dc_channel}${dc_port}`
                                if (dcDenomPath === currentIbcTxDenomPath) {
                                    isBaseDenom = false;
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
                                create_at: dateNow,
                                update_at: ''
                            })
                        }
                        // console.log(JSON.stringify(ibcTx),'这个东西是啥--------------------')
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
