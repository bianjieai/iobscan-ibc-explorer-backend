import * as mongoose from 'mongoose';
import {IbcChainConfigType} from '../types/schemaTypes/ibc_chain_config.interface';

export const IbcChainConfigSchema = new mongoose.Schema({
    chain_id: String,
    icon: String,
    chain_name: String,
    lcd: String,
    lcd_api_path: Object,
    ibc_info: Object,
    is_manual: Boolean
});

IbcChainConfigSchema.index({chain_id: 1}, {unique: true});
IbcChainConfigSchema.statics = {
    async findCount(query): Promise<number> {
        return this.count(query);
    },

    // async aggregateFindChannels(): Promise<any> {
    //     return this.aggregate([{$group: {_id: '$ibc_info.paths.channel_id'}}]);
    // },

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

    // async findDcChain(
    //     query: {
    //         sc_chain_id: string,
    //         sc_port: string,
    //         sc_channel: string,
    //     },
    // ): Promise<{ _id: string; ibc_info: { chain_id: string }[] } | null> {
    //     // search dc_chain_config by sc_chain_id、sc_port、sc_channel
    //     const {sc_chain_id, sc_port, sc_channel} = query;
    //     return this.findOne(
    //         {
    //             chain_id: sc_chain_id,
    //             'ibc_info.paths.channel_id': sc_channel,
    //             'ibc_info.paths.port_id': sc_port,
    //         }
    //     );
    // },

    // async findScChain(
    //     query: {
    //         dc_chain_id: string,
    //         dc_port: string,
    //         dc_channel: string,
    //     },
    // ): Promise<{ _id: string; ibc_info: { chain_id: string }[] } | null> {
    //     // search dc_chain_config by sc_chain_id、sc_port、sc_channel
    //     const {dc_chain_id, dc_port, dc_channel} = query;
    //     return this.findOne(
    //         {
    //             'ibc_info.chain_id': dc_chain_id,
    //             'ibc_info.paths.counterparty.channel_id': dc_channel,
    //             'ibc_info.paths.counterparty.port_id': dc_port,
    //         }
    //     );
    // },

    async updateChain(chain: IbcChainConfigType) {
        const {chain_id} = chain;
        const options = {upsert: true, new: false, setDefaultsOnInsert: true};
        return this.findOneAndUpdate({chain_id}, chain, options);
    },
};
