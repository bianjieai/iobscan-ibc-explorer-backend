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
            const result = await ibcChainConfigTaskService.findAllConfig()
            console.log(result,'----')
        });
    });
})