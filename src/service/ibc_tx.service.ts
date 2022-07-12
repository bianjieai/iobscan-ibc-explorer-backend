/* eslint-disable @typescript-eslint/camelcase */
import {Injectable, Logger} from '@nestjs/common';
import {Connection} from 'mongoose';
import {InjectConnection} from '@nestjs/mongoose';
import {ListStruct, Result} from '../api/ApiResult';
import {IbcTxDetailsResDto, IbcTxListReqDto, IbcTxResDto, TxWithHashReqDto} from '../dto/ibc_tx.dto';
import {IbcDenomSchema} from '../schema/ibc_denom.schema';
import {IbcBaseDenomSchema} from '../schema/ibc_base_denom.schema';
import {IbcTxSchema} from '../schema/ibc_tx.schema';
import {unAuth, TaskEnum, IbcTxTable, TxType,} from '../constant';
import {cfg} from '../config/config';
import {IbcTxQueryType, IbcTxType} from "../types/schemaTypes/ibc_tx.interface";
import {IbcStatisticsSchema} from "../schema/ibc_statistics.schema";
import {TxSchema} from "../schema/tx.schema";


@Injectable()
export class IbcTxService {
    private ibcDenomModel;
    private ibcBaseDenomModel;
    private ibcTxLatestModel;
    private ibcTxModel;
    private ibcStatisticsModel;

    constructor(@InjectConnection() private connection: Connection) {
        this.getModels();
    }

    async getModels(): Promise<void> {
        this.ibcTxLatestModel = await this.connection.model(
            'ibcTxLatestModel',
            IbcTxSchema,
            IbcTxTable.IbcTxLatestTableName,
        );
        this.ibcTxModel = await this.connection.model(
            'ibcTxModel',
            IbcTxSchema,
            IbcTxTable.IbcTxTableName,
        );
        this.ibcDenomModel = await this.connection.model(
            'ibcDenomModel',
            IbcDenomSchema,
            'ibc_denom',
        );
        this.ibcBaseDenomModel = await this.connection.model(
            'ibcBaseDenomModel',
            IbcBaseDenomSchema,
            'ibc_base_denom',
        );
        // ibcStatisticsModel
        this.ibcStatisticsModel = await this.connection.model(
            'ibcStatisticsModel',
            IbcStatisticsSchema,
            'ibc_statistics',
        );
    }

    async getStartTxTime(): Promise<number> {
        const startTx = await this.ibcTxLatestModel.findFirstTx()
        return startTx?.tx_time;
    }

    async getTxCount(query: IbcTxQueryType, token): Promise<number> {
        const count = await this.ibcTxLatestModel.countTxList({...query, token});
        if (count >= cfg.serverCfg.displayIbcRecordMax) {
            return cfg.serverCfg.displayIbcRecordMax
        }
        return count
    }

    async getIbcTxs(query: IbcTxQueryType, token): Promise<IbcTxType[]> {
        return await this.ibcTxLatestModel.findTxList({...query, token})
    }

    async findStatisticTxsCount(): Promise<number> {
        const statisticData = await this.ibcStatisticsModel.findStatisticsRecord(
            TaskEnum.staticsTxAll,
        );
        if (statisticData.count_latest >= cfg.serverCfg.displayIbcRecordMax) {
            return cfg.serverCfg.displayIbcRecordMax
        }
        return statisticData.count_latest
    }

    async getTokenBySymbol(symbol): Promise<any> {
        const result = await this.ibcDenomModel.findRecordBySymbol(symbol);
        return result.map(item => {
            return item.base_denom;
            // return {
            //     denom: item.denom,
            //     base_denom: item.base_denom,
            //     chain_id: item.chain_id
            // };
        });
    }

    async getBaseDenomMap() : Promise<any>{
        let  baseDenomMap = new Map
        const baseDenom = await this.ibcBaseDenomModel.findAllRecord()
        if (baseDenom?.length) {
            for (const item of baseDenom) {
                baseDenomMap.set(`${item?.denom}`,item?.symbol)
            }
        }
        return baseDenomMap
    }


    async queryIbcTxList(
        query: IbcTxListReqDto,
    ): Promise<ListStruct<IbcTxResDto[]> | number> {
        const {use_count, page_num, page_size, symbol, denom, start_time} = query;
        const date_range = query?.date_range?.split(","),
            status = query?.status?.split(",") || [1, 2, 3, 4]
        let queryData: IbcTxQueryType = {
            useCount: query.use_count,
            date_range: [],
            chain_id: query.chain_id,
            status: [],
            // token?: { denom: string; chain_id: string }[];
            page_num: page_num,
            page_size: page_size > cfg.serverCfg.maxPageSize ? cfg.serverCfg.maxPageSize : page_size,
        }
        if (date_range?.length && !date_range.includes("0")) {
            for (const one of date_range) {
                queryData?.date_range.push(Number(one))
            }
        }
        for (const one of status) {
            queryData?.status.push(Number(one))
        }
        let token = undefined;
        if (symbol === unAuth) {
            token = await this.getTokenBySymbol('')
            // filter token which base_denom exist in ibc_base_denom
            if (token.length) {
                const baseDenomMap = await this.getBaseDenomMap()

                let  tokensFilter = []
                for  (const one of token) {
                    //only push token which base_denom not in ibc_base_denom
                    // if (baseDenomMap && !baseDenomMap?.has(`${one?.base_denom}`)) {
                    if (baseDenomMap && !baseDenomMap?.has(one)) {
                        tokensFilter.push(one)
                    }
                }

                if (tokensFilter?.length) {
                    token = [...tokensFilter]
                }
            }

        } else if (symbol) {
            token = await this.getTokenBySymbol(symbol)
            //no found the symbol
            if (!token) {
                return new ListStruct(null, page_num, page_size);
            }
        }
        if (denom) {
            token = [denom];
        }
        if (start_time) {
            // const startTx = await this.ibcTxModel.findFirstTx()
            // return startTx?.tx_time;
            //todo this value get by setting data
            return await this.getStartTxTime();
        }
        if (query?.chain_id) {
            const chains:string[] = query?.chain_id?.split(",")
            if (chains?.length > 2) {
                return new ListStruct(null, page_num, page_size);
            }
        }
        if (use_count) {
            if (query.symbol || query.chain_id || query.denom || (!queryData.date_range.includes(0)) || (queryData.status?.length !== 4)) {
                return await this.getTxCount(queryData, token)
            }
            // get statistic data
            return await this.findStatisticTxsCount()
            // return await this.ibcTxModel.countTxList({...query, token});
        } else {
            const ibcTxDatas: IbcTxResDto[] = IbcTxResDto.bundleData(
                // await this.ibcTxModel.findTxList({...query, token}),
                await this.getIbcTxs(queryData, token),
            );
            return new ListStruct(ibcTxDatas, page_num, page_size);
        }
    };

    async getConnectByTransferEventNews(eventNews,txMsgIndex) {
        let connect = '', timeout_timestamp = ''
        if (eventNews?.events_new?.length) {
            eventNews.events_new.forEach(item => {
                if (item?.events?.length && item?.msg_index === txMsgIndex) {
                    item.events.forEach(event => {
                        if (event?.type === 'send_packet') {
                            if (event?.attributes?.length) {
                                event.attributes.forEach(attribute => {
                                    switch (attribute.key) {
                                        case 'packet_connection':
                                            connect = attribute.value
                                            break
                                        case 'packet_timeout_timestamp':
                                            timeout_timestamp = attribute.value
                                            break
                                    }
                                })
                            }
                        }
                    })
                }
            })
        }
        return {connect, timeout_timestamp}
    }

    async getConnectByRecvPacketEventsNews(eventNews,txMsgIndex) {
        let connect = '', ackData = ''
        if (eventNews?.events_new?.length) {
            eventNews.events_new.forEach(item => {
                if (item?.events?.length && item?.msg_index === txMsgIndex) {
                    item.events.forEach(event => {
                        if (event?.type === 'write_acknowledgement' || event?.type === 'recv_packet') {
                            if (event?.attributes?.length) {
                                event.attributes.forEach(attribute => {
                                    switch (attribute.key) {
                                        case 'packet_connection':
                                            connect = attribute.value
                                            break
                                        case 'packet_ack':
                                            ackData = attribute.value
                                            break
                                    }
                                })
                            }
                        }
                    })
                }
            })
        }
        return {connect, ackData}
    }

    async getMsgIndex(tx,txType,packetId) :Promise<number>{
        for (const index in tx?.msgs) {
            if (tx?.msgs[index]?.type === txType && tx?.msgs[index]?.msg?.packet_id === packetId){
                return Number(index)
            }
        }
        return -1
    }
    async getScTxInfo(scChainID, scTxHash,packetId) {
        let scSigners = null, scConnect = null, timeOutTimestamp = null;
        if (scChainID && scTxHash) {
            const txModel = await this.connection.model(
                'txModel',
                TxSchema,
                `sync_${scChainID}_tx`,
            );
            let scTxData = await txModel.queryTxByHash(scTxHash)

            if (scTxData?.length) {
                const scTx = scTxData[scTxData?.length - 1]
                scSigners = scTx?.signers
                if (scTx?.events_new) {
                    const txMsgIndex = await this.getMsgIndex(scTx,TxType.transfer, packetId)
                    const {connect, timeout_timestamp} = await this.getConnectByTransferEventNews(scTx,txMsgIndex)
                    scConnect = connect
                    timeOutTimestamp = timeout_timestamp
                }

            }
        }
        return {
            scSigners,
            scConnect,
            timeOutTimestamp,
        }
    }

    async getDcTxInfo(dcChainID, dcTxHash, packetId) {
        let ack = null, dcConnect = null;
        if (dcChainID && dcTxHash) {
            const txModel = await this.connection.model(
                'txModel',
                TxSchema,
                `sync_${dcChainID}_tx`,
            );
            let dcTxData = await txModel.queryTxByHash(dcTxHash)

            if (dcTxData?.length) {
                const dcTx = dcTxData[dcTxData?.length - 1]
                if (dcTx?.events_new){
                    const txMsgIndex = await this.getMsgIndex(dcTx,TxType.recv_packet, packetId)
                    const {connect, ackData} = await this.getConnectByRecvPacketEventsNews(dcTx,txMsgIndex);
                    dcConnect = connect
                    ack = ackData
                }
            }
        }
        return {
            ack,
            dcConnect
        }
    }

    async getIbcTxDetail(txHash) {
        let ibcTx = await this.ibcTxLatestModel.queryTxDetailByHash(txHash)
        if (ibcTx.length === 0) {
            ibcTx =  await this.ibcTxModel.queryTxDetailByHash(txHash)
        }
        let ibcTxDetail = [],setMap = new Map
        for (const one of ibcTx) {
            if (setMap.has(one.record_id)) {
                continue
            }
            ibcTxDetail.push(one)
            setMap.set(one.record_id,'')
        }
        return ibcTxDetail
    }

    async getTxInfo(tx) {
        const packetId = `${tx.sc_port}${tx.sc_channel}${tx.dc_port}${tx.dc_channel}${tx.sequence}`
        if (tx.sc_chain_id && tx?.sc_tx_info?.hash) {
            const {scSigners, scConnect} = await this.getScTxInfo(tx.sc_chain_id, tx?.sc_tx_info?.hash, packetId)
            if (scSigners) {
                tx.sc_signers = scSigners
            }
            if (scConnect) {
                tx.sc_connect = scConnect
            }
        }

        if (tx.dc_chain_id && tx?.dc_tx_info?.hash) {
            const {ack, dcConnect} = await this.getDcTxInfo(tx.dc_chain_id, tx?.dc_tx_info?.hash, packetId)
            if (ack) {
                tx.dc_tx_info.ack = ack
            }
            if (dcConnect) {
                tx.dc_connect = dcConnect
            }
        }
        return tx
    }


    async queryIbcTxDetailsByHash(query: TxWithHashReqDto): Promise<IbcTxDetailsResDto[]> {
        let txDetailsData = await this.getIbcTxDetail(query)
        //todo not support detail return many data currently
        if (txDetailsData.length > 1) {
            return []
        }

        for (let tx of txDetailsData) {
            tx = await this.getTxInfo(tx)
            tx.dc_signers = []
            if (tx?.dc_tx_info?.msg?.msg?.signer) {
                tx.dc_signers.push(tx?.dc_tx_info?.msg?.msg?.signer)
            }
        }
        return IbcTxDetailsResDto.bundleData(
            txDetailsData,
        );
    }
}
