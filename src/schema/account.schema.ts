import * as mongoose from 'mongoose';
import {
    IAccountStruct,
    ITokenTotal
} from '../types/schemaTypes/account.interface';
import {IListStruct} from "../types";

export const AccountSchema = new mongoose.Schema({
    address: String,
    account_total: Number,
    total: Object,
    balance: Object,
    delegation: Object,
    unbonding_delegation: Object,
    rewards: Object,
    create_time: Number,
    update_time: Number,
    handled_block_height: Number,
})
AccountSchema.index({ address: 1 }, { unique: true });
AccountSchema.index({ account_total: -1 }, { background: true });
AccountSchema.index({ handled_block_height: -1 }, { background: true });

AccountSchema.statics = {
    async queryHandledBlockHeight(): Promise<IAccountStruct>{
        return await this.find({},{handled_block_height: 1,address: 1}).sort({handled_block_height: -1}).limit(1);
    },
    async insertManyAccount(AccountList): Promise<IAccountStruct[]>{
        return await this.insertMany(AccountList, { ordered: false })
    },
    async updateAccount(account:IAccountStruct) {
        let { address } = account
        const options = {upsert: true, new: false, setDefaultsOnInsert: true}
        await this.findOneAndUpdate({address}, account, options)
    },
    async queryAllAddress(): Promise<IAccountStruct>{
        return await this.find({}, {address:1,_id:0});
    },
    async queryAccountsLimit(): Promise<IAccountStruct[]> {
        return await this.find({account_total : {$gt : 0}}, {address:1,total:1,_id:0,update_time:1}).sort({account_total: -1}).limit(100);
    },
    async queryTokenTotal(): Promise<ITokenTotal> {
        let data = await this.aggregate( [
            { $group: { _id: null, "account_totals" : { $sum: "$account_total" } } }
        ])
        return data[0]
    },
    async queryAccountsTotalLimit(): Promise<IAccountStruct[]> {
        return await this.find({}, {total:1,_id:0}).sort({account_total: -1}).limit(1000);
    },
}
