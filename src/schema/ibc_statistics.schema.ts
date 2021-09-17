import * as mongoose from 'mongoose';

export const IbcStatisticsSchema = new mongoose.Schema(
  {
    statistics_name: String,
    count: Number,
    create_at: String,
    update_at: String,
  },
  { versionKey: false },
);

IbcStatisticsSchema.index({ statistics_name: 1 }, { unique: true });
IbcStatisticsSchema.index({ update_at: -1 }, { background: true });

IbcStatisticsSchema.statics = {
  // 查
  async findStatisticsRecord(statistics_name, cb) {
    return  this.findOne({ statistics_name }, { _id: 0 }, cb);
  },

  async findAllRecord() {
    return  this.find({})
  },

  // 改
  async updateStatisticsRecord(statisticsRecord, cb) {
    const { statistics_name } = statisticsRecord;
    const options = { upsert: true, new: false, setDefaultsOnInsert: true };
    return  this.findOneAndUpdate(
      { statistics_name },
      statisticsRecord,
      options,
      cb,
    );
  },

  // 增
  async insertManyStatisticsRecord(statisticsRecord) {
    return  this.insertMany(statisticsRecord, { ordered: false });
  },
};
