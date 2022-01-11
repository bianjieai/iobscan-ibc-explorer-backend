/* eslint-disable @typescript-eslint/camelcase */
import {Injectable, Logger} from '@nestjs/common';
import {Connection} from 'mongoose';
import {InjectConnection} from '@nestjs/mongoose';
import {ListStruct, Result} from '../api/ApiResult';
import {IbcTxDetailsResDto, IbcTxListReqDto, IbcTxResDto, TxWithHashReqDto} from '../dto/ibc_tx.dto';
import {IbcDenomSchema} from '../schema/ibc_denom.schema';
import {IbcTxSchema} from '../schema/ibc_tx.schema';
import {unAuth, TaskEnum, IbcTxTable} from '../constant';
import {cfg} from '../config/config';
import {IbcTxQueryType, IbcTxType} from "../types/schemaTypes/ibc_tx.interface";
import {IbcStatisticsType} from "../types/schemaTypes/ibc_statistics.interface";
import {IbcStatisticsSchema} from "../schema/ibc_statistics.schema";
import {TxSchema} from "../schema/tx.schema";
import {ChainHttp} from "../http/lcd/chain.http";

@Injectable()
export class IbcTxService {
    private ibcDenomModel;
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
        return await this.ibcTxModel.findTxList({...query, token})
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
            return {
                denom: item.denom,
                chain_id: item.chain_id
            };
        });
    }


    async queryIbcTxList(
        query: IbcTxListReqDto,
    ): Promise<ListStruct<IbcTxResDto[]> | number> {
        const {use_count, page_num, page_size, symbol, denom, start_time} = query;
        const date_range = query?.date_range?.split(",") || [0, new Date().getTime() / 1000],
            status = query?.status?.split(",") || [1, 2, 3, 4]
        let queryData: IbcTxQueryType = {
            useCount: query.use_count,
            date_range: [],
            chain_id: query.chain_id,
            status: [],
            // token?: { denom: string; chain_id: string }[];
            page_num: page_num,
            page_size: page_size,
        }
        for (const one of date_range) {
            queryData?.date_range.push(Number(one))
        }
        for (const one of status) {
            queryData?.status.push(Number(one))
        }
        let token = undefined;
        if (symbol === unAuth) {
            // const resultUnAuth = await this.ibcDenomModel.findRecordBySymbol('');
            // token = resultUnAuth.map(item => {
            //     return {
            //         denom: item.denom,
            //         chain_id: item.chain_id
            //     };
            // });
            token = await this.getTokenBySymbol('')
        } else if (symbol) {
            token = await this.getTokenBySymbol(symbol)
            // const result = await this.ibcDenomModel.findRecordBySymbol(symbol);
            // token = result.map(item => {
            //     return {
            //         denom: item.denom,
            //         chain_id: item.chain_id
            //     };
            // });
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

    async getConnectByTransferEventNews(eventNews) {
        let connect = '', timeout_timestamp = ''
        if (eventNews?.events_new?.length) {
            eventNews.events_new.forEach(item => {
                if (item?.events?.length) {
                    item.events.forEach(event => {
                        if (event?.type === 'send_packet') {
                            if (event?.attributes?.length) {
                                event.attributes.forEach(attribute => {
                                    switch (attribute.key) {
                                        case 'packet_connection':
                                            connect = attribute.value
                                        case 'packet_timeout_timestamp':
                                            timeout_timestamp = attribute.value
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

    async getConnectByRecvPacketEventsNews(eventNews) {
        let connect = '', ackData = ''
        if (eventNews?.events_new?.length) {
            eventNews.events_new.forEach(item => {
                if (item?.events?.length) {
                    item.events.forEach(event => {
                        if (event?.type === 'write_acknowledgement') {
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

    async getScTxInfo(scChainID, scTxHash) {
        let scSigners = null, scConnect = null, timeOutTimestamp = null;
        if (scChainID && scTxHash) {
            const txModel = await this.connection.model(
                'txModel',
                TxSchema,
                `sync_${scChainID}_tx`,
            );
            let scTxData = await txModel.queryTxByHash(scTxHash)

            if (scTxData?.length) {
                scSigners = scTxData[scTxData?.length - 1]?.signers
                if (scTxData[scTxData?.length - 1]?.events_new) {
                    const {connect, timeout_timestamp} = await this.getConnectByTransferEventNews(scTxData[scTxData?.length - 1])
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

    async getDcTxInfo(dcChainID, dcTxHash) {
        let ack = null, dcConnect = null;
        if (dcChainID && dcTxHash) {
            const txModel = await this.connection.model(
                'txModel',
                TxSchema,
                `sync_${dcChainID}_tx`,
            );
            let dcTxData = await txModel.queryTxByHash(dcTxHash)

            if (dcTxData?.length && dcTxData[dcTxData?.length - 1]?.events_new) {
                const {connect, ackData} = await this.getConnectByRecvPacketEventsNews(dcTxData[dcTxData?.length - 1]);
                dcConnect = connect
                ack = ackData
            }
        }
        return {
            ack,
            dcConnect
        }
    }

    async getIbcTxDetail(query) {
        return await this.ibcTxModel.queryTxDetailByHash(query)
    }

    async getTxInfo(tx) {
        if (tx.sc_chain_id && tx?.sc_tx_info?.hash) {
            const {scSigners, scConnect} = await this.getScTxInfo(tx.sc_chain_id, tx?.sc_tx_info?.hash)
            if (scSigners) {
                tx.sc_signers = scSigners
            }
            if (scConnect) {
                tx.sc_connect = scConnect
            }
        }

        if (tx.dc_chain_id && tx?.dc_tx_info?.hash) {
            const {ack, dcConnect} = await this.getDcTxInfo(tx.dc_chain_id, tx?.dc_tx_info?.hash)
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
