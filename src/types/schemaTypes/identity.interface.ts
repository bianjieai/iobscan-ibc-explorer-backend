import { IQueryBase } from '../index';
export interface IBaseIdentityStruct {
    identities_id:string
}

export interface IIdentityStruct extends IBaseIdentityStruct{
    credentials: string,
    owner: string,
    create_block_time: string,
    create_block_height: string,
    create_tx_hash: string,
    update_block_time: string,
    update_block_height: string,
    update_tx_hash: string,
    create_time: number,
    update_time: number
}
export interface IIdentityPubKeyStruct extends IBaseIdentityStruct {
    pubkey: object,
    hash: string,
    height: number,
    time: number,
    msg_index: number,
    pubkey_hash: string,
    certificate_hash: string,
    create_time: number,
}
export interface IIdentityCertificateStruct extends IBaseIdentityStruct{
    certificate_hash:string,
    certificate:string,
    hash: string,
    height: number,
    time: number,
    msg_index:number,
    create_time: number,
}

export interface IUpDateIdentityCredentials extends IBaseIdentityStruct {
    credentials?:string
    update_block_time: string,
    update_block_height: string,
    update_tx_hash: string,
    update_time?: number
}

export interface IIdentityPubKeyAndCertificateQuery extends IQueryBase{
    id:string
}
export interface IIdentityByAddressQuery extends IQueryBase{
    address:string
}
export interface IIdentityInfoQuery extends IQueryBase{
    id:string
}
export interface IIdentityInfoResponse {
    identities_id: string
    owner: string
    credentials: string
    create_block_height: number
    create_block_time: number
    create_tx_hash: string
}
