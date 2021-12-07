import {Injectable} from '@nestjs/common';
import {TaskEnum} from "../constant";
import {IbcTxHandler} from "../util/IbcTxHandler";
import {dateNow} from "../helper/date.helper";

@Injectable()
export class IbcSyncTransferTxTaskService {
    constructor(private readonly taskCommonService: IbcTxHandler) {
        this.doTask = this.doTask.bind(this);
    }
    async doTask(taskName?: TaskEnum): Promise<void> {
        await this.taskCommonService.parseIbcTx(dateNow)
    }
}