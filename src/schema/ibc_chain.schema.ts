import * as mongoose from 'mongoose';
import { IbcChainType } from '../types/schemaTypes/ibc_chain.interface';
export const IbcChainSchema = new mongoose.Schema({
  chain_id: String,
  chain_name: String,
  icon: String,
  create_at: {
    type: Number,
    default: Math.floor(new Date().getTime() / 1000),
  },
  update_at: {
    type: Number,
    default: Math.floor(new Date().getTime() / 1000),
  },
  tx_time: {
    type: Number,
    default: Math.floor(new Date().getTime() / 1000),
  },
});

IbcChainSchema.index({ chain_id: 1 }, { unique: true });
IbcChainSchema.index({ tx_time: -1 }, { unique: true });

IbcChainSchema.statics = {
  async findAll(): Promise<IbcChainType[]> {
    return this.find();
  },

  async findActive(): Promise<IbcChainType[]> {
    return this.find({
      tx_time: { $gte: Math.floor(new Date().getTime() / 1000) - 24 * 60 * 60 },
    }).collation( { locale: 'en_US' } ).sort({'chain_id': 1});
  },

  async countActive(): Promise<number> {
    return this.count({
      tx_time: { $gte: Math.floor(new Date().getTime() / 1000) - 24 * 60 * 60 },
    });
  },

  async findById(chain_id): Promise<IbcChainType> {
    return this.findOne({ chain_id });
  },

  async updateChainRecord(chain): Promise<void> {
    const { chain_id, update_at, tx_time } = chain;
    const options = { upsert: true, new: true, setDefaultsOnInsert: true };
    return this.findOneAndUpdate({ chain_id }, { $set: { update_at, tx_time }}, options);
  },

  async insertManyChain(chain): Promise<void> {
    return this.insertMany(chain, { ordered: false });
  },
};
