import {IbcChainConfigTaskService} from "./ibc_chain_config.task.service";
import {Test} from "@nestjs/testing";
import {AppModule} from "../app.module";

describe('IbcChainConfigTaskService', () => {
    let ibcChainConfigTaskService: IbcChainConfigTaskService;
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports: [
                AppModule
            ]
        }).compile();
        ibcChainConfigTaskService = module.get<IbcChainConfigTaskService>(IbcChainConfigTaskService);
    })


    describe('parseChainConfig', () => {
        it('parseChainConfig Test', async () => {
            const result = await ibcChainConfigTaskService.parseChainConfig()
            console.log(result,'----')
        });
    });

    describe('findAllConfig', () => {
        it('findAllConfig Test', async () => {
            const allChains = await ibcChainConfigTaskService.findAllConfig()
            console.log(allChains,'--allChains--')
        });
    });


    describe('handleChain', () => {
        it('handleChain Test', async () => {
            const allChains = await ibcChainConfigTaskService.handleChain()
            console.log(allChains,'--allChains--')
        });
    });
})