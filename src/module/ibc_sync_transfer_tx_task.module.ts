import { Module } from '@nestjs/common';
import {IbcSyncTransferTxTaskService} from "../task/ibc_sync_transfer_tx_task.service";
import {IbcTxHandler} from "../util/IbcTxHandler";
@Module({
    providers: [IbcSyncTransferTxTaskService,IbcTxHandler],
    exports: [IbcSyncTransferTxTaskService],
})
export class IbcSyncTransferTxTaskModule {}