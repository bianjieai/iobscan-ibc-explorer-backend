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
            jest.setTimeout(100000000)
            const query: IbcTxListReqDto = { page_num: 1, page_size: 10,use_count:true,
                status:"1,2,3,4",
                date_range:"0,1642495903",
                symbol:'ATOM',
                // chain_id:'cosmoshub_4',
            };
            const time1 = Math.floor(new Date().getTime()/1000)

            const result = await ibcTxService.queryIbcTxList(query)
            console.log("time cost====>:",Math.floor(new Date().getTime()/1000) - time1)

            console.log('====count==>>:',result)
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
    describe('getConnectByTransferEventNews', () => {
        it('getConnectByTransferEventNews Test', async () => {
            const eventNews = {
                events_new: [
                    {
                        "msg_index":0,
                        "events":[
                            {
                                "type":"coin_received",
                                "attributes":[
                                    {
                                        "key":"receiver",
                                        "value":"bostrom12k2pyuylm9t7ugdvz67h9pg4gmmvhn5vu43p0n"
                                    },
                                    {
                                        "key":"amount",
                                        "value":"965326boot"
                                    }
                                ]
                            },
                            {
                                "type":"coin_spent",
                                "attributes":[
                                    {
                                        "key":"spender",
                                        "value":"bostrom1et80vz76fn5w946z864cg9j8yuwp298gc2n826"
                                    },
                                    {
                                        "key":"amount",
                                        "value":"965326boot"
                                    }
                                ]
                            },
                            {
                                "type":"ibc_transfer",
                                "attributes":[
                                    {
                                        "key":"sender",
                                        "value":"bostrom1et80vz76fn5w946z864cg9j8yuwp298gc2n826"
                                    },
                                    {
                                        "key":"receiver",
                                        "value":"osmo1et80vz76fn5w946z864cg9j8yuwp298gnz5yz0"
                                    }
                                ]
                            },
                            {
                                "type":"message",
                                "attributes":[
                                    {
                                        "key":"action",
                                        "value":"/ibc.applications.transfer.v1.MsgTransfer"
                                    },
                                    {
                                        "key":"sender",
                                        "value":"bostrom1et80vz76fn5w946z864cg9j8yuwp298gc2n826"
                                    },
                                    {
                                        "key":"module",
                                        "value":"ibc_channel"
                                    },
                                    {
                                        "key":"module",
                                        "value":"transfer"
                                    }
                                ]
                            },
                            {
                                "type":"send_packet",
                                "attributes":[
                                    {
                                        "key":"packet_data",
                                        "value":"{\"amount\":\"965326\",\"denom\":\"boot\",\"receiver\":\"osmo1et80vz76fn5w946z864cg9j8yuwp298gnz5yz0\",\"sender\":\"bostrom1et80vz76fn5w946z864cg9j8yuwp298gc2n826\"}"
                                    },
                                    {
                                        "key":"packet_data_hex",
                                        "value":"7b22616d6f756e74223a22393635333236222c2264656e6f6d223a22626f6f74222c227265636569766572223a226f736d6f3165743830767a3736666e35773934367a3836346367396a3879757770323938676e7a35797a30222c2273656e646572223a22626f7374726f6d3165743830767a3736666e35773934367a3836346367396a38797577703239386763326e383236227d"
                                    },
                                    {
                                        "key":"packet_timeout_height",
                                        "value":"1-2608471"
                                    },
                                    {
                                        "key":"packet_timeout_timestamp",
                                        "value":"0"
                                    },
                                    {
                                        "key":"packet_sequence",
                                        "value":"2114"
                                    },
                                    {
                                        "key":"packet_src_port",
                                        "value":"transfer"
                                    },
                                    {
                                        "key":"packet_src_channel",
                                        "value":"channel-2"
                                    },
                                    {
                                        "key":"packet_dst_port",
                                        "value":"transfer"
                                    },
                                    {
                                        "key":"packet_dst_channel",
                                        "value":"channel-95"
                                    },
                                    {
                                        "key":"packet_channel_ordering",
                                        "value":"ORDER_UNORDERED"
                                    },
                                    {
                                        "key":"packet_connection",
                                        "value":"connection-2"
                                    }
                                ]
                            },
                            {
                                "type":"transfer",
                                "attributes":[
                                    {
                                        "key":"recipient",
                                        "value":"bostrom12k2pyuylm9t7ugdvz67h9pg4gmmvhn5vu43p0n"
                                    },
                                    {
                                        "key":"sender",
                                        "value":"bostrom1et80vz76fn5w946z864cg9j8yuwp298gc2n826"
                                    },
                                    {
                                        "key":"amount",
                                        "value":"965326boot"
                                    }
                                ]
                            }
                        ]
                    }
                ]
            };
            const result = await ibcTxService.getConnectByTransferEventNews(eventNews,1)
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
            // console.log(result, '----')
            console.log(result.length, '----')
        });

        it('getScTxInfo Test', async () => {
            const result = await ibcTxService.getScTxInfo("qa_iris_snapshot","A7B69456C9C34B477FA021D6781F8F95A704BEC001532AF5D833354961573C98","transferchannel-1transferchannel-541")
            console.log(result, '----')
        });

        it('getDcTxInfo Test', async () => {
            const result = await ibcTxService.getDcTxInfo("cosmoshub_4","E3DFBCF3A15BF971C9C97CBCEBBD81C610B028CE18F285FDC3175A2968CFA2AB","transferchannel-12transferchannel-1823433")
            console.log(result, '----')
        });
    });

})