import * as mongoose from 'mongoose';
import { IbcChainConfigType } from '../types/schemaTypes/ibc_chain_config.interface';

export const IbcChainConfigSchema = new mongoose.Schema({
  chain_id: String,
  icon: String,
  chain_name: String,
  lcd: String,
  ibc_info: Object,
});

IbcChainConfigSchema.index({ chain_id: 1 }, { unique: true });
IbcChainConfigSchema.statics = {
  async findCount(query): Promise<Number> {
    return this.count(query);
  },

  async aggregateFindChannels(): Promise<any> {
    return this.aggregate([{ $group: { _id: '$ibc_info.paths.channel_id' } }]);
  },

  async findAll(): Promise<IbcChainConfigType[]> {
    return this.find({})
  },

  async findList(): Promise<IbcChainConfigType[]> {
    return this.find().collation( { locale: 'en_US' } ).sort({'chain_id': 1});
  },

  async findDcChain(
    query,
  ): Promise<{ _id: string; ibc_info: { chain_id: string }[] } | null> {
    // search dc_chain_id by sc_chain_id、sc_port、sc_channel、dc_port、dc_channel
    const { sc_chain_id, sc_port, sc_channel, dc_port, dc_channel } = query;
    return this.findOne(
      {
        chain_id: sc_chain_id,
        'ibc_info.paths.channel_id': sc_channel,
        'ibc_info.paths.port_id': sc_port,
        'ibc_info.paths.counterparty.channel_id': dc_channel,
        'ibc_info.paths.counterparty.port_id': dc_port,
      }
    );
  },

  async updateChain(chain: IbcChainConfigType) {
    const { chain_id } = chain;
    const options = { upsert: true, new: false, setDefaultsOnInsert: true };
    return this.findOneAndUpdate({ chain_id }, chain, options);
  },
};
