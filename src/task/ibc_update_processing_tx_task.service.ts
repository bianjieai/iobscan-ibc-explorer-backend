import {Injectable} from '@nestjs/common';
import {TaskEnum} from "../constant";
import {IbcTxHandler} from "../util/IbcTxHandler";
import {dateNow} from "../helper/date.helper";
@Injectable()
export class IbcUpdateProcessingTxTaskService {
    constructor(private readonly taskCommonService: IbcTxHandler) {
        this.doTask = this.doTask.bind(this);
    }
    async doTask(taskName?: TaskEnum): Promise<void> {
        const defaultSubstate = 0
        const ibcTxLatestModel = this.taskCommonService.getIbcTxLatestModel()
        await this.taskCommonService.changeIbcTxState(ibcTxLatestModel, dateNow,[defaultSubstate],false,[])
    }
}