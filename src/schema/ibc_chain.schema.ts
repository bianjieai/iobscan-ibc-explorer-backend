import * as mongoose from 'mongoose';
import { dateNow } from 'src/helper/date.helper';

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
  // 查
  async findAll(): Promise<any> {
    return this.find({});
  },

  async findActive(): Promise<any> {
    return this.find({ update_at: { $gte: String(Number(dateNow) - 24 * 60 * 60) }});
  },

  async countActive(): Promise<any> {
    return this.count({ update_at: { $gte: String(Number(dateNow) - 24 * 60 * 60) }});
  },

  async findById(chain_id): Promise<any> {
    return this.findOne({ chain_id });
  },

  // 改
  async updateChainRecord(chain): Promise<any> {
    const { chain_id } = chain;
    const options = { upsert: true, new: false, setDefaultsOnInsert: true };
    return this.findOneAndUpdate({ chain_id }, chain, options);
  },

  // 增
  async insertManyChain(chain): Promise<any> {
    return this.insertMany(chain, { ordered: false });
  },
};
