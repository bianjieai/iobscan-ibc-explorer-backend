import * as mongoose from 'mongoose';
import { getIpAddress, getTimestamp } from '../util/util';
import { ITaskDispatchStruct } from '../types/schemaTypes/task.dispatch.interface';
import { TaskEnum } from 'src/constant';
import { Logger } from '../logger';
import { IRandomKey } from '../types';
import { taskLoggerHelper } from '../helper/task.log.helper';
import moment from 'moment';

export const TaskDispatchSchema = new mongoose.Schema({
    name: { type: String, unique: true },
    is_locked: Boolean,
    device_ip: String,
    create_time: Number,
    task_begin_time: Number,
    task_end_time: Number,
    heartbeat_update_time: Number,
},{versionKey: false});

TaskDispatchSchema.statics = {
    async findOneByName(name: TaskEnum): Promise<ITaskDispatchStruct | null> {
        return await this.findOne({ name }).exec();
    },

    async createOne(t: ITaskDispatchStruct): Promise<ITaskDispatchStruct | null> {
        return new this(t).save();
    },
    async lock(name: TaskEnum): Promise<boolean> {
        return new Promise(async (res)=>{
            return await this.update({ name, is_locked: false }, {
                // condition: is_locked: false, those server whose query's is_locked is true should not to be updated;
                is_locked: true,
                task_begin_time: getTimestamp(),
                device_ip: getIpAddress(),
            }, null, (error,effect)=>{
                if(error) {
                    res(false);
                    return;
                }
                if(effect && effect.nModified === 1){
                    res(true);
                    Logger.log(`From task.dispatch.schema ${name} task begin time: ${new Date().getTime()}`);
                }else {
                    res(false);
                }
            }).exec();
        });


    },

    async unlock(name: TaskEnum, randomKey?: IRandomKey): Promise<boolean> {
        return new Promise(async (res)=>{
            //console.log('before updated unlock:',await this.find({ name}))
            return await this.update({ name, is_locked: true }, {
                is_locked: false,
                task_end_time: getTimestamp(),
            }, null, (error,effect)=>{
                if(error) {
                    taskLoggerHelper(`${name}: unlock error`, randomKey);
                    res(false);
                    return;
                }

                //TODO(lvshenchao) it seems like that there is a bug in mongoose, the params of nModified returned 0 while record was be updated
                // console.log('updated unlock:',effect);
                //console.log('updated unlock:',await this.find({ name}))
                if(effect && effect.nModified === 1){
                    taskLoggerHelper(`${name}: unlock successful, From task.dispatch.schema ${name} task end time: ${new Date().getTime()}`, randomKey);
                    res(true);
                }else {
                    taskLoggerHelper(`${name}: unlock failed`, randomKey);
                    res(false);
                }
            }).exec();
        })
    },

    async releaseLockByName(name: TaskEnum): Promise<ITaskDispatchStruct | null> {
        return await this.updateOne({
            name,
            is_locked:true
        }, {
            is_locked: false,
        }).exec();
    },

    async findAllLocked(): Promise<ITaskDispatchStruct[]> {
        return await this.find({is_locked:true}).exec();
    },

    async updateHeartbeatUpdateTime(name: TaskEnum, randomKey?: IRandomKey): Promise<ITaskDispatchStruct | null> {
        if(randomKey){
            taskLoggerHelper(`${name}: update hearbeat time`, randomKey);
        }else{
            //定时任务打印的日志, 不需要对step递增
            Logger.log(`${name}: update hearbeat time: ${moment(new Date().getTime()).format('YYYY-MM-DD HH:mm:ss')}`);
        }

        return await this.updateOne({ name, is_locked: true }, {
            heartbeat_update_time: getTimestamp(),
        }).exec();
    },
    async deleteOneByName(name: TaskEnum) {
        await this.deleteOne({name: name})
    }
};