import * as mongoose from 'mongoose';

export const IbcDenomSchema = new mongoose.Schema(
  {
    chain_id: String,
    denom: String,
    base_denom: String,
    base_denom_chain_id: String,
    denom_path: String,
    is_source_chain: String,
    create_at: String,
    update_at: String,
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

  // 增
  async insertManyDenom(ibcDenom, cb) {
    return this.insertMany(ibcDenom, { ordered: false }, cb);
  },
};
