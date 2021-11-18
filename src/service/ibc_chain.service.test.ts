import {IbcChainService} from "./ibc_chain.service";
import {Test} from "@nestjs/testing";
import {AppModule} from "../app.module";

describe('IbcChainService', () => {
    let ibcChainService: IbcChainService;
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports: [
                AppModule
            ]
        }).compile();
        ibcChainService = module.get<IbcChainService>(IbcChainService);
    })

    describe('getAllChainConfigs', () => {
        it('getAllChainConfigs Test', async () => {
            const result = await ibcChainService.getAllChainConfigs()
            console.log(result, '----')
        });
    });

    describe('findActiveChains24hr', () => {
        it('findActiveChains24hr Test', async () => {
            const result = await ibcChainService.findActiveChains24hr(Math.floor(1623955689))
            console.log(result, '----')
        });
    });

    describe('queryChains', () => {
        it('queryChains Test', async () => {
            const result = await ibcChainService.queryChainsByDatetime(Math.floor(1623955689))
            console.log(result, '----')
        });
    });

    describe('handleActiveChains', () => {
        it('handleActiveChains Test', async () => {
            const allChainConfigs = await ibcChainService.getAllChainConfigs()

            const result = await ibcChainService.handleActiveChains(Math.floor(1623955689),allChainConfigs)
            console.log(result, '----')
        });
    });
})