import * as mongoose from 'mongoose';
import {Logger} from "../logger";
import {IListStruct} from "../types";
import {
    IQueryValidatorByStatus,
    IStakingValidator,
    IStakingValidatorBlock,
    IDetailByValidatorAddress
} from "../types/schemaTypes/staking.interface";
import {activeValidatorLabel, candidateValidatorLabel, jailedValidatorLabel, ValidatorStatus} from "../constant";

export const StakingValidatorSchema = new mongoose.Schema({
    operator_address: String,
    consensus_pubkey: String,
    jailed: Boolean,
    status: Number,
    tokens: Number,
    delegator_shares: String,
    description: Object,
    bond_height: String,
    unbonding_height: String,
    unbonding_time: String,
    commission: Object,
    uptime: Number,
    self_bond: Object,
    delegator_num: Number,
    proposer_addr: String,
    voting_power: Number,
    min_self_delegation: Number,
    icon: String,
    start_height: String,
    index_offset: String,
    jailed_until: String,
    tombstoned: Boolean,
    missed_blocks_counter: String,
    create_time: Number,
    update_time: Number,
    is_black: Boolean,
    moniker_m: String
})
StakingValidatorSchema.index({operator_address: 1}, {unique: true})
StakingValidatorSchema.index({proposer_addr: 1})	
StakingValidatorSchema.index({jailed: 1, status: 1,tokens:1})

StakingValidatorSchema.statics = {

    async queryAllValidators() {
        return await this.find({})
    },

    async findValidatorByPropopserAddr(PropopserAddr: string): Promise<IStakingValidator>{
        return await this.find({proposer_addr:PropopserAddr});
    },

    async insertValidator(insertValidator:IStakingValidator) {
        let {operator_address} = insertValidator
        //设置 options 查询不到就插入操作
        const options = {upsert: true, new: false, setDefaultsOnInsert: true}
        await this.findOneAndUpdate({operator_address}, insertValidator, options)
    },

    async deleteValidator(deleteValidator: IStakingValidator) {
        const {operator_address} = deleteValidator
        await this.deleteOne({operator_address})
    },

    async queryAllValCommission(query): Promise<IListStruct> {
        const result: IListStruct = {}
        if (query.useCount && query.useCount == true) {
            result.count = await this.find({}).countDocuments();
        }
        result.data = await this.find({}).select({'_id': 0, '__v': 0})
        return result
    },
    async queryValidatorsByStatus(query: IQueryValidatorByStatus): Promise<IListStruct> {
        const queryParameters: any = {};
        const result: IListStruct = {}

        if(query.status === jailedValidatorLabel){
            queryParameters.jailed = true
        }else if(query.status === activeValidatorLabel){
            queryParameters.jailed = false
            queryParameters.status = ValidatorStatus['bonded']
        }else if(query.status === candidateValidatorLabel){
            queryParameters.jailed = false
            queryParameters.status ={'$in':[ ValidatorStatus['Unbonded'], ValidatorStatus['Unbonding']]}
        }
        if (query.useCount && query.useCount == true) {
            result.count = await this.find(queryParameters).countDocuments();
        }
        // console.log('查询条件',queryParameters)
        result.data = await this.find(queryParameters)
            .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
            .limit(Number(query.pageSize)).sort({ tokens: -1 });
        return result
    },

    async queryDetailByValidator(operator_address: string): Promise<any> {
        let queryParameters: any = {};
        queryParameters.operator_address = operator_address
        let result = await this.findOne(queryParameters)
        return result
    },

    async queryActiveValCount(): Promise<any> {
        let count: number;
        return count = await this.find({'status':ValidatorStatus['bonded'],'jailed':false}).countDocuments();
    },

    async queryActiveVal(): Promise<any> {
        return await this.find({'status':ValidatorStatus['bonded'],'jailed':false});
    },

    async updateBlcakValidator(updateValidator: IStakingValidatorBlock) {
        let { ivaAddr, monikerM,isBlack } = updateValidator
        return await this.updateOne({ operator_address: ivaAddr }, { $set: {is_black: isBlack,moniker_m: monikerM}})
    },

}

