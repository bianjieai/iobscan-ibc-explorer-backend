import {Injectable} from '@nestjs/common';
import {SubState, TaskEnum} from "../constant";
import {IbcTxHandler} from "../util/IbcTxHandler";
import {dateNow} from "../helper/date.helper";

@Injectable()
export class IbcUpdateSubStateTxTaskService {
    constructor(private readonly taskCommonService: IbcTxHandler) {
        this.doTask = this.doTask.bind(this);
    }
    async doTask(taskName?: TaskEnum): Promise<void> {
        const substate = [SubState.SuccessRecvPacketNotFound,SubState.RecvPacketAckFailed,SubState.SuccessTimeoutPacketNotFound]
        const ibcTxLatestModel = this.taskCommonService.getIbcTxLatestModel()
        await this.taskCommonService.changeIbcTxState(ibcTxLatestModel,dateNow,substate,false)
    }
}