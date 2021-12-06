import { Module } from '@nestjs/common';
import {TaskCommonService} from "../util/taskCommonService";
import {IbcUpdateSubStateTxTaskService} from "../task/ibc_update_substate_tx_task.service";
@Module({
    providers: [IbcUpdateSubStateTxTaskService,TaskCommonService],
    exports: [IbcUpdateSubStateTxTaskService],
})
export class IbcUpdateSubstateTxTaskModule {}