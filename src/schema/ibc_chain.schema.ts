import * as mongoose from 'mongoose';
import { dateNow } from 'src/helper/date.helper';
import { IbcChainType } from '../types/schemaTypes/ibc_chain.interface';
export const IbcChainSchema = new mongoose.Schema({
  chain_id: String,
  chain_name: String,
  icon: String,
  create_at: {
    type: String,
    default: dateNow,
  },
  update_at: {
    type: String,
    default: dateNow,
  },
});

IbcChainSchema.index({ chain_id: 1 }, { unique: true });

IbcChainSchema.statics = {
  async findAll(): Promise<IbcChainType[]> {
    return this.find();
  },

  async findActive(): Promise<IbcChainType[]> {
    return this.find({
      update_at: { $gte: String(Number(dateNow) - 24 * 60 * 60) },
    }).collation( { locale: 'en_US' } ).sort({'chain_id': 1});
  },

  async countActive(): Promise<number> {
    return this.count({
      update_at: { $gte: String(Number(dateNow) - 24 * 60 * 60) },
    });
  },

  async findById(chain_id): Promise<IbcChainType> {
    return this.findOne({ chain_id });
  },

  async updateChainRecord(chain): Promise<void> {
    const { chain_id } = chain;
    const options = { upsert: true, new: false, setDefaultsOnInsert: true };
    return this.findOneAndUpdate({ chain_id }, chain, options);
  },

  async insertManyChain(chain): Promise<void> {
    return this.insertMany(chain, { ordered: false });
  },
};
