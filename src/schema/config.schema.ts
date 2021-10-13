import * as mongoose from 'mongoose';
import { ConfigType } from '../types/schemaTypes/config.interface';

export const ConfigSchema = new mongoose.Schema(
  {
    iobscan: String,
  },
  { versionKey: false },
);

ConfigSchema.index({ name: 1 }, { unique: true });

ConfigSchema.statics = {
  async findRecord(): Promise<ConfigType[]> {
    return this.findOne();
  },
};
