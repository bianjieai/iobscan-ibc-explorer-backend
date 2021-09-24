import * as mongoose from 'mongoose';
import { IbcTxType } from '../types/schemaTypes/ibc_tx.interface';
import { dateNow } from '../helper/date.helper';

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
  create_at: {
    type: String,
    default: dateNow,
  },
  update_at: {
    type: String,
    default: dateNow,
  },
});

IbcTxSchema.index({ record_id: -1 }, { unique: true });

IbcTxSchema.statics = {
  // todo 方法命名规范  明确query入参类型
  // 查
  async findCount(query) {
    return this.count(query);
  },

  async findTxList(pageNum: number, pageSize: number): Promise<IbcTxType> {
    return this.find({ status: { "$in": [ 1, 2, 3, 4 ] }}, { _id: 0 })
      .skip((Number(pageNum) - 1) * Number(pageSize))
      .limit(Number(pageSize));
  },

  async queryTxList(query): Promise<any> {
    const { status, limit } = query;
    return this.find({ status }, { _id: 0 })
      .sort({ update_at: 1 })
      .limit(Number(limit));
  },

  async distinctChainList(query) {
    const { type, dateNow, status } = query;
    return this.distinct(type, {
      update_at: { $gte: String(dateNow - 24 * 60 * 60) },
      status: { $in: status },
    });
  },

  // 改
  async updateIbcTx(ibcTx) {
    const { record_id } = ibcTx;
    const options = { upsert: true, new: false, setDefaultsOnInsert: true };
    return this.findOneAndUpdate({ record_id }, ibcTx, options);
  },

  // 增
  async insertManyIbcTx(ibcTx, cb) {
    return this.insertMany(ibcTx, cb);
  },
};
