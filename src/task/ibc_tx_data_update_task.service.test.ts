import {IbcTxDataUpdateTaskService} from "./ibc_tx_data_update_task.service";
import {Test} from "@nestjs/testing";
import {AppModule} from "../app.module";

describe('IbcTxLatestMigrateTaskService', () => {
    let ibcTxDataUpdateTaskService: IbcTxDataUpdateTaskService;
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports: [
                AppModule
            ]
        }).compile();
        ibcTxDataUpdateTaskService = module.get<IbcTxDataUpdateTaskService>(IbcTxDataUpdateTaskService);
    })
    describe('handleUpdateIbcTx', () => {
        it('handleUpdateIbcTx', async () => {
            jest.setTimeout(100000000)
            await ibcTxDataUpdateTaskService.handleUpdateIbcTx()
            // console.log(ibcTxTaskService,'----')
        });
    });
})