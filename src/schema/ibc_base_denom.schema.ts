import * as mongoose from 'mongoose';
import { IbcBaseDenomType } from '../types/schemaTypes/ibc_base_denom.interface';

export const IbcBaseDenomSchema = new mongoose.Schema(
  {
    chain_id: String,
    denom: String,
    symbol: String,
    scale: Number,
    icon: String,
    is_main_token: Boolean,
    create_at: {
      type: Number,
      default: Math.floor(new Date().getTime() / 1000),
    },
    update_at: {
      type: Number,
      default: Math.floor(new Date().getTime() / 1000),
    },
  },
  { versionKey: false },
);

IbcBaseDenomSchema.index({ chain_id: 1, denom: 1 }, { unique: true });
IbcBaseDenomSchema.index({ update_at: -1 }, { background: true });

IbcBaseDenomSchema.statics = {
  async findCount(): Promise<number> {
    return this.count();
  },

  async findAllRecord(): Promise<IbcBaseDenomType[]> {
    return this.find();
  },
};
