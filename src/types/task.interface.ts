import { IRandomKey } from './index';
import { TaskEnum } from '../constant';

export interface TaskCallback {
    (taskName?: TaskEnum, randomKey?: IRandomKey): Promise<void>;
}

export interface DispatchFaultTolerance {
    (taskName: TaskEnum): void;
}

