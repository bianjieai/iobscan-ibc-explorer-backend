import { Injectable } from '@nestjs/common';
import { Cron, SchedulerRegistry } from '@nestjs/schedule';
import { TaskDispatchService } from '../service/task.dispatch.service';
import { TaskEnum } from '../constant';
import { getIpAddress } from '../util/util';
import { cfg } from '../config/config';
import { TaskCallback } from '../types/task.interface';
import { Logger } from '../logger';
import { IRandomKey } from '../types';
import { taskLoggerHelper } from '../helper/task.log.helper';
import { IbcChainConfigTaskService } from './ibc_chain_config.task.service';
import { IbcStatisticsTaskService } from './ibc_statistics.task.service';
import {IbcSyncTransferTxTaskService} from "./ibc_sync_transfer_tx_task.service";
import {IbcUpdateProcessingTxTaskService} from "./ibc_update_processing_tx_task.service";
import {IbcUpdateSubStateTxTaskService} from "./ibc_update_substate_tx_task.service";
import {IbcTxDataUpdateTaskService} from "./ibc_tx_data_update_task.service";
import {IbcTxLatestMigrateTaskService} from "./ibc_tx_latest_migrate_task.service";
import {IbcDenomCaculateTaskService} from "./ibc_denom_caculate.task.service";
import {IbcDenomUpdateTaskService} from "./ibc_denom_update.task.service";
import {IbcMonitorService} from "../monitor/ibc_monitor.service";
@Injectable()
export class TasksService {
  constructor(
    private readonly taskDispatchService: TaskDispatchService,
    private readonly ibcChainConfigTaskService: IbcChainConfigTaskService,
    private readonly ibcStatisticsTaskService: IbcStatisticsTaskService,
    private readonly ibcSyncTransferTxTaskService : IbcSyncTransferTxTaskService,
    private readonly ibcUpdateProcessingTxService : IbcUpdateProcessingTxTaskService,
    private readonly ibcUpdateSubstateTxService: IbcUpdateSubStateTxTaskService,
    private readonly ibcTxDataUpdateTaskService: IbcTxDataUpdateTaskService,
    private readonly ibcTxLatestMigrateTaskService: IbcTxLatestMigrateTaskService,
    private readonly ibcDenomCaculateTaskService: IbcDenomCaculateTaskService,
    private readonly ibcDenomUpdateTaskService: IbcDenomUpdateTaskService,
    private readonly ibMonitorService: IbcMonitorService
  ) {
    // this[`${TaskEnum.denom}_timer`] = null;
  }

  // chainConfig
  @Cron(cfg.taskCfg.executeTime.chain, {
    name: TaskEnum.chain,
  })
  // @Cron('*/10 * * * * *')
  async syncChain() {
    this.handleDoTask(TaskEnum.chain, this.ibcChainConfigTaskService.doTask);
  }

  @Cron(cfg.taskCfg.executeTime.ibcMonitorCronjob, {
    name: TaskEnum.monitor,
  })
  // @Cron('*/10 * * * * *')
  async syncMonitor() {
    this.handleDoTask(TaskEnum.monitor, this.ibMonitorService.doTask);
  }

  //????????????
  @Cron(cfg.taskCfg.executeTime.transferTx, {
    name: TaskEnum.transferTx,
  })
  // @Cron('*/1 * * * * *')
  async syncTransferTx() {
    this.handleDoTask(TaskEnum.transferTx, this.ibcSyncTransferTxTaskService.doTask);
  }


  @Cron(cfg.taskCfg.executeTime.updateProcessingTx, {
     name: TaskEnum.updateProcessingTx,
   })
  // @Cron('*/15 * * * * *')
  async updateProcessingTx() {
    this.handleDoTask(TaskEnum.updateProcessingTx, this.ibcUpdateProcessingTxService.doTask);
  }


  @Cron(cfg.taskCfg.executeTime.updateSubStateTx, {
       name: TaskEnum.updateSubStateTx,
     })
  // @Cron('*/15 * * * * *')
  async upSubstateTx() {
    this.handleDoTask(TaskEnum.updateSubStateTx, this.ibcUpdateSubstateTxService.doTask);
  }


  // ex_ibc_statistics
  @Cron(cfg.taskCfg.executeTime.statistics, {
    name: TaskEnum.statistics,
  })
  async syncStatistics() {
    this.handleDoTask(
      TaskEnum.statistics,
      this.ibcStatisticsTaskService.doTask,
    );
  }

  @Cron(cfg.taskCfg.executeTime.faultTolerance)
  //@Cron('18 * * * * *')
  async taskDispatchFaultTolerance() {
    this.taskDispatchService.taskDispatchFaultTolerance((name: TaskEnum) => {
      if (this[`${name}_timer`]) {
        clearInterval(this[`${name}_timer`]);
        this[`${name}_timer`] = null;
      }
    });
  }

  @Cron(cfg.taskCfg.executeTime.ibcTxUpdateCronjob)
  // @Cron('* */2 * * * *')
  async ibcTxUpdateCronjob() {
    this.handleDoTask(TaskEnum.ibcTxUpdateCronJob, this.ibcTxDataUpdateTaskService.doTask);
  }

  @Cron(cfg.taskCfg.executeTime.ibcTxLatestMigrate)
  // @Cron('*/30 * * * * *')
  async ibcTxMigrateCronjob() {
    this.handleDoTask(TaskEnum.ibcTxMigrateCronJob, this.ibcTxLatestMigrateTaskService.doTask);
  }

  @Cron(cfg.taskCfg.executeTime.ibcDenomCaculateCronjob)
  // @Cron('*/30 * * * * *')
  async ibcDenomCaculateCronjob() {
    this.handleDoTask(TaskEnum.ibcDenomCaculateCronJob, this.ibcDenomCaculateTaskService.doTask);
  }

  @Cron(cfg.taskCfg.executeTime.ibcDenomUpdateCronjob)
  // @Cron('*/30 * * * * *')
  async ibcDenomUpdateCronjob() {
    this.handleDoTask(TaskEnum.ibcDenomUpdateCronjob, this.ibcDenomUpdateTaskService.doTask);
  }

  async handleDoTask(taskName: TaskEnum, doTask: TaskCallback) {
    if (
      cfg &&
      cfg.taskCfg &&
      cfg.taskCfg.CRON_JOBS &&
      cfg.taskCfg.CRON_JOBS.indexOf(taskName) === -1
    ) {
      return;
    }
    // ?????????????????????????????????
    // if (this['once'] && cfg.taskCfg.DELETE_CRON_JOBS && cfg.taskCfg.DELETE_CRON_JOBS.length) {
    //     cfg.taskCfg.DELETE_CRON_JOBS.forEach(async item => {
    //         this.schedulerRegistry.deleteCronJob(item)
    //         await this.taskDispatchService.deleteOneByName(item)
    //     })
    //     this['once'] = false
    // }
    const needDoTask: boolean = await this.taskDispatchService.needDoTask(
      taskName,
    );
    Logger.log(
      `the ip ${getIpAddress()} (process pid is ${
        process.pid
      }) should do task ${taskName}? ${needDoTask}`,
    );
    if (needDoTask) {
      //???????????????????????????????????????????????????:
      //?????? ---> ?????????????????? ---> ???????????????????????????????????? ---> ?????????????????? ---> ????????? ---> ????????????????????????
      //????????????????????????????????????????????????????????????, ????????????????????????????????????????????????????????????
      //????????????????????????????????????????????????, ?????????????????????????????????????????????key(????????????, ?????????????????????????????????????????????),?????????????????????
      let randomKey: IRandomKey = {
        key: String(Math.random()),
        step: 0,
      };

      try {
        //???????????????????????????????????????????????????????????????, ?????????heartbeat_update_time??????????????????,
        //??????????????????????????????????????????heartbeat_update_time;
        await this.updateHeartbeatUpdateTime(taskName, randomKey);
        const beginTime: number = new Date().getTime();
        this[`${taskName}_timer`] = setInterval(() => {
          this.updateHeartbeatUpdateTime(taskName);
        }, cfg.taskCfg.interval.heartbeatRate);
        await doTask(taskName, randomKey);
        //weather task is completed successfully, lock need to be released;
        const unlock: boolean = await this.taskDispatchService.unlock(
          taskName,
          randomKey,
        );
        taskLoggerHelper(
          `${taskName}: (ip: ${getIpAddress()}, pid: ${
            process.pid
          }) has released the lock? ${unlock}`,
          randomKey,
        );
        if (this[`${taskName}_timer`]) {
          clearInterval(this[`${taskName}_timer`]);
          this[`${taskName}_timer`] = null;
          taskLoggerHelper(
            `${taskName}: timer has been cleared out`,
            randomKey,
          );
        }

        taskLoggerHelper(
          `${taskName}: current task executes end, took ${new Date().getTime() -
            beginTime}ms`,
          randomKey,
        );
      } catch (e) {
        Logger.error(
          `${taskName}: task executes error, should release lock`,
          e,
        );
        await this.taskDispatchService.unlock(taskName, randomKey);
        if (this[`${taskName}_timer`]) {
          clearInterval(this[`${taskName}_timer`]);
          this[`${taskName}_timer`] = null;
        }
      }
    }
  }

  async updateHeartbeatUpdateTime(
    name: TaskEnum,
    randomKey?: IRandomKey,
  ): Promise<void> {
    await this.taskDispatchService.updateHeartbeatUpdateTime(name, randomKey);
  }
}
