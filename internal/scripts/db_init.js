db.ibc_chain.createIndex({'chain_id': -1}, {background: true, unique: true});

db.ibc_relayer.createIndex({
    "chain_a": -1,
    "channel_a": -1,
    "chain_a_address": -1
}, {background: true, unique: true});

db.ibc_relayer.createIndex({
    "chain_b": -1,
    "channel_b": -1,
    "chain_b_address": -1
}, {background: true, unique: true});

db.ibc_relayer_config.createIndex({
    "relayer_channel_pair": -1
}, {background: true, unique: true});

db.ibc_relayer_statistics.createIndex({
    "transfer_base_denom": 1,
    "address": 1,
    "statistic_id": 1,
    "segment_start_time": -1,
    "segment_end_time": -1
}, {
    unique: true,
    background: true
});


db.ibc_channel.createIndex({
    "channel_id": 1
}, {background: true, unique: true});


db.ibc_channel_statistics.createIndex({
    "channel_id": 1,
    "base_denom": 1,
    "base_denom_chain_id": 1,
    "segment_start_time": -1,
    "segment_end_time": -1
}, {
    name: "channel_sstatistics_unique",
    unique: true,
    background: true
});

db.ibc_token.createIndex({
    "base_denom": 1,
    "chain_id": 1
}, {background: true, unique: true});


db.ibc_token_statistics.createIndex({
    "base_denom": 1,
    "base_denom_chain_id": 1,
    "segment_start_time": -1,
    "segment_end_time": -1
}, {
    unique: true,
    background: true
});

db.ibc_token_trace.createIndex({
    "denom": 1,
    "chain_id": 1,
}, {
    background: true,
    unique: true
});

db.ibc_token_trace_statistics.createIndex({
    "denom": 1,
    "chain_id": 1,
    "segment_start_time": -1,
    "segment_end_time": -1
}, {
    unique: true,
    background: true
});