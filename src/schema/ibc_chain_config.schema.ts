import * as mongoose from 'mongoose';
import {IbcChainConfigType} from '../types/schemaTypes/ibc_chain_config.interface';

export const IbcChainConfigSchema = new mongoose.Schema({
    chain_id: String,
    icon: String,
    chain_name: String,
    lcd: String,
    lcd_api_path: Object,
    ibc_info: Object,
    ibc_info_hash_lcd: String,
    ibc_info_hash_caculate: String,
    is_manual: Boolean
});

IbcChainConfigSchema.index({chain_id: 1}, {unique: true});
IbcChainConfigSchema.statics = {
    async findCount(query): Promise<number> {
        return this.count(query);
    },

    async findAll(): Promise<IbcChainConfigType[]> {
        return this.find({}).sort({'chain_name': 1});
    },

    async findList(): Promise<IbcChainConfigType[]> {
        return this.find().select({
            "_id": 0,
            'chain_id': 1,
            'chain_name': 1,
            'icon': 1
        }).collation({locale: 'en_US'}).sort({'chain_name': 1});
    },



    async updateChain(chain: IbcChainConfigType) {
        const {chain_id} = chain;
        const options = {upsert: true, new: false, setDefaultsOnInsert: true};
        return this.findOneAndUpdate({chain_id}, chain, options);
    },
    async updateChainCfgWithSession(chain: IbcChainConfigType,session) {
        const {chain_id} = chain;
        const options = {session,upsert: true, new: false, setDefaultsOnInsert: true};
        return this.findOneAndUpdate({chain_id}, chain, options);
    },
};
