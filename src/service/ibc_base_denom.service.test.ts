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

    describe('insertBaseDenom', () => {
        it('insertBaseDenom Test', async () => {
            let data = {
                chain_id:"irishub_2",
                denom:"CAT",
                symbol:"TomCat",
                scale:8,
                icon:"",
                is_main_token:false,
            }
            await ibcBaseDenomService.insertBaseDenom(data)

        });
    });


})