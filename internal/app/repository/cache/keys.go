package cache

import (
	"time"
)

// redis key expiration
const (
	oneHour                    = 1 * time.Hour
	threeHours                 = 3 * time.Hour
	oneDay                     = 24 * time.Hour
	oneMin                     = 60 * time.Second
	FiveMin                    = 5 * time.Minute
	NoExpiration time.Duration = -1
)

// redis key
const (
	tokenPrice                  = "token_price"
	denomSupply                 = "denom_supply:%s"
	denomTransAmount            = "denom_trans_amount:%s"
	ibcInfoHash                 = "ibc_info_hash"
	ibcInfo                     = "ibc_info:%s"
	lcdTxData                   = "lcd_tx_data:%s"
	lcdAddr                     = "lcd_addr:%s"
	clientIdInfo                = "client_id_info:%s"
	ibcRelayerCfg               = "ibc_relayer_config"
	ibcRelayerCfgPairIds        = "ibc_relayer_config_pair_ids"
	ibcChainsConnection         = "ibc_chains_connection"
	ibcRelayerTotalTxs          = "relayer_total_txs"
	ibcRelayerTransferTypeTxs   = "relayer_transfer_type_txs"
	ibcRelayerTotalRelayedValue = "relayer_total_relayed_value"
	ibcRelayerTotalFeeCost      = "relayer_total_fee_cost"
	relayerRelayedTrend         = "relayer_relayed_trend"
	ibcRelayer                  = "ibc_relayer"
	baseDenom                   = "base_denom"
	ibcChain                    = "ibc_chain"
	chainUnbondTime             = "chain_unbond_time"
	statisticsCheck             = "statistics_check:%s_%s"
	BaseDenomUnauth             = "base_denom_unauth"
	baseDenomSymbol             = "base_denom:%s"
	clientState                 = "client_state:%s"
	DenomSymbol                 = "denom_symbol"
	chainInflowVolume           = "chain_inflow_volume_%d"
	chainOutflowVolume          = "chain_outflow_volume_%d"
)

// Keys of open api
// The prefix is "oa"
const (
	openApiKey = "oa_open_api_key:%s"
	rateLimit  = "oa_rate_limit:%s"
)
