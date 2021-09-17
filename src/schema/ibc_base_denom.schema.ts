import * as mongoose from 'mongoose';

export const IbcBaseDenomSchema = new mongoose.Schema(
  {
    chain_id: String,
    denom: String,
    symbol: String,
    scale: String,
    icon: String,
    is_main_token: Boolean,
    create_at: String,
    update_at: String,
  },
  { versionKey: false },
);

IbcBaseDenomSchema.index({ chain_id: 1, denom: 1 }, { unique: true });
IbcBaseDenomSchema.index({ update_at: -1 }, { background: true });

IbcBaseDenomSchema.statics = {
  // 查
  async findCount() {
    return await this.count();
  },
};
