import {IbcBaseDenomService} from "./ibc_base_denom.service";
import {Test} from "@nestjs/testing";
import {AppModule} from "../app.module";

describe('IbcDenomService', () => {
    let ibcBaseDenomService: IbcBaseDenomService;
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports: [
                AppModule
            ]
        }).compile();
        ibcBaseDenomService = module.get<IbcBaseDenomService>(IbcBaseDenomService);
    })

    describe('findAllRecord', () => {
        it('findAllRecord Test', async () => {
            const result = await ibcBaseDenomService.findAllRecord()
            console.log(result, '----')
        });
    });
})