import {Test} from "@nestjs/testing";
import {AppModule} from "../app.module";
import {IbcDenomUpdateTaskService} from "./ibc_denom_update.task.service";


describe('IbcDenomUpdateTaskService', () => {
    let ibcDenomUpdateTaskService: IbcDenomUpdateTaskService;
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports: [
                AppModule
            ]
        }).compile();
        ibcDenomUpdateTaskService = module.get<IbcDenomUpdateTaskService>(IbcDenomUpdateTaskService);
    })

    describe('Uint Test', () => {
        it('getEmptySymbolDenom Test', async () => {
            jest.setTimeout(100000000)
            const AllDenom = await ibcDenomUpdateTaskService.getNeedhandleIbcDenoms(1,100)
            const ret = await ibcDenomUpdateTaskService.collectChainDenomsMap(AllDenom)
            console.log(ret,'--ret--')
        });

        it('getNeedhandleIbcDenoms Test', async () => {
            jest.setTimeout(100000000)
            const AllDenom = await ibcDenomUpdateTaskService.getNeedhandleIbcDenoms(1,100)
            console.log(AllDenom,'--AllDenom--')
        });

        it('handleChainDenoms Test', async () => {
            jest.setTimeout(100000000)
            await ibcDenomUpdateTaskService.handleChainDenoms()
        });

        it('handleChain Test', async () => {
            const ret = await ibcDenomUpdateTaskService.getIbcDenoms("bigbang",["ibc/FD01DE9421BAC4DA062DF60B0750248026D1EA64C7990A4BB58D91EC2BBEAB56","ibc/4BAD958A5FF565DEADCF82F5B5D95B823EE9A7796683D3CCFE31A761BF931513"])
            console.log(ret,"====ret===")
        });
    });
})