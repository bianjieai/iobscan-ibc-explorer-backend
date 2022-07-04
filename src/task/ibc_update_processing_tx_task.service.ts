import {Injectable} from '@nestjs/common';
import {TaskEnum} from "../constant";
import {IbcTxHandler} from "../util/IbcTxHandler";
@Injectable()
export class IbcUpdateProcessingTxTaskService {
    constructor(private readonly taskCommonService: IbcTxHandler) {
        this.doTask = this.doTask.bind(this);
    }
    async doTask(taskName?: TaskEnum): Promise<void> {
        const defaultSubstate = 0
        const ibcTxLatestModel = this.taskCommonService.getIbcTxLatestModel()
        const dateNow = Math.floor(new Date().getTime() / 1000)
        await this.taskCommonService.changeIbcTxState(ibcTxLatestModel, dateNow,[defaultSubstate],false,[])
    }
}