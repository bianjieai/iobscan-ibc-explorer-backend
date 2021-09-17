import { IQueryBase } from '.';
import { Document } from 'mongoose';
import { ITxStruct } from './schemaTypes/tx.interface';

export interface ITxsQueryParams extends IQueryBase {
    type?: string,
    $or?: object[],
    status?: number,
    addrs?:any,
    time?: {
        $gte?: number,
        $lte?: number,
    }
}

export interface IListStruct {
    data?: any[],
    count?: number
}

export interface IBindTx {
    provider: string;
    respondTimes?: number,
    bindTime: string,
}

export interface ExternalIBindTx {
    provider: string;
    respondTimes?: number,
}

export interface IServiceName {
    serviceName: string;
    description?: string;
    bind: number;
    bindList?:IBindTx[],
}

export interface ExternalIServiceName {
    serviceName: string;
    description?: string;
    bind: number;
    bindList?:ExternalIBindTx[],
}

export interface IExFieldTx extends ITxStruct{
    ex_service_name?:string,
    ex_call_hash?: string,
    ex_request_context_id?: string,
    ex_consumer?: string,
    ex_bind?:number,
}

export interface IIdentity extends  IQueryBase{
    id?:string
}
