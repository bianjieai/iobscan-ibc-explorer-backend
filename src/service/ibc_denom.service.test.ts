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

    describe('updateIbcDenom', () => {
        it('updateIbcDenom Test', async () => {
            const result = await ibcDenomService.updateIbcDenom({
                chain_id:"microtick_1",
                denom:"ibc/5F78C42BCC76287AE6B3185C6C1455DFFF8D805B1847F94B9B625384B93885C7",
                symbol:"ATOM"
            })
            console.log(result, '----')
        });
    });
})