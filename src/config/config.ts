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
    MAX_PAGE_SIZE,
    DISPLAY_IBC_RECORD_MAX,
    FAULT_TOLERANCE_EXECUTE_TIME,
    HEARTBEAT_RATE,
    DisableLog,
    CRON_JOBS,
    UPDATE_DENOM_BATCH_LIMIT,
    UPDATE_IBC_TX_BATCH_LIMIT,
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
    IBC_MONITOR_CRONJOB,
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
        updateDenomBatchLimit: Number(UPDATE_DENOM_BATCH_LIMIT || 100),
        updateIbcTxBatchLimit: Number(UPDATE_IBC_TX_BATCH_LIMIT || 100),
        maxPageSize: Number(MAX_PAGE_SIZE || 100),
    },
    taskCfg:{
        interval:{
            heartbeatRate:Number(HEARTBEAT_RATE || 10000),
        },
        executeTime:{
            // tx: '*/10 * * * * *',
            // chain: '*/10 * * * * *',
            // statistics: '*/10 * * * * *',
            // tx: IBCTX_EXECUTE_TIME || '15 * * * * *',
            chain: IBCCHAIN_EXECUTE_TIME || '0 0 */1 * * *',
            statistics: IBCSTATISTICS_EXECUTE_TIME || '0 */10 * * * *',

            faultTolerance:FAULT_TOLERANCE_EXECUTE_TIME || '41 * * * * *',
            transferTx: SYNC_TRANSFER_TX_TIME || '*/15 * * * * *',
            updateProcessingTx: UPDATE_PROCESSING_TX_TIME || '*/15 * * * * *',
            updateSubStateTx: UPDATE_SUB_STATE_TX_TIME || '*/15 * * * * *',
            ibcTxLatestMigrate: IBC_TX_LATEST_MIGRATE || '0 */30 * * * *',
            ibcTxUpdateCronjob: IBC_TX_UPDATE_CRONJOB || '0 0 */2 * * *',
            ibcDenomCaculateCronjob: IBC_DENOM_CACULATE_CRONJOB || '0 0 */1 * * *',
            ibcDenomUpdateCronjob: IBC_DENOM_UPDATE_CRONJOB || '*/30 * * * * *',
            ibcMonitorCronjob: IBC_MONITOR_CRONJOB || '0 */1 * * * *',
        },
        CRON_JOBS: CRON_JOBS ? JSON.parse(CRON_JOBS) : [],
    },
};

