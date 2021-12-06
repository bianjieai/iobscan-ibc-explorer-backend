import {Injectable} from '@nestjs/common';
import {SubState, TaskEnum} from "../constant";
import {TaskCommonService} from "../util/taskCommonService";
import {dateNow} from "../helper/date.helper";
@Injectable()
export class IbcUpdateSubStateTxTaskService {
    constructor(private readonly taskCommonService: TaskCommonService) {
        this.doTask = this.doTask.bind(this);
    }
    async doTask(taskName?: TaskEnum): Promise<void> {
        const substate = [SubState.SuccessRecvPacketNotFound,SubState.SuccessTimeoutPacketNotFound]
        await this.taskCommonService.changeIbcTxState(dateNow,substate)
    }
}