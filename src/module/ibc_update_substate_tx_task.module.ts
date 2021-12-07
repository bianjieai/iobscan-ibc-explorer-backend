import { Module } from '@nestjs/common';
import {IbcTxHandler} from "../util/IbcTxHandler";
import {IbcUpdateSubStateTxTaskService} from "../task/ibc_update_substate_tx_task.service";
@Module({
    providers: [IbcUpdateSubStateTxTaskService,IbcTxHandler],
    exports: [IbcUpdateSubStateTxTaskService],
})
export class IbcUpdateSubstateTxTaskModule {}