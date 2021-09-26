import * as mongoose from 'mongoose';
export const SyncTaskSchema = new mongoose.Schema({
  start_height: Number,
  end_height: Number,
  current_height: Number,
  status: String,
  worker_id: String,
  worker_logs: Object,
  last_update_time: Number,
});

SyncTaskSchema.statics = {
  async queryTaskCount(end_height: number = 0, status: string = 'underway') {
    return this.find({
      end_height: end_height,
      status: status,
    }).countDocuments();
  },
};
