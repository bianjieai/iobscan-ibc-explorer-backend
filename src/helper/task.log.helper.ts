import {Logger} from '../logger';
import { IRandomKey } from '../types';

export function taskLoggerHelper(log: string, randomKey?: IRandomKey): void{
    if(randomKey){
        let str: string = `${log}, random key: ${randomKey.key}_${randomKey.step++}`;
        Logger.log(str);
    }
}