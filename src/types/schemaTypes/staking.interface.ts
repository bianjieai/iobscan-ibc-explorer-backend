import {IQueryBase} from "../index";

export interface IStakingValidator {
    operator_address:string,
    consensus_pubkey:string,
    jailed:boolean,
    status:number,
    tokens:string,
    delegator_shares:string,
    description:object,
    bond_height:string,
    unbonding_height:string,
    unbonding_time:string,
    commission:object,
    uptime:number,
    self_bond:string,
    delegator_num:number,
    proposer_addr:string,
    voting_power:number,
    min_self_delegation: number,
    icons: string,
    start_height: string,
    index_offset: string,
    jailed_until: string,
    tombstoned: boolean,
    missed_blocks_counter: string,
    create_time: number,
    update_time: number,
}
export interface IStakingValidatorLcdMap {
    operator_address:string,
    consensus_pubkey:string,
    status:number,
    tokens:string,
    delegator_shares:string,
    description:object,
    unbonding_time:number,
    commission:object,
    min_self_delegation:string
}
//TODO 设置DB的map 的interface
export interface IStakingValidatorDbMap {

}
export interface IQueryValidatorByStatus extends IQueryBase{
    status:string
}
export interface IDetailByValidatorAddress{
    address:string
}
export interface IStakingValidatorFromLcd {
    operator_address:string,
    consensus_pubkey:string,
    status:number,
    tokens:string,
    description:object | null,
    unbonding_time:string,
    commission:object,
    min_self_delegation:string,
}
export interface IStakingValidatorBlock {
    ivaAddr:string,
    monikerM: string,
    isBlack:boolean
}
