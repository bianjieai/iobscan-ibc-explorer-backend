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
  //todo  1.rename： findBlockByLastHeight => findLatestBlock    2.方法声明未明确返回值类型
  async findBlockByLastHeight() {
    return this.findOne({})
      .sort({ height: -1 })
      .limit(1);
  },
};
