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
    return await this.count(query);
  },

  async aggregateFindChannels() {
    return await this.aggregate([
      { $group: { _id: '$ibc_info.paths.channel_id' } },
    ]);
  },

  async findAll(): Promise<IbcChainType[]> {
    return await this.find({});
  },

  async findList(
    pageNum: number,
    pageSize: number,
    chain_name?: String,
  ): Promise<IbcChainType[]> {
    const result = await this.find(
      chain_name ? { chain_name: { $regex: chain_name } } : undefined,
    )
      .skip((Number(pageNum) - 1) * Number(pageSize))
      .limit(Number(pageSize));
    return result;
  },

  async findDcChainId(query): Promise<String> {
    // search dc_chain_id by sc_chain_id、sc_port、sc_channel、dc_port、dc_channel
    const { sc_chain_id, sc_port, sc_channel, dc_port, dc_channel } = query;
    const result = await this.findOne(
      {
        chain_id: sc_chain_id,
        'ibc_info.paths.channel_id': sc_channel,
        'ibc_info.paths.port_id': sc_port,
        'ibc_info.paths.counterparty.channel_id': dc_channel,
        'ibc_info.paths.counterparty.port_id': dc_port,
      },
      { 'ibc_info.chain_id': 1 },
    );
    if (result && result.ibc_info && result.ibc_info.length) {
      return result.ibc_info[0]['chain_id'];
    } else {
      return '';
    }
  },

  // 改
  async updateChain(chain: IbcChainType) {
    const { chain_id } = chain;
    const options = { upsert: true, new: false, setDefaultsOnInsert: true };
    await this.findOneAndUpdate({ chain_id }, chain, options);
  },
};
