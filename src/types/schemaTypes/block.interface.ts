import { Document } from 'mongoose';

export interface IBlockStruct {
    height?:number,
    hash?:string,
    txn?:number,
    time?:string,
    proposer?: string,
    total_validator_num?: number;
    total_voting_power?: number;
    precommit_voting_power?: number;
    precommit_validator_num?: number;
    proposer_moniker?: string;
    proposer_addr?: string;
}

export interface IBlock extends IBlockStruct, Document {
    
}