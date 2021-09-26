import { Document } from 'mongoose';

export interface ITaskDispatchStruct {
    name?: string,
    is_locked?: boolean,
    device_ip?: string,
    create_time?: number,
    task_begin_time?: number,
    task_end_time?: number,
    heartbeat_update_time?: number,
}

export interface ITaskDispatch extends ITaskDispatchStruct, Document {
}