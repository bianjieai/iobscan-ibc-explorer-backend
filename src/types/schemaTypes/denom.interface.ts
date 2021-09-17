import { Document } from 'mongoose';

export interface IDenomStruct {
    name?: string,
    denom_id?: string,
    json_schema?: string,
    creator?: string,
    tx_hash?: string,
    height?: number,
    time?: number,
    create_time?: number,
    last_block_height?: number,
    last_block_time?:number
}

export interface IDenom extends IDenomStruct, Document {
    
}

export interface IDenomMapStruct {
    name: string,
    denomId: string,
    jsonSchema: string,
    creator: string,
    height: number,
    txHash: string,
    createTime: number,
}

export interface IDenomMsgStruct {
    id?: string,
    name?: string,
    sender?: string,
    schema?: string
}