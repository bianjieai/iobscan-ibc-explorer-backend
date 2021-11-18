import {IbcStatisticsTaskService} from "./ibc_statistics.task.service";
import {Test} from "@nestjs/testing";
import {AppModule} from "../app.module";
import {IbcTxTaskService} from "./ibc_tx.task.service"

describe('IbcTxTaskService', () => {
    let ibcTxTaskService: IbcTxTaskService;
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports: [
                AppModule
            ]
        }).compile();
        ibcTxTaskService = module.get<IbcTxTaskService>(IbcTxTaskService);
    })

    describe('parseIbcTx', () => {
        it('parseIbcTxBychain Test', async () => {
            const dateNow = Math.floor(1623955689);
            await ibcTxTaskService.parseIbcTx(dateNow)
            console.log('----')
        });
    });

    describe('getRecordLimitTx', () => {
        it('getRecordLimitTx Test', async () => {
            let chain_id = "emoney_3"
            const result  = await ibcTxTaskService.getRecordLimitTx(chain_id,13106,10)
            console.log(result,'----')
        });
    });

    describe('handlerSourcesTx', () => {
        it('handlerSourcesTx Test', async () => {
            let chain_id = "emoney_3"
            const txs  = await ibcTxTaskService.getRecordLimitTx(chain_id,13106,10)
            const allChainsMap = await ibcTxTaskService.getAllChainsMap()
            const dateNow = Math.floor(1623955689);
            const result  = await ibcTxTaskService.handlerSourcesTx(txs,chain_id,dateNow,allChainsMap)
            console.log(result,'----')
        });
    });
})