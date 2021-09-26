import * as mongoose from 'mongoose';
import { IbcChannelType } from '../types/schemaTypes/ibc_channel.interface';
import { dateNow } from '../helper/date.helper';

export const IbcChannelSchema = new mongoose.Schema(
  {
    channel_id: String,
    record_id: String,
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

IbcChannelSchema.index({ record_id: 1 }, { unique: true });
IbcChannelSchema.index({ update_at: -1 }, { background: true });

IbcChannelSchema.statics = {
  async findCount(query): Promise<number> {
    return this.count(query);
  },

  async findChannelRecord(record_id): Promise<IbcChannelType> {
    return this.findOne({ record_id }, { _id: 0 });
  },

  async updateChannelRecord(channelRecord): Promise<void> {
    const { record_id } = channelRecord;
    const options = { upsert: true, new: false, setDefaultsOnInsert: true };
    return this.findOneAndUpdate({ record_id }, channelRecord, options);
  },

  async insertManyChannel(ibcChannel): Promise<void> {
    return this.insertMany(ibcChannel, { ordered: false });
  },
};
