import * as mongoose from 'mongoose';
export const TxSchema = new mongoose.Schema(
  {
    time: Number,
    height: Number,
    tx_hash: String,
    memo: String,
    status: Number,
    log: String,
    complex_msg: Boolean,
    type: String,
    from: String,
    to: String,
    coins: Array,
    signer: String,
    events: Array,
    events_new: Array,
    msgs: Array,
    signers: Array,
    addrs: Array,
    fee: Object,
    tx_index: Number,
  },
  { versionKey: false },
);
TxSchema.index({ time: -1, 'msgs.type': -1, status: -1 }, { background: true });
TxSchema.index({ height: 1 }, { background: true });
TxSchema.index({ addrs: -1, time: -1, status: -1 }, { background: true });
TxSchema.index(
  { 'msgs.type': -1, height: -1, 'msgs.msg.ex.service_name': -1 },
  { background: true },
);

// 	txs
TxSchema.statics = {
  async queryTxListSortHeight(query): Promise<any> {
    const { type, height, limit, status } = query;
    return await this.find({ type, height: { $gt: height } })
      .sort({ height: 1 })
      .limit(Number(limit));
  },

  async queryTxListSortUpdateAt(query): Promise<any> {
    const { type, limit, status } = query;
    return await this.find({ type, status })
      .sort({ update_at: 1 })
      .limit(Number(limit));
  },

  async queryTxListByPacketId(query): Promise<any> {
    const { type, limit, status, packet_id } = query;
    return await this.find({
      'msgs.type': type,
      status,
      'msgs.msg.packet_id': packet_id,
    })
      .sort({ update_at: 1 })
      .limit(Number(limit));
  },
};
