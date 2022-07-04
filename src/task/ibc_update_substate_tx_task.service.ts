import {Injectable} from '@nestjs/common';
import {SubState, TaskEnum} from "../constant";
import {IbcTxHandler} from "../util/IbcTxHandler";

@Injectable()
export class IbcUpdateSubStateTxTaskService {
    constructor(private readonly taskCommonService: IbcTxHandler) {
        this.doTask = this.doTask.bind(this);
    }
    async doTask(taskName?: TaskEnum): Promise<void> {
        const substate = [SubState.SuccessRecvPacketNotFound,SubState.RecvPacketAckFailed,SubState.SuccessTimeoutPacketNotFound]
        const ibcTxLatestModel = this.taskCommonService.getIbcTxLatestModel()
        const dateNow = Math.floor(new Date().getTime() / 1000)
        await this.taskCommonService.changeIbcTxState(ibcTxLatestModel,dateNow,substate,false,[])
    }
}