## Description

<p align="center">IRITA service end</p>

## Installation

```bash
$ npm install
```

## Running the app

```bash
# development
$ npm run start

# production mode
$ npm run start:prod
```
## Env Variables

### Db config
- DB_ADDR: `required` `string` db addr（example: `127.0.0.1:27017, 127.0.0.2:27017, ...`）
- DB_USER: `required` `string` db user（example: `user`）
- DB_PASSWD: `required` `string` db password（example: `DB_PASSWD`）
- DB_DATABASE：`required` `string` database name（example：`DB_DATABASE`）

### Server config

- LCD_ADDR: `required` `string`  lcd address（example: `http://192.168.150.32:11317`）
- RPC_ADDR: `required` `string`  rpc address（example: `http://192.168.150.32:16657`）

### Task config

- HEARTBEAT_RATE: `Optional` `number`  hearbeat rate for monitor（example: `10000`）
- DENOM_EXECUTE_TIME: `Optional`  execute time for denom pull（example: "01 * * * * *"）
- NFT_EXECUTE_TIME: `Optional`  execute time for nft pull（example: "21 * * * * *"）
- TX_SERVICE_NAME_EXECUTE_TIME: `Optional`  execute time for nft pull（example: "30 * * * * *"）
- FAULT_TOLERANCE_EXECUTE_TIME: `Optional` `string`  execute time for fault tolerance（example: "41 * * * * *"）
- VALIDATORS_EXECTUTE_TIME `Optional` execute time for validators pull（example: "1 * * * * *"）
- IDENTITY_EXECUTE_TIME `Optional` execute time for identity pull（example: "1 * * * * *"）
- STAKING_VALIDATORS_INFO_TIME `Optional` execute time for stakingValidators info pull（example: "15 * * * * *"）
- STAKING_VALIDATORS_MORE_INFO_TIME `Optional` execute time for stakingValidators more info pull（example: "0 */5 * * * *"）
- STAKING_PARAMETERS `Optional` execute time for parameters pull（example: "10 * * * * *"）
- TOKENS `Optional` execute time for parameters pull（example: "5 * * * * *"）
- PROPLSAL `Optional` execute time for proplsal pull（example: "25 * * * * *"）
- ACCOUNT_EXECUTE_TIME `Optional` execute time for account pull（example: "35 * * * * *"）
- ACCOUNT_INFO_EXECUTE_TIME `Optional` execute time for account info pull（example: "* */10 * * * *"）
- SYNC_TX_SERVICE_NAME_SIZE: `Optional` `number`  execute time for fault tolerance（default 200）
- INCREASE_HEIGHT `Optional` `number` increase height for sync nft (default 1000)
- MAX_OPERATE_TX_COUNT `Optional` `number` limit operate tx count (default 100)
- CRON_JOBS `Optional` `array` name of the synchronization task performed (example: ["ex_sync_denom","ex_sync_nft","sync_tx_service_name","sync_validators","sync_identity","staking_sync_validators_info","staking_sync_validators_more_info","staking_sync_parameters","tokens","ex_sync_proposal"])
- PROPOSALS_LIMIT `Optional` `number` proposals limit for sync proposal (default 1000)

### Chain config
- CURRENT_CHAIN `Optional` `string` chain name (example: iris/cosmos/binance)

### log configure
- DisableLog: `Optional` `string` disable Log `true/false`
