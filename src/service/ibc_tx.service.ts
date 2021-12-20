/* eslint-disable @typescript-eslint/camelcase */
import {Injectable} from '@nestjs/common';
import {Connection} from 'mongoose';
import {InjectConnection} from '@nestjs/mongoose';
import {ListStruct} from '../api/ApiResult';
import {IbcTxListReqDto, IbcTxResDto} from '../dto/ibc_tx.dto';
import {IbcDenomSchema} from '../schema/ibc_denom.schema';
import {IbcTxSchema} from '../schema/ibc_tx.schema';
import {unAuth} from '../constant';
import {IbcTxQueryType, IbcTxType} from "../types/schemaTypes/ibc_tx.interface";
import {IbcStatisticsType} from "../types/schemaTypes/ibc_statistics.interface";
import {IbcStatisticsSchema} from "../schema/ibc_statistics.schema";

@Injectable()
export class IbcTxService {
    private ibcDenomModel;
    private ibcTxModel;
    private ibcStatisticsModel;

    constructor(@InjectConnection() private connection: Connection) {
        this.getModels();
    }

    async getModels(): Promise<void> {
        this.ibcTxModel = await this.connection.model(
            'ibcTxModel',
            IbcTxSchema,
            'ex_ibc_tx',
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
        const startTx = await this.ibcTxModel.findFirstTx()
        return startTx?.tx_time;
    }

    async getTxCount(query: IbcTxQueryType, token): Promise<number> {
        return await this.ibcTxModel.countTxList({...query, token});
    }

    async getIbcTxs(query: IbcTxQueryType, token): Promise<IbcTxType[]> {
        return await this.ibcTxModel.findTxList({...query, token})
    }

    async findStatisticTxsCount():Promise<IbcStatisticsType> {
        return await this.ibcStatisticsModel.findStatisticsRecord(
            'tx_all',
        );
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

            if (query.symbol || query.chain_id ||query.denom || (query.date_range?.length === 2)) {
                return await this.getTxCount(query,token)
            }
            // get statistic data
            const statisticData = await this.findStatisticTxsCount()
            return statisticData.count
            // return await this.ibcTxModel.countTxList({...query, token});
        } else {
            const ibcTxDatas: IbcTxResDto[] = IbcTxResDto.bundleData(
                // await this.ibcTxModel.findTxList({...query, token}),
                await this.getIbcTxs(query,token),
            );
            return new ListStruct(ibcTxDatas, page_num, page_size);
        }
    }
}
