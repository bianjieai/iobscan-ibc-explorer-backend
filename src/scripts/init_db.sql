-- create db index

db.chain_config.createIndex({"chain_id": 1}, {"unique": true});

db.ibc_base_denom.createIndex({"chain_id": 1, "denom": 1},{"unique": true});

db.ibc_denom.createIndex({"chain_id": 1, "denom": 1}, {"unique": true});
db.ibc_denom.createIndex({"symbol": -1}, {"background": true});

db.ex_ibc_tx.createIndex({"record_id": -1}, {"unique": true});
db.ex_ibc_tx.createIndex({"status": -1},{"background":true});
db.ex_ibc_tx.createIndex({"status": -1,"tx_time": -1, "sc_chain_id": -1, "denoms.sc_denom": -1},{"background":true});
db.ex_ibc_tx.createIndex({"status": -1,"tx_time": -1, "dc_chain_id": -1, "denoms.dc_denom": -1},{"background":true});

db.ibc_statistics.createIndex({ "statistics_name": 1 }, { "unique": true });

db.ibc_task_record.createIndex({ "task_name": 1 }, { "unique": true });


db.sync_irishub_1_tx.createIndex({"types": -1,"height": -1},{"background":true});
db.sync_irishub_1_tx.createIndex({"msgs.type": -1,"status": -1,"msgs.msg.packet_id":-1},{"background":true});

db.sync_cosmoshub_4_tx.createIndex({"types": -1,"height": -1},{"background":true});
db.sync_cosmoshub_4_tx.createIndex({"msgs.type": -1,"status": -1,"msgs.msg.packet_id":-1},{"background":true});

db.sync_crypto_org_chain_mainnet_1_tx.createIndex({"types": -1,"height": -1},{"background":true});
db.sync_crypto_org_chain_mainnet_1_tx.createIndex({"msgs.type": -1,"status": -1,"msgs.msg.packet_id":-1},{"background":true});

db.sync_core_1_tx.createIndex({"types": -1,"height": -1},{"background":true});
db.sync_core_1_tx.createIndex({"msgs.type": -1,"status": -1,"msgs.msg.packet_id":-1},{"background":true});


db.sync_akashnet_2_tx.createIndex({"types": -1,"height": -1},{"background":true});
db.sync_akashnet_2_tx.createIndex({"msgs.type": -1,"status": -1,"msgs.msg.packet_id":-1},{"background":true});

db.sync_emoney_3_tx.createIndex({"types": -1,"height": -1},{"background":true});
db.sync_emoney_3_tx.createIndex({"msgs.type": -1,"status": -1,"msgs.msg.packet_id":-1},{"background":true});

db.sync_iov_mainnet_ibc_tx.createIndex({"types": -1,"height": -1},{"background":true});
db.sync_iov_mainnet_ibc_tx.createIndex({"msgs.type": -1,"status": -1,"msgs.msg.packet_id":-1},{"background":true});

db.sync_microtick_1_tx.createIndex({"types": -1,"height": -1},{"background":true});
db.sync_microtick_1_tx.createIndex({"msgs.type": -1,"status": -1,"msgs.msg.packet_id":-1},{"background":true});

db.sync_osmosis_1_tx.createIndex({"types": -1,"height": -1},{"background":true});
db.sync_osmosis_1_tx.createIndex({"msgs.type": -1,"status": -1,"msgs.msg.packet_id":-1},{"background":true});

db.sync_regen_1_tx.createIndex({"types": -1,"height": -1},{"background":true});
db.sync_regen_1_tx.createIndex({"msgs.type": -1,"status": -1,"msgs.msg.packet_id":-1},{"background":true});

db.sync_sentinelhub_2_tx.createIndex({"types": -1,"height": -1},{"background":true});
db.sync_sentinelhub_2_tx.createIndex({"msgs.type": -1,"status": -1,"msgs.msg.packet_id":-1},{"background":true});

db.sync_sifchain_1_tx.createIndex({"types": -1,"height": -1},{"background":true});
db.sync_sifchain_1_tx.createIndex({"msgs.type": -1,"status": -1,"msgs.msg.packet_id":-1},{"background":true});

-- init data

db.getCollection("ibc_config").insertMany([{
    iobscan: "https://www.iobscan.io",
}]);

db.getCollection("chain_config").insertMany([
{
    chain_id: 'irishub_1',
    icon : "https://iobscan.io/resources/xp-chains/irishub-1.png",
    chain_name: 'irishub_1',
    lcd: 'https://irishub.stakesystems.io/',
}, {
    chain_id: 'cosmoshub_4',
    icon: 'https://iobscan.io/resources/xp-chains/cosmoshub-4.png',
    chain_name: 'cosmoshub_4',
    lcd: 'https://cosmoshub.stakesystems.io/',
},{
    chain_id: 'crypto_org_chain_mainnet_1',
    icon: 'https://iobscan.io/resources/xp-chains/crypto-org-chain-mainnet-1.png',
    chain_name: 'crypto_org_chain_mainnet_1',
    lcd: 'https://mainnet.crypto.org:1317/',
},{
    chain_id: 'core_1',
    icon: 'https://iobscan.io/resources/xp-chains/core-1.png',
    chain_name: 'core_1',
    lcd: 'https://rest.core.persistence.one/',
},{
    chain_id: 'akashnet_2',
    icon: 'https://iobscan.io/resources/xp-chains/akashnet-2.png',
    chain_name: 'akashnet_2',
    lcd: 'https://lcd.akash.forbole.com/',
},{
    chain_id: 'emoney_3',
    icon: 'https://iobscan.io/resources/xp-chains/emoney-3.png',
    chain_name: 'emoney_3',
    lcd: 'https://emoney.stakesystems.io/',
},{
    chain_id: 'iov_mainnet_ibc',
    icon: 'https://iobscan.io/resources/xp-chains/iov-mainnet-ibc.png',
    chain_name: 'iov_mainnet_ibc',
    lcd: 'https://lcd-iov.keplr.app',
},{
    chain_id: 'microtick_1',
    icon: 'https://iobscan.io/resources/xp-chains/microtick-1.png',
    chain_name: 'microtick_1',
    lcd: 'https://microtick.stakesystems.io/',
},{
    chain_id: 'osmosis_1',
    icon: 'https://iobscan.io/resources/xp-chains/osmosis-1.png',
    chain_name: 'osmosis_1',
    lcd: 'https://osmosis.stakesystems.io/',
},{
    chain_id: 'regen_1',
    icon: 'https://iobscan.io/resources/xp-chains/regen-1.png',
    chain_name: 'regen_1',
    lcd: 'https://regen.stakesystems.io/',
},{
    chain_id: 'sentinelhub_2',
    icon: 'https://iobscan.io/resources/xp-chains/sentinelhub-2.png',
    chain_name: 'sentinelhub_2',
    lcd: 'https://lcd.sentinel.co/',
},{
    chain_id: 'sifchain_1',
    icon: 'https://iobscan.io/resources/xp-chains/sifchain-1.png',
    chain_name: 'sifchain_1',
    lcd: 'https://api.sifchain.finance/',
}
]);

db.getCollection("ibc_base_denom").insertMany([
{
    "chain_id" : "irishub_1",
    "denom" : "uiris",
    "symbol" : "iris",
    "scale" : "6",
    "icon" : "https://iobscan.io/resources/xp-tokens/irishub-1.png",
    "is_main_token" : true,
    "create_at" : "1631775090023",
    "update_at" : "1631775090023"
},
{
    "chain_id" : "cosmoshub_4",
    "denom" : "uatom",
    "symbol" : "atom",
    "scale" : "6",
    "icon" : "https://iobscan.io/resources/xp-tokens/cosmoshub-4.png",
    "is_main_token" : true,
    "create_at" : "1631775090023",
    "update_at" : "1631775090023"
},
{
    "chain_id" : "crypto_org_chain_mainnet_1",
    "denom" : "basecro",
    "symbol" : "cro",
    "scale" : "8",
    "icon" : "https://iobscan.io/resources/xp-tokens/crypto-org-chain-mainnet-1.png",
    "is_main_token" : true,
    "create_at" : "1631775090023",
    "update_at" : "1631775090023"
},
{
    "chain_id" : "core_1",
    "denom" : "uxprt",
    "symbol" : "xprt",
    "scale" : "6",
    "icon" : "https://iobscan.io/resources/xp-tokens/core-1.png",
    "is_main_token" : true,
    "create_at" : "1631775090023",
    "update_at" : "1631775090023"
},
{
    "chain_id" : "akashnet_2",
    "denom" : "uakt",
    "symbol" : "akt",
    "scale" : "6",
    "icon" : "https://iobscan.io/resources/xp-tokens/akashnet-2.png",
    "is_main_token" : true,
    "create_at" : "1631775090023",
    "update_at" : "1631775090023"
},
{
    "chain_id" : "emoney_3",
    "denom" : "ungm",
    "symbol" : "ngm",
    "scale" : "6",
    "icon" : "https://iobscan.io/resources/xp-tokens/emoney-3.png",
    "is_main_token" : true,
    "create_at" : "1631775090023",
    "update_at" : "1631775090023"
},
{
    "chain_id" : "iov_mainnet_ibc",
    "denom" : "uiov",
    "symbol" : "iov",
    "scale" : "6",
    "icon" : "https://iobscan.io/resources/xp-tokens/iov-mainnet-ibc.png",
    "is_main_token" : true,
    "create_at" : "1631775090023",
    "update_at" : "1631775090023"
},
{
    "chain_id" : "microtick_1",
    "denom" : "utick",
    "symbol" : "tick",
    "scale" : "6",
    "icon" : "https://iobscan.io/resources/xp-tokens/microtick-1.png",
    "is_main_token" : true,
    "create_at" : "1631775090023",
    "update_at" : "1631775090023"
},
{
    "chain_id" : "osmosis_1",
    "denom" : "uosmo",
    "symbol" : "osmo",
    "scale" : "6",
    "icon" : "https://iobscan.io/resources/xp-tokens/osmosis-1.png",
    "is_main_token" : true,
    "create_at" : "1631775090023",
    "update_at" : "1631775090023"
},
{
    "chain_id" : "regen_1",
    "denom" : "uregen",
    "symbol" : "regen",
    "scale" : "6",
    "icon" : "https://iobscan.io/resources/xp-tokens/regen-1.png",
    "is_main_token" : true,
    "create_at" : "1631775090023",
    "update_at" : "1631775090023"
},
{
    "chain_id" : "sentinelhub_2",
    "denom" : "udvpn",
    "symbol" : "dvpn",
    "scale" : "6",
    "icon" : "https://iobscan.io/resources/xp-tokens/sentinelhub-2.png",
    "is_main_token" : true,
    "create_at" : "1631775090023",
    "update_at" : "1631775090023"
},
{
    "chain_id" : "sifchain_1",
    "denom" : "rowan",
    "symbol" : "rowan",
    "scale" : "18",
    "icon" : "https://iobscan.io/resources/xp-tokens/sifchain-1.png",
    "is_main_token" : true,
    "create_at" : "1631775090023",
    "update_at" : "1631775090023"
}
]);