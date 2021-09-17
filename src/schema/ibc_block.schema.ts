import * as mongoose from 'mongoose';

export const IbcBlockSchema = new mongoose.Schema(
  {
    height: Number,
    hash: String,
    txn: Number,
    time: Number,
    proposer: String,
  },
  { versionKey: false },
);

IbcBlockSchema.statics = {
  async findBlockByLastHeight() {
    return await this.findOne({})
      .sort({ height: -1 })
      .limit(1);
  },
};
