// chain_config表
db.chain_config.createIndex({'current_chain_id': -1}, {background: true, unique: true});
// ibc_chain表
db.ibc_chain.createIndex({"chain": -1}, {background: true, unique: true});

// chain_registry
db.chain_registry.createIndex({
    "chain": 1
}, {
    unique: true,
    background: true
});

db.getCollection("auth_denom").createIndex({
    "chain": 1,
    "denom": 1
}, {
    background: true,
    unique: true
});

db.getCollection("ibc_denom").createIndex({
    "chain": 1,
    "denom": 1
}, {
    background: true,
    unique: true
});

// ibc_relayer表
db.ibc_relayer.createIndex({
    "channel_pair_info.pair_id": -1,
}, {background: true, unique: true});

db.ibc_relayer.createIndex({
    "relayer_id": -1,
}, {background: true, unique: true});

db.ibc_relayer.createIndex({
    "relayer_name": -1,
}, {background: true});

db.ibc_relayer.createIndex({
    "channel_pair_info.chain_a_address": -1
}, {background: true});

db.ibc_relayer.createIndex({
    "channel_pair_info.chain_b_address": -1
}, {background: true});

db.ibc_relayer.createIndex({
    "channel_pair_info.chain_a": -1,
    "channel_pair_info.channel_a": -1,
}, {background: true});

db.ibc_relayer.createIndex({
    "channel_pair_info.chain_b": -1,
    "channel_pair_info.channel_b": -1,
}, {background: true});


// ibc_channel表

db.ibc_channel.createIndex({
    "channel_id": 1
}, {background: true, unique: true});


// ibc_channel_statistics表

db.ibc_channel_statistics.createIndex({
    "channel_id": 1,
    "base_denom": 1,
    "base_denom_chain": 1,
    "status": 1,
    "segment_start_time": -1,
    "segment_end_time": -1
}, {
    name: "channel_statistics_unique",
    unique: true,
    background: true
});

// ibc_token表

db.ibc_token.createIndex({
    "base_denom": 1,
    "chain": 1
}, {background: true, unique: true});

// ibc_token_statistics表

db.ibc_token_statistics.createIndex({
    "base_denom": 1,
    "base_denom_chain": 1,
    "segment_start_time": -1,
    "segment_end_time": -1
}, {
    unique: true,
    background: true
});

// ibc_token_trace表
db.ibc_token_trace.createIndex({
    "denom": 1,
    "chain": 1,
}, {
    background: true,
    unique: true
});

// ibc_token_trace_statistics表
db.ibc_token_trace_statistics.createIndex({
    "denom": 1,
    "chain": 1,
    "segment_start_time": -1,
    "segment_end_time": -1
}, {
    unique: true,
    background: true
});


// ex_ibc_tx表
db.ex_ibc_tx.createIndex({
    "sc_tx_info.hash": 1,
    "sc_tx_info.height": 1,
    "sc_chain": 1,
    "sc_tx_info.msg.msg.packet_id": 1
}, {
    name: "sc_tx_unique",
    background: true,
    unique: true
});

db.ex_ibc_tx.createIndex({
    "dc_tx_info.hash": -1,
}, {
    background: true
});

db.ex_ibc_tx.createIndex({
    "ack_timeout_tx_info.hash": -1,
}, {
    background: true
});

db.ex_ibc_tx.createIndex({
    "status": 1,
    "sc_tx_info.status": 1
}, {
    background: true
});

db.ex_ibc_tx.createIndex({
    "dc_tx_info.status": 1
}, {
    background: true
});

db.ex_ibc_tx.createIndex({
    "tx_time": 1,
    "status": 1
}, {
    background: true
});

db.ex_ibc_tx.createIndex({
    "create_at": 1,
}, {
    background: true
});

db.ex_ibc_tx.createIndex({
    "sc_chain": 1,
    "sc_channel": 1,
    "status": 1
}, {
    background: true
});

db.getCollection("ex_ibc_tx").createIndex({
    "base_denom": 1,
    "base_denom_chain": 1,
    "status": 1
}, {
    background: true
});

db.getCollection("ex_ibc_tx").createIndex({
    "sc_chain": 1,
    "status": 1,
    "next_try_time": 1
}, {
    background: true
});

// ex_ibc_tx_latest表

db.getCollection("ex_ibc_tx_latest").createIndex({
    "denoms.sc_denom": 1,
    "sc_chain": 1
}, {
    background: true
});

db.getCollection("ex_ibc_tx_latest").createIndex({
    "denoms.dc_denom": 1,
    "dc_chain": 1
}, {
    background: true
});

db.ex_ibc_tx_latest.createIndex({
    "sc_tx_info.hash": -1,
}, {
    background: true
});

db.ex_ibc_tx_latest.createIndex({
    "dc_tx_info.hash": -1,
}, {
    background: true
});

db.ex_ibc_tx_latest.createIndex({
    "ack_timeout_tx_info.hash": -1,
}, {
    background: true
});

db.ex_ibc_tx_latest.createIndex({
    "tx_time": 1,
    "status": 1
}, {
    background: true
});

db.ex_ibc_tx_latest.createIndex({
    "dc_tx_info.status": 1
}, {
    background: true
});
db.ex_ibc_tx_latest.createIndex({
    "status": 1,
    "sc_tx_info.status": 1
}, {
    background: true
});

db.ex_ibc_tx_latest.createIndex({
    "dc_chain": 1,
    "status": 1
}, {
    background: true
});

db.ex_ibc_tx_latest.createIndex({
    "sc_chain": 1,
    "dc_chain": 1,
    "status": 1
}, {
    background: true
});

db.ex_ibc_tx_latest.createIndex({
    "create_at": 1
}, {
    background: true
});

db.getCollection("ex_ibc_tx_latest").createIndex({
    "sc_chain": 1,
    "sc_channel": 1,
    "status": 1
}, {
    background: true
});

db.getCollection("ex_ibc_tx_latest").createIndex({
    "base_denom": 1,
    "base_denom_chain": 1,
    "status": 1
}, {
    background: true
});

db.getCollection("ex_ibc_tx_latest").createIndex({
    "sc_chain": 1,
    "status": 1,
    "next_try_time": 1
}, {
    background: true
});

// sync_{chain}_tx表
db.sync_xxxx_tx.createIndex({"tx_hash": -1, "height": -1}, {unique: true, background: true});
db.sync_xxxx_tx.createIndex({"height": -1}, {background: true});
db.sync_xxxx_tx.createIndex({"types": -1, "height": -1}, {background: true});
db.sync_xxxx_tx.createIndex({"msgs.msg.packet_id": -1}, {background: true});
db.sync_xxxx_tx.createIndex({"msgs.msg.signer": 1, "msgs.type": 1, "time": 1}, {background: true});
db.sync_xxxx_tx.createIndex({"time": -1, "msgs.type": -1}, {background: true});

// uba_search_record
db.uba_search_record.createIndex({
    "create_at": 1
}, {
    expireAfterSeconds: 31536000
})

// relayer statistics
db.ibc_relayer_fee_statistics.createIndex({
    "chain_address_comb": 1,
    "tx_type": 1,
    "tx_status": 1,
    "fee_denom": 1,
    "segment_start_time": 1,
    "segment_end_time": 1,
}, {name: "statistics_unique", background: true, unique: true});

db.ibc_relayer_fee_statistics.createIndex({
    "statistics_chain": 1,
    "segment_start_time": 1,
    "segment_end_time": 1,
}, {background: true});

db.ibc_relayer_denom_statistics.createIndex({
    "chain_address_comb": 1,
    "segment_start_time": 1,
    "segment_end_time": 1,
}, {background: true});

db.ibc_relayer_denom_statistics.createIndex({
    "chain_address_comb": 1,
    "tx_type": 1,
    "tx_status": 1,
    "base_denom": 1,
    "base_denom_chain": 1,
    "segment_start_time": 1,
    "segment_end_time": 1,
}, {name: "statistics_unique", background: true, unique: true});

db.ibc_relayer_denom_statistics.createIndex({
    "statistics_chain": 1,
    "segment_start_time": 1,
    "segment_end_time": 1,
}, {background: true});

db.ibc_relayer_address_channel.createIndex({
    "relayer_address": 1,
    "chain": 1,
    "channel": 1
}, {background: true, unique: true});
