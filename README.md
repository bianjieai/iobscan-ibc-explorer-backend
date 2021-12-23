# iobscan-ibc-explorer-backend
IBC Explorer Backend

## Installation

```bash
$ npm install
```

npm run build

## environment configure:

NODE_ENV=development
DB_ADDR='192.168.150.40:27017'
DB_USER='ibc'
DB_PASSWD='ibcpassword'
DB_DATABASE='iobscan-ibc'
LCD_ADDR='http://192.168.150.40:1317'
CRON_JOBS = '["ex_sync_tx", "ex_sync_chain", "ex_sync_statistics"]'
IBCTX_EXECUTE_TIME = '15 * * * * *'
IBCCHAIN_EXECUTE_TIME = '30 * * * * *'
IBCSTATISTICS_EXECUTE_TIME = '45 * * * * *'
