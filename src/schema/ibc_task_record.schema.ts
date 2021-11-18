import * as mongoose from 'mongoose';
import { IbcTaskRecordType } from '../types/schemaTypes/ibc_task_record.interface';
import { IbcTaskRecordStatus } from '../constant';

export const IbcTaskRecordSchema = new mongoose.Schema({
  _id: String,
  task_name: String,
  status: {
    type: String,
    default: IbcTaskRecordStatus.OPEN,
  },
  height: Number,
  create_at: {
    type: Number,
    default: Math.floor(new Date().getTime() / 1000),
  },
  update_at: {
    type: Number,
    default: Math.floor(new Date().getTime() / 1000),
  },
});

IbcTaskRecordSchema.index({ task_name: 1 }, { unique: true });

IbcTaskRecordSchema.statics = {
  async findAllChainConfig(): Promise<IbcTaskRecordType[]> {
    return this.find();
  },
  async findTaskRecord(task_id): Promise<IbcTaskRecordType> {
    return await this.findOne(
      { task_name: `sync_${task_id}_transfer` },
      { _id: 0 },
    );
  },

  async updateTaskRecord(taskRecord: IbcTaskRecordType): Promise<void> {
    const { task_name } = taskRecord;
    const options = { upsert: true, new: false, setDefaultsOnInsert: true };
    return this.findOneAndUpdate({ task_name }, taskRecord, options);
  },

  async insertManyTaskRecord(taskRecords: IbcTaskRecordType[]): Promise<void> {
    return this.insertMany(taskRecords, { ordered: false });
  },
};
