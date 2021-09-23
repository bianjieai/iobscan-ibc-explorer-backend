import * as mongoose from 'mongoose';
import { IbcChainType } from '../types/schemaTypes/ibc_chain.interface';

export const IbcChainSchema = new mongoose.Schema({
  chain_id: String,
  icon: String,
  chain_name: String,
  lcd: String,
  ibc_info: Object,
});

IbcChainSchema.index({ chain_id: 1 }, { unique: true });

IbcChainSchema.statics = {
  // 查
  async findCount(query): Promise<Number> {
    return this.count(query);
  },

  async aggregateFindChannels() {
    return this.aggregate([{ $group: { _id: '$ibc_info.paths.channel_id' } }]);
  },

  async findAll(): Promise<IbcChainType[]> {
    return this.find({});
  },

  async findList(): Promise<IbcChainType[]> {
    return this.find()
  },

  async findDcChainId(
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
      },
      { 'ibc_info.chain_id': 1 },
    );
  },

  // 改
  async updateChain(chain: IbcChainType) {
    const { chain_id } = chain;
    const options = { upsert: true, new: false, setDefaultsOnInsert: true };
    return this.findOneAndUpdate({ chain_id }, chain, options);
  },
};
