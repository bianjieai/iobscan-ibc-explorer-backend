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
import { IbcTxTaskService } from './ibc_tx.task.service';
import { IbcChainConfigTaskService } from './ibc_chain_config.task.service';
import { IbcStatisticsTaskService } from './ibc_statistics.task.service';
@Injectable()
export class TasksService {
  constructor(
    private readonly taskDispatchService: TaskDispatchService,
    private readonly ibcChainConfigTaskService: IbcChainConfigTaskService,
    private readonly ibcTxTaskService: IbcTxTaskService,
    private readonly ibcStatisticsTaskService: IbcStatisticsTaskService,
  ) {
    this[`${TaskEnum.denom}_timer`] = null;
  }

  // chainConfig
  @Cron(cfg.taskCfg.executeTime.chain, {
    name: TaskEnum.chain,
  })
  async syncChain() {
    this.handleDoTask(TaskEnum.chain, this.ibcChainConfigTaskService.doTask);
  }

  // ex_ibc_tx
  @Cron('*/5 * * * * *', {
    name: TaskEnum.tx,
  })
  async syncTx() {
    this.handleDoTask(TaskEnum.tx, this.ibcTxTaskService.doTask);
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

  async handleDoTask(taskName: TaskEnum, doTask: TaskCallback) {
    if (
      cfg &&
      cfg.taskCfg &&
      cfg.taskCfg.CRON_JOBS &&
      cfg.taskCfg.CRON_JOBS.indexOf(taskName) === -1
    ) {
      return;
    }
    // 只执行一次删除定时任务
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
      //一个定时任务的完整周期必须严格按照:
      //上锁 ---> 更新心率时间 ---> 启动定时更新心率时间任务 ---> 执行定时任务 ---> 释放锁 ---> 清除心率定时任务
      //否则如果在执行定时任务未结束之前将锁打开, 那么有可能后面的实例会重新执行同样的任务
      //为了清晰的看出完整的周期执行顺序, 为每一次的定时任务新增一个唯一key(大致唯一, 只要跟最近的定时任务不重复即可),并标注执行步骤
      let randomKey: IRandomKey = {
        key: String(Math.random()),
        step: 0,
      };

      try {
        //因为一般情况下定时任务执行时间要小于心跳率, 为防止heartbeat_update_time一直不被更新,
        //所以在任务开始之前先更新一下heartbeat_update_time;
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
