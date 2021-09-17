import { IRandomKey } from './index';
import { TaskEnum } from '../constant';

export interface TaskCallback {
    (taskName?: TaskEnum, randomKey?: IRandomKey): Promise<void>;
}

export interface DispatchFaultTolerance {
    (taskName: TaskEnum): void;
}

export interface ILcdNftStruct {
    id: string;
    name: string;
    owner: string;
    data: string;
    uri?: string;
    hash?: string;
    denom_id?: string;
    denom_name?: string;
    time?: number
}