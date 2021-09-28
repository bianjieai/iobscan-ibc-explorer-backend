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
  log: {
    sc_log: String,
    dc_log: String,
  },
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
IbcTxSchema.index({ update_at: -1 }, { background: true });

IbcTxSchema.statics = {
  // todo query
  async findCount(query): Promise<number> {
    return this.count(query);
  },

  async findTxList(pageNum: number, pageSize: number): Promise<IbcTxType[]> {
    return this.find({ status: { $in: [1, 2, 3, 4] } }, { _id: 0 })
      .skip((Number(pageNum) - 1) * Number(pageSize))
      .limit(Number(pageSize))
      .sort({ update_at: -1 });
  },

  async queryTxList(query): Promise<IbcTxType[]> {
    const { status, limit } = query;
    return this.find({ status }, { _id: 0 })
      .sort({ update_at: -1 })
      .limit(Number(limit));
  },

  async distinctChainList(query): Promise<any> {
    const { type, dateNow, status } = query;
    return this.distinct(type, {
      update_at: { $gte: String(dateNow - 24 * 60 * 60) },
      status: { $in: status },
    });
  },

  async updateIbcTx(ibcTx): Promise<void> {
    const { record_id } = ibcTx;
    const options = { upsert: true, new: false, setDefaultsOnInsert: true };
    return this.findOneAndUpdate({ record_id }, ibcTx, options);
  },

  async insertManyIbcTx(ibcTx, cb): Promise<void> {
    return this.insertMany(ibcTx, cb);
  },
};
