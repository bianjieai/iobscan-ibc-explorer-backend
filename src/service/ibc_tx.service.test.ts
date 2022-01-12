import {IbcTxService} from "./ibc_tx.service";
import {Test} from "@nestjs/testing";
import {AppModule} from "../app.module";
import {IbcTxListReqDto} from "../dto/ibc_tx.dto";
import {IbcTxQueryType} from "../types/schemaTypes/ibc_tx.interface";

describe('IbcTxService', () => {
    let ibcTxService: IbcTxService;
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports: [
                AppModule
            ]
        }).compile();
        ibcTxService = module.get<IbcTxService>(IbcTxService);
    })

    describe('queryIbcTxList', () => {
        it('queryIbcTxList Test', async () => {
            const query: IbcTxListReqDto = { page_num: 1, page_size: 10,use_count:true,
                // status:"1,2,3,4",
                // date_range:"0,1640074000"
            };
            const result = await ibcTxService.queryIbcTxList(query)
            console.log(result, '----')
        });
    });

    describe('getTokenBySymbol', () => {
        it('getTokenBySymbol Test', async () => {
            const symbol = "osmo"
            // const symbol = unAuth
            const result = await ibcTxService.getTokenBySymbol(symbol)
            console.log(result, '----')
        });
    });

    describe('getStartTxTime', () => {
        it('getStartTxTime Test', async () => {
            const result = await ibcTxService.getStartTxTime()
            console.log(result, '----')
        });
    });

    describe('getTxCount', () => {
        it('getTxCount Test', async () => {
            const query: IbcTxListReqDto = {status: "1", chain_id: "osmosis_1"};
            const token = await ibcTxService.getTokenBySymbol("osmo")
            const date_range = query.date_range?.split(","),status = query.status?.split(",")
            let queryData :IbcTxQueryType = {
                useCount: query.use_count,
                // date_range?: number[];
                chain_id: query.chain_id,
                // status?: number[];
                // token?: { denom: string; chain_id: string }[];
                page_num:1,
                page_size: 10,
            }
            for (const one of date_range) {
                queryData?.date_range.push(Number(one))
            }
            for (const one of status) {
                queryData?.status.push(Number(one))
            }
            const result = await ibcTxService.getTxCount(queryData, token)
            console.log(result, '----')
        });
    });

    describe('getIbcTxs', () => {
        it('getIbcTxs Test', async () => {
            const query: IbcTxListReqDto = {status: "1", page_num: 1, page_size: 10, chain_id: "osmosis_1"};
            const token = await ibcTxService.getTokenBySymbol("osmo")
            const date_range = query.date_range?.split(","),status = query.status?.split(",")
            let queryData :IbcTxQueryType = {
                useCount: query.use_count,
                // date_range?: number[];
                chain_id: query.chain_id,
                // status?: number[];
                // token?: { denom: string; chain_id: string }[];
                page_num:1,
                page_size: 10,
            }
            for (const one of date_range) {
                queryData?.date_range.push(Number(one))
            }
            for (const one of status) {
                queryData?.status.push(Number(one))
            }
            const result = await ibcTxService.getIbcTxs(queryData, token)
            console.log(result, '----')
        });
    });

    describe('findStatisticTxsCount', () => {
        it('findStatisticTxsCount Test', async () => {
            const result = await ibcTxService.findStatisticTxsCount()
            console.log(result, '----')
        });
    });

    describe('queryIbcTxDetailsByHash', () => {
        it('queryIbcTxDetailsByHash Test', async () => {
            const result = await ibcTxService.queryIbcTxDetailsByHash({hash:'A7B69456C9C34B477FA021D6781F8F95A704BEC001532AF5D833354961573C98'})
            console.log(result, '----')
        });

        it('getIbcTxDetail Test', async () => {
            const result = await ibcTxService.getIbcTxDetail({hash:'A7B69456C9C34B477FA021D6781F8F95A704BEC001532AF5D833354961573C98'})
            console.log(result, '----')
        });

        it('getScTxInfo Test', async () => {
            const result = await ibcTxService.getScTxInfo("qa_iris_snapshot","A7B69456C9C34B477FA021D6781F8F95A704BEC001532AF5D833354961573C98","transferchannel-1transferchannel-541")
            console.log(result, '----')
        });

        it('getDcTxInfo Test', async () => {
            const result = await ibcTxService.getDcTxInfo("bigbang","AF5CB680C80BC8B48158CD44537A08C4193AB4FDAE406A623F142443456EE6C1","transferchannel-1transferchannel-541")
            console.log(result, '----')
        });
    });

})