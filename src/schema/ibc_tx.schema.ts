import * as mongoose from 'mongoose';
import { IbcTxType } from '../types/schemaTypes/ibc_tx.interface';

export const IbcTxSchema = new mongoose.Schema({
  record_id: String,
  sc_addr: String,
  dc_addr: String,
  sc_port: String,
  sc_channel: String,
  sc_chain_id: String,
  dc_port: String,
  dc_channel: String,
  dc_chain_id: String,
  sequence: String,
  status: Number,
  sc_tx_info: Object,
  dc_tx_info: Object,
  refunded_tx_info: Object,
  log: Array,
  denoms: Array,
  base_denom: String,
  create_at: String,
  update_at: String,
});

IbcTxSchema.index({ record_id: -1 }, { unique: true });

IbcTxSchema.statics = {
  // 查
  async findCount(query) {
    const result = await this.count(query);
    return result;
  },

  async findTxList(pageNum: number, pageSize: number): Promise<IbcTxType> {
    return await this.find({}, { _id: 0 })
      .skip((Number(pageNum) - 1) * Number(pageSize))
      .limit(Number(pageSize));
  },

  async queryTxList(query): Promise<any> {
    const { status, limit } = query;
    return await this.find({ status }, { _id: 0 })
      .sort({ update_at: 1 })
      .limit(Number(limit));
  },

  async distinctChainList(query) {
    const { type, dateNow, status } = query;
    return await this.distinct(type, {
      update_at: { $gte: String(dateNow - 24 * 60 * 60 * 1000) },
      status: { $in: status },
    });
  },

  // 改
  async updateIbcTx(ibcTx, cb) {
    const { record_id } = ibcTx;
    const options = { upsert: true, new: false, setDefaultsOnInsert: true };
    return await this.findOneAndUpdate({ record_id }, ibcTx, options, cb);
  },

  // 增
  async insertTx(ibcTx, cb) {
    return await this.insertMany(ibcTx, cb);
  },
};
