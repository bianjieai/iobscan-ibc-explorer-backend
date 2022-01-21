import {IbcDenomCaculateTaskService} from "./ibc_denom_caculate.task.service";
import {Test} from "@nestjs/testing";
import {AppModule} from "../app.module";


describe('IbcDenomHashTaskService', () => {
    let ibcDenomHashTaskService: IbcDenomCaculateTaskService;
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports: [
                AppModule
            ]
        }).compile();
        ibcDenomHashTaskService = module.get<IbcDenomCaculateTaskService>(IbcDenomCaculateTaskService);
    })

    describe('findConfig', () => {
        it('findAllBaseDenom Test', async () => {
            const AllBaseDenom = await ibcDenomHashTaskService.findAllBaseDenom()
            console.log(AllBaseDenom,'--findAllBaseDenom--')
        });

        it('findAllChainConfig Test', async () => {
            const allChains = await ibcDenomHashTaskService.findAllChainConfig()
            console.log(allChains,'--findAllChainConfig--')
        });

        it('handleChain Test', async () => {
            await ibcDenomHashTaskService.handleChain()
        });

        it('getCaculateDenomMap Test', async () => {
           const data = await ibcDenomHashTaskService.getCaculateDenomMap("sifchain_1")
            console.log('--data-->:',data)
        });

        it('caculateBaseDenom Test', async () => {
            const chainConfig = await ibcDenomHashTaskService.findAllChainConfig()
            const AllBaseDenom = await ibcDenomHashTaskService.findAllBaseDenom()
            let denomMap = new Map, channelMap = new Map

            for (const one of AllBaseDenom) {
                denomMap.set(`${one.chain_id}`, one)
            }
            for (const one of chainConfig) {
                channelMap.set(`${one.chain_id}`, one)
            }
            await ibcDenomHashTaskService.caculateBaseDenom(denomMap.get("cosmoshub_4"),channelMap)
        });

        it('findCount Test', async () => {
            const chainConfig = await ibcDenomHashTaskService.findAllChainConfig()
            let Total = 0
            for (const one of chainConfig) {
                let num = 0
                if (one?.ibc_info?.length > 0) {
                    for (const ibcInfo of one?.ibc_info) {
                        if (ibcInfo?.chain_id && ibcInfo?.paths?.length > 0) {
                            ibcInfo.paths.forEach(item => {
                                if (item?.counterparty?.port_id && item?.counterparty?.channel_id) {
                                    Total++
                                    num++
                                }
                            })
                        }
                    }
                }
                console.log("chain:",one.chain_id," total:",num)
            }
            console.log('--Total-->:',Total)
        });
    });
})