import * as mongoose from 'mongoose';

export const IbcChannelSchema = new mongoose.Schema(
  {
    channel_id: String,
    record_id: String,
    update_at: String,
    create_at: String,
  },
  { versionKey: false },
);

IbcChannelSchema.index({ record_id: 1, denom: 1 }, { unique: true });
IbcChannelSchema.index({ update_at: -1 }, { background: true });

IbcChannelSchema.statics = {
  // 查
  async findCount(query) {
    return this.count(query);
  },

  async findChannelRecord(record_id) {
    return this.findOne({ record_id }, { _id: 0 });
  },

  // 改
  async updateChannelRecord(channelRecord) {
    const { record_id } = channelRecord;
    const options = { upsert: true, new: false, setDefaultsOnInsert: true };
    return this.findOneAndUpdate({ record_id }, channelRecord, options);
  },

  // 增
  async insertManyChannel(ibcChannel) {
    return this.insertMany(ibcChannel, { ordered: false });
  },
};
