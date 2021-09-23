import * as mongoose from 'mongoose';
import { dateNow } from '../helper/date.helper';

export const IbcDenomSchema = new mongoose.Schema(
  {
    chain_id: String,
    denom: String,
    base_denom: String,
    denom_path: String,
    is_source_chain: String,
    create_at: {
      type: String,
      default: dateNow,
    },
    update_at: {
      type: String,
      default: dateNow,
    },
  },
  { versionKey: false },
);

IbcDenomSchema.index({ chain_id: 1, denom: 1 }, { unique: true });
IbcDenomSchema.index({ update_at: -1 }, { background: true });

IbcDenomSchema.statics = {
  // 查
  async findCount() {
    return this.count();
  },

  async findDenomRecord(chain_id, denom) {
    return this.findOne({ chain_id, denom }, { _id: 0 });
  },

  // 改
  async updateDenomRecord(denomRecord) {
    const { chain_id, denom } = denomRecord;
    const options = { upsert: true, new: false, setDefaultsOnInsert: true };
    return this.findOneAndUpdate({ chain_id, denom }, denomRecord, options);
  },

  // 增
  async insertManyDenom(ibcDenom) {
    return this.insertMany(ibcDenom, { ordered: false });
  },


};
