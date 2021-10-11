import * as mongoose from 'mongoose';
import { IbcStatisticsType } from '../types/schemaTypes/ibc_statistics.interface';
import { dateNow } from '../helper/date.helper';

export const IbcStatisticsSchema = new mongoose.Schema(
  {
    statistics_name: String,
    count: Number,
    create_at: {
      type: Number,
      default: dateNow,
    },
    update_at: {
      type: Number,
      default: dateNow,
    },
  },
  { versionKey: false },
);

// todo 冗余的索引

IbcStatisticsSchema.index({ statistics_name: 1 }, { unique: true });

IbcStatisticsSchema.statics = {
  async findStatisticsRecord(
    statistics_name: string,
  ): Promise<IbcStatisticsType> {
    return this.findOne({ statistics_name }, { _id: 0 });
  },

  async findAllRecord(): Promise<IbcStatisticsType[]> {
    return this.find();
  },

  async updateStatisticsRecord(
    statisticsRecord: IbcStatisticsType,
    cb,
  ): Promise<void> {
    const { statistics_name } = statisticsRecord;
    const options = { upsert: true, new: false, setDefaultsOnInsert: true };
    return this.findOneAndUpdate(
      { statistics_name },
      statisticsRecord,
      options,
      cb,
    );
  },

  async insertManyStatisticsRecord(
    statisticsRecord: IbcStatisticsType,
  ): Promise<void> {
    return this.insertMany(statisticsRecord, { ordered: false });
  },
};
