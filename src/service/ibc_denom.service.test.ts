import {IbcDenomService} from "./ibc_denom.service";
import {Test} from "@nestjs/testing";
import {AppModule} from "../app.module";

describe('IbcDenomService', () => {
    let ibcDenomService: IbcDenomService;
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports: [
                AppModule
            ]
        }).compile();
        ibcDenomService = module.get<IbcDenomService>(IbcDenomService);
    })

    describe('findAllRecord', () => {
        it('findAllRecord Test', async () => {
            const result = await ibcDenomService.findAllRecord()
            console.log(result, '----')
        });
    });
})