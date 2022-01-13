import { Module } from '@nestjs/common';
import {IbcTxLatestMigrateTaskService} from "../task/ibc_tx_latest_migrate_task.service";
@Module({
    providers: [IbcTxLatestMigrateTaskService],
    exports: [IbcTxLatestMigrateTaskService],
})
export class IbcTxLatestMigrateTaskModule {}