import { Injectable } from '@nestjs/common';
import { Model } from 'mongoose';
import { InjectModel } from '@nestjs/mongoose';
import { ITaskDispatch, ITaskDispatchStruct } from '../types/schemaTypes/task.dispatch.interface';
import { getIpAddress, getTimestamp } from '../util/util';
import { TaskEnum } from '../constant';
import { cfg } from '../config/config';
import { Logger } from '../logger';
import { IRandomKey } from '../types';
import { DispatchFaultTolerance } from '../types/task.interface';

@Injectable()
export class TaskDispatchService {

    constructor(@InjectModel('TaskDispatch') private taskDispatchModel: Model<ITaskDispatch>) {
        this.updateHeartbeatUpdateTime = this.updateHeartbeatUpdateTime.bind(this);
    }

    async needDoTask(name: TaskEnum): Promise<boolean> {
        const task: ITaskDispatchStruct | null = await (this.taskDispatchModel as any).findOneByName(name);
        if (task) {
            if (task.is_locked) {
                Logger.log(`${name}: previous task is executing, the next should not be executed!`);
                return false;
            } else {
                const updated: boolean = await this.lock(name);
                if (updated) {
                    return true;
                } else {
                    return false;
                }
            }
        } else {
            //it should be register if there is no this type of task;
            const registered = await this.registerTask(name);
            Logger.log(`${name}: register successfully? ${registered}`);
            if (registered) {
                const updated: boolean = await this.lock(name);
                if (updated) {
                    return true;
                } else {
                    return false;
                }
            } else {
                Logger.warn(`${name}: task has not been registered, but it couldn't register successfully!`);
                return false;

            }
        }
    }

    async registerTask(name: TaskEnum): Promise<ITaskDispatchStruct | null> {
        const task: ITaskDispatchStruct = {
            name,
            is_locked: false,
            device_ip: getIpAddress(),
            create_time: getTimestamp(),
            task_begin_time: 0,
            task_end_time: 0,
            heartbeat_update_time: getTimestamp(),
        };
        return await (this.taskDispatchModel as any).createOne(task);
    }

    private async lock(name: TaskEnum): Promise<boolean> {
        return await (this.taskDispatchModel as any).lock(name);
    }

    async unlock(name: TaskEnum, randomKey?: IRandomKey): Promise<boolean> {
        return await (this.taskDispatchModel as any).unlock(name, randomKey);
    }

    async taskDispatchFaultTolerance(cb: DispatchFaultTolerance): Promise<void> {
        const taskList: ITaskDispatchStruct[] = await (this.taskDispatchModel as any).findAllLocked();
        if (taskList && taskList.length > 0) {
            for(let task of taskList){
                //对比当前时间跟上次心跳更新时间的差值与 心率, 当大于两个心率周期的时候, 认为上一个执行task的实例发生故障
                if ((getTimestamp() - task.heartbeat_update_time) >= (cfg.taskCfg.interval.heartbeatRate * 2)/1000) {
                    Logger.error(`${(task as any).name}: task executed breakdown, and reed to be released the lock`);
                    await this.releaseLockByName((task as any).name);
                    //出故障后, 释放锁, 也需要将出故障的timer清除;
                    cb((task as any).name);
                }
            }
        }
    }

    private async releaseLockByName(name: TaskEnum): Promise<ITaskDispatchStruct | null> {
        return await (this.taskDispatchModel as any).releaseLockByName(name);
    }

    public async updateHeartbeatUpdateTime(name: TaskEnum, randomKey?: IRandomKey): Promise<void> {
        await (this.taskDispatchModel as any).updateHeartbeatUpdateTime(name, randomKey);
    }

    public async deleteOneByName(name: TaskEnum): Promise<void> {
        await (this.taskDispatchModel as any).deleteOneByName(name);
    }

}

