import { IQueryBase } from '../index';

export interface ITxsQuery extends IQueryBase {
	type?:string,
	status?:string,
    address?:string,
	beginTime?:string,
  	endTime?:string,
}

export interface ITxsWithHeightQuery extends IQueryBase {
	height?:string,
}

export interface ITxsWithAddressQuery extends IQueryBase {
	address?:string,
    type?:string,
    status?:string,
}

export interface ITxsWithContextIdQuery extends IQueryBase {
    contextId?:string,
    type?:string,
    status?:string,
}

export interface ITxsWithNftQuery extends IQueryBase {
	denomId?:string,
	tokenId?:string,
}

export interface ITxsWithServiceNameQuery extends IQueryBase {
	serviceName?:string,
}

export interface ITxStruct {
    time:number,
    height:number,
    tx_hash:string,
    memo:string,
    status:number,
    log:string,
    complex_msg:boolean,
    type:string,
    from:string,
    to:string,
    coins:object[],
    signer:string,
    events:object[],
    msgs:object[],
    signers: string[],
    addrs: string[],
    fee: object[]
}

export interface callServiceStruct extends ITxStruct {
    respond?:object[],
}

export interface bindServiceStruct extends ITxStruct {
    respond_times?:number,
}

export interface ITxStructMsgs {
    events:object[],
    msgs:object[],
}

export interface ITxStructHash {
    tx_hash:object[],
    msgs:object[],
}

export interface IExFieldQuery {
    requestContextId?: string;
    consumer?: string;
    serviceName?: string;
    callHash?: string;
    hash: string;
    bind?: number;
}
export interface ITXWithIdentity extends IQueryBase {
    search?: string
}
export interface IIdentityTx extends IQueryBase {
    id?: string
}

export interface ITxsWithAssetQuery {
	type?:string
    symbol?:string
    pageNum?:number
    pageSize?:number
    useCount?:boolean
}

export interface ITxSubmitProposal  {
    tx_hash: string
    proposer: string
    initial_deposit: object[]
    content: string|object
} 

export interface ITxVoteProposal  {
    _id: string
    msg: object
    count: number
}

export interface ITxVoteProposalAll  {
    height: number
    tx_hash: string
    msgs: object
}

export interface ITxVoteALL {
    msgs: object
    tx_hash: string
}