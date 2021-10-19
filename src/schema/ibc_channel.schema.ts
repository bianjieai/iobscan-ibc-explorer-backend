/* eslint-disable @typescript-eslint/camelcase */
import * as mongoose from 'mongoose';
import { IbcChannelType } from '../types/schemaTypes/ibc_channel.interface';

export const IbcChannelSchema = new mongoose.Schema(
  {
    channel_id: String,
    record_id: String,
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
  },
  { versionKey: false },
);

IbcChannelSchema.index({ record_id: 1 }, { unique: true });
IbcChannelSchema.index({ update_at: -1 }, { background: true });

IbcChannelSchema.statics = {

  async countActive(): Promise<number> {
    return this.count({
      tx_time: { $gte: Math.floor(new Date().getTime() / 1000) - 24 * 60 * 60 },
    });
  },

  async findChannelRecord(record_id): Promise<IbcChannelType> {
    return this.findOne({ record_id }, { _id: 0 });
  },

  async updateChannelRecord(channelRecord): Promise<void> {
    const { record_id, update_at, tx_time } = channelRecord;
    const options = { upsert: true, new: true, setDefaultsOnInsert: true };
    return this.findOneAndUpdate({ record_id }, { $set: { update_at, tx_time } }, options);
  },

  async insertManyChannel(ibcChannel): Promise<void> {
    return this.insertMany(ibcChannel, { ordered: false });
  },
};
