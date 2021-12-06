import { Module } from '@nestjs/common';
import {TaskCommonService} from "../util/taskCommonService";
import {IbcUpdateProcessingTxTaskService} from "../task/ibc_update_processing_tx_task.service";
@Module({
    providers: [IbcUpdateProcessingTxTaskService,TaskCommonService],
    exports: [IbcUpdateProcessingTxTaskService],
})
export class IbcUpdateProcessingTxModule {}