/* eslint-disable @typescript-eslint/camelcase */
import * as mongoose from 'mongoose';
import { IbcChannelType } from '../types/schemaTypes/ibc_channel.interface';
import { dateNow } from '../helper/date.helper';

export const IbcChannelSchema = new mongoose.Schema(
  {
    channel_id: String,
    record_id: String,
    create_at: {
      type: Number,
      default: dateNow,
    },
    update_at: {
      type: Number,
      default: dateNow,
    },
    tx_time: {
      type: Number,
      default: dateNow,
    },
  },
  { versionKey: false },
);

IbcChannelSchema.index({ record_id: 1 }, { unique: true });
IbcChannelSchema.index({ update_at: -1 }, { background: true });

IbcChannelSchema.statics = {

  async countActive(): Promise<number> {
    return this.count({
      tx_time: { $gte: dateNow - 24 * 60 * 60 },
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
