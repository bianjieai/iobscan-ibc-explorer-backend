import { Module } from '@nestjs/common';
import {IbcTxHandler} from "../util/IbcTxHandler";
import {IbcTxDataUpdateTaskService} from "../task/ibc_tx_data_update_task.service";
@Module({
    providers: [IbcTxDataUpdateTaskService,IbcTxHandler],
    exports: [IbcTxDataUpdateTaskService],
})
export class IbcTxDataUpdateModule {}