import { Module } from '@nestjs/common';
import {IbcSyncTransferTxTaskService} from "../task/ibc_sync_transfer_tx_task.service";
import {IbcTxHandler} from "../util/IbcTxHandler";
import {TransferTaskStatusMetric, TransferTaskStatusProvider} from "../monitor/metrics/ibc_transfer_task_status.metric";
@Module({
    providers: [IbcSyncTransferTxTaskService,IbcTxHandler,TransferTaskStatusMetric,
        TransferTaskStatusProvider(),],
    exports: [IbcSyncTransferTxTaskService],
})
export class IbcSyncTransferTxTaskModule {}