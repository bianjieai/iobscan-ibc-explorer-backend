import { Module } from '@nestjs/common';
import {IbcSyncTransferTxTaskService} from "../task/ibc_sync_transfer_tx_task.service";
import {TaskCommonService} from "../util/taskCommonService";
@Module({
    providers: [IbcSyncTransferTxTaskService,TaskCommonService],
    exports: [IbcSyncTransferTxTaskService],
})
export class IbcSyncTransferTxTaskModule {}