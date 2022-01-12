import dotenv from 'dotenv';
dotenv.config();
const {
    LCD_ADDR,
    RPC_ADDR,
    DB_USER,
    DB_PASSWD,
    DB_ADDR,
    ICONURI,
    DB_DATABASE,
    NODE_ENV,
    EXECUTE_KEY,
    DISPLAY_IBC_RECORD_MAX,
    FAULT_TOLERANCE_EXECUTE_TIME,
    SYNC_TX_SERVICE_NAME_SIZE,
    HEARTBEAT_RATE,
    DisableLog,
    INCREASE_HEIGHT,
    MAX_OPERATE_TX_COUNT,
    CRON_JOBS,
    PROPOSALS_LIMIT,
    IBCTX_EXECUTE_TIME,
    IBCCHAIN_EXECUTE_TIME,
    IBCSTATISTICS_EXECUTE_TIME,
    CHANNELS_LIMITS,
    CHANNELS_OFFSET,
    SYNC_TRANSFER_TX_TIME,
    UPDATE_PROCESSING_TX_TIME,
    UPDATE_SUB_STATE_TX_TIME,
    IBC_TX_LATEST_MIGRATE,
    IBC_TX_UPDATE_CRONJOB,
    IBC_DENOM_CACULATE_CRONJOB,
    IBC_DENOM_UPDATE_CRONJOB

} = process.env;
export const cfg = {
    channels:{
        limit: CHANNELS_LIMITS || 1000,
        offset: CHANNELS_OFFSET || 0,
    },
    env: NODE_ENV,
    disableLog:Boolean(DisableLog=='true'),
    dbCfg: {
        user: DB_USER,
        psd: DB_PASSWD,
        dbAddr: DB_ADDR,
        dbName: DB_DATABASE,
    },
    serverCfg:{
        lcdAddr:LCD_ADDR,
        rpcAddr:RPC_ADDR,
        iconUri: ICONURI || 'https://keybase.io/_/api/1.0/user/lookup.json',
        executeKey: EXECUTE_KEY,
        displayIbcRecordMax:Number(DISPLAY_IBC_RECORD_MAX || 500000),
    },
    taskCfg:{
        interval:{
            heartbeatRate:Number(HEARTBEAT_RATE || 10000),
        },
        executeTime:{
            // tx: '*/10 * * * * *',
            // chain: '*/10 * * * * *',
            // statistics: '*/10 * * * * *',
            tx: IBCTX_EXECUTE_TIME || '15 * * * * *',
            chain: IBCCHAIN_EXECUTE_TIME || '* * */1 * * *',
            statistics: IBCSTATISTICS_EXECUTE_TIME || '* */10 * * * *',

            faultTolerance:FAULT_TOLERANCE_EXECUTE_TIME || '41 * * * * *',
            transferTx: SYNC_TRANSFER_TX_TIME || '*/15 * * * * *',
            updateProcessingTx: UPDATE_PROCESSING_TX_TIME || '*/15 * * * * *',
            updateSubStateTx: UPDATE_SUB_STATE_TX_TIME || '*/15 * * * * *',
            ibcTxLatestMigrate: IBC_TX_LATEST_MIGRATE || '* */30 * * * *',
            ibcTxUpdateCronjob: IBC_TX_UPDATE_CRONJOB || '* * */2 * * *',
            ibcDenomCaculateCronjob: IBC_DENOM_CACULATE_CRONJOB || '* * */1 * * *',
            ibcDenomUpdateCronjob: IBC_DENOM_UPDATE_CRONJOB || '*/30 * * * * *',
        },
        syncTxServiceNameSize: Number(SYNC_TX_SERVICE_NAME_SIZE) || 200,
        increaseHeight: INCREASE_HEIGHT || 1000,
        maxOperateTxCount: MAX_OPERATE_TX_COUNT || 100,
        CRON_JOBS: CRON_JOBS ? JSON.parse(CRON_JOBS) : [],
        proposalsLimit: PROPOSALS_LIMIT || 1000,
    },
};

