import {Injectable} from '@nestjs/common';
import {SubState, TaskEnum} from "../constant";
import {IbcTxHandler} from "../util/IbcTxHandler";
import {dateNow} from "../helper/date.helper";

@Injectable()
export class IbcTxDataUpdateTaskService {
    constructor(private readonly taskCommonService: IbcTxHandler) {
        this.doTask = this.doTask.bind(this);
    }
    async doTask(taskName?: TaskEnum): Promise<void> {
        const defaultSubstate = 0
        const substate = [defaultSubstate,SubState.SuccessRecvPacketNotFound,SubState.RecvPacketAckFailed,SubState.SuccessTimeoutPacketNotFound]
        const ibcTxModel = this.taskCommonService.getIbcTxModel()
        await this.taskCommonService.changeIbcTxState(ibcTxModel,dateNow,substate)
    }
}