import * as mongoose from 'mongoose';
import { dateNow } from '../helper/date.helper'

export const IbcBaseDenomSchema = new mongoose.Schema(
  {
    chain_id: String,
    denom: String,
    symbol: String,
    scale: String,
    icon: String,
    is_main_token: Boolean,
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

IbcBaseDenomSchema.index({ chain_id: 1, denom: 1 }, { unique: true });
IbcBaseDenomSchema.index({ update_at: -1 }, { background: true });

// todo 未确定返回值类型
IbcBaseDenomSchema.statics = {
  // 查
  async findCount() {
    return this.count();
  },

  async findAllRecord() {
    return this.find();
  },
};
