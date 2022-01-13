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
    });
})