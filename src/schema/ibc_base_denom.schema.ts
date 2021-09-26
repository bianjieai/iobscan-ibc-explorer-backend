import * as mongoose from 'mongoose';
import { IbcBaseDenomType } from '../types/schemaTypes/ibc_base_denom.interface';
import { dateNow } from '../helper/date.helper';

export const IbcBaseDenomSchema = new mongoose.Schema(
  {
    chain_id: String,
    denom: String,
    symbol: String,
    scale: Number,
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

IbcBaseDenomSchema.statics = {
  async findCount(): Promise<number> {
    return this.count();
  },

  async findAllRecord(): Promise<IbcBaseDenomType[]> {
    return this.find();
  },
};
