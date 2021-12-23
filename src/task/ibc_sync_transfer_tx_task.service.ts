import {Injectable} from '@nestjs/common';
import {TaskEnum} from "../constant";
import {IbcTxHandler} from "../util/IbcTxHandler";
import {dateNow} from "../helper/date.helper";
import {TransferTaskStatusMetric} from "../monitor/metrics/ibc_transfer_task_status.metric";

@Injectable()
export class IbcSyncTransferTxTaskService {
    constructor(private readonly taskCommonService: IbcTxHandler,private readonly transferTaskStatusMetric: TransferTaskStatusMetric) {
        this.doTask = this.doTask.bind(this);
    }
    async doTask(taskName?: TaskEnum): Promise<void> {
        this.transferTaskStatusMetric.collect(1)
        await this.taskCommonService.parseIbcTx(dateNow)
    }
}