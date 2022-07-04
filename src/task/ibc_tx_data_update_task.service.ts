import {Injectable, Logger} from '@nestjs/common';
import {RecordLimit, SubState, TaskEnum} from "../constant";
import {IbcTxHandler} from "../util/IbcTxHandler";
import {cfg} from "../config/config";

@Injectable()
export class IbcTxDataUpdateTaskService {
    constructor(private readonly taskCommonService: IbcTxHandler) {
        this.doTask = this.doTask.bind(this);
    }

    async doTask(taskName?: TaskEnum): Promise<void> {
        await this.handleUpdateIbcTx()
    }

    async handleUpdateIbcTx() {
        let pageNum = 1, handIbcTxs = []
        const defaultSubstate = 0
        const substate = [defaultSubstate, SubState.SuccessRecvPacketNotFound, SubState.RecvPacketAckFailed, SubState.SuccessTimeoutPacketNotFound]
        const ibcTxModel = this.taskCommonService.getIbcTxModel()

        const chainsCfg = await this.taskCommonService.getChainsCfg()
        // 依次遍历每条链去处理ibc_tx表中历史状态的数据
        for (const chainCfg of chainsCfg) {

            //judge  handIbcTxs size for handle batch limit
            while (handIbcTxs?.length < cfg.serverCfg.updateIbcTxBatchLimit) {
                const ibcTxs = await this.taskCommonService.getProcessingTxsByPage(chainCfg.chain_id, substate, pageNum)
                if (ibcTxs.length) {
                    handIbcTxs.push(...ibcTxs)
                }
                // break when finish collect all the ibc tx.
                if (ibcTxs?.length < RecordLimit) {
                    break;
                }
                pageNum++
            }
            const dateNow = Math.floor(new Date().getTime() / 1000)
            await this.taskCommonService.changeIbcTxState(ibcTxModel, dateNow, substate, true, handIbcTxs)
            Logger.log("finish update chain "+chainCfg.chain_id+ " ibc_tx processing data")
        }

    }
}