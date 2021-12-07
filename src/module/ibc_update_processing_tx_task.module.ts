import { Module } from '@nestjs/common';
import {IbcTxHandler} from "../util/IbcTxHandler";
import {IbcUpdateProcessingTxTaskService} from "../task/ibc_update_processing_tx_task.service";
@Module({
    providers: [IbcUpdateProcessingTxTaskService,IbcTxHandler],
    exports: [IbcUpdateProcessingTxTaskService],
})
export class IbcUpdateProcessingTxModule {}