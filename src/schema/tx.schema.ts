import * as mongoose from 'mongoose';
export const TxSchema = new mongoose.Schema(
  {
    time: Number,
    height: Number,
    tx_hash: String,
    memo: String,
    status: Number,
    log: String,
    type: String,
    types: Array,
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

TxSchema.index({ 'types': -1, height: -1 }, { background: true });
TxSchema.index({ 'msgs.type': -1,status: -1 }, { background: true });
TxSchema.index({ 'msgs.msg.packet_id': -1 }, { background: true });

// 	txs
TxSchema.statics = {
  async queryTxListSortHeight(query): Promise<any> {
    const { type, height, limit } = query;
    return this.find({ 'types': type, height: { $gte: height } })
      .sort({ height: 1 })
      .limit(Number(limit));
  },

    async queryTxsByPacketId(query): Promise<any> {
        const { type,limit, packet_id } = query;
        return this.find({
            'msgs.type': type,
            'msgs.msg.packet_id': {$in:packet_id},
        }).limit(Number(limit));
    },

  async queryTxListByPacketId(query): Promise<any> {
    const { type,limit, status, packet_id } = query;
    return this.find({
      'msgs.type': type,
      status,
      'msgs.msg.packet_id': {$in:packet_id},
    }).limit(Number(limit));
  },

  async queryTxListByHeight(type, height): Promise<any> {
    return this.find({'types': type, height });
  },
};
