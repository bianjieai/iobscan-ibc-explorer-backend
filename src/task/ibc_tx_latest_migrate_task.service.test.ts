import {IbcTxLatestMigrateTaskService} from "./ibc_tx_latest_migrate_task.service";
import {Test} from "@nestjs/testing";
import {AppModule} from "../app.module";

describe('IbcTxLatestMigrateTaskService', () => {
    let ibcTxLatestMigrateTaskService: IbcTxLatestMigrateTaskService;
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports: [
                AppModule
            ]
        }).compile();
        ibcTxLatestMigrateTaskService = module.get<IbcTxLatestMigrateTaskService>(IbcTxLatestMigrateTaskService);
    })

    describe('migrateData', () => {
        it('migrateData Test', async () => {
            jest.setTimeout(1000000)
            await ibcTxLatestMigrateTaskService.migrateData(10)
            console.log("=====")
        });

    });

})