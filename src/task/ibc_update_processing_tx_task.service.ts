import {Injectable} from '@nestjs/common';
import {TaskEnum} from "../constant";
import {TaskCommonService} from "../util/taskCommonService";
import {dateNow} from "../helper/date.helper";
@Injectable()
export class IbcUpdateProcessingTxTaskService {
    constructor(private readonly taskCommonService: TaskCommonService) {
        this.doTask = this.doTask.bind(this);
    }
    async doTask(taskName?: TaskEnum): Promise<void> {
        console.log('怎么感觉有点冷了呢')
        const defaultSubstate = 0
        await this.taskCommonService.changeIbcTxState(dateNow,[defaultSubstate])
    }
}