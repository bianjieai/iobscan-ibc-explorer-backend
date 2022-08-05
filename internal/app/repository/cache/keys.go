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
	tokenPrice       = "token_price"
	denomSupply      = "denom_supply:%s"
	denomTransAmount = "denom_trans_amount:%s"
	ibcInfoHash      = "ibc_info_hash"
	ibcInfo          = "ibc_info:%s"
	lcdInfo          = "lcd_info"
	clientIdInfo     = "client_id_info:%s"
	ibcRelayerCfg    = "ibc_relayer_config"
	ibcRelayer       = "ibc_relayer"
	baseDenom        = "base_denom"
	ibcChain         = "ibc_chain"
	chainUnbondTime  = "chain_unbond_time"
	statisticsCheck  = "statistics_check:%s_%s"

	BaseDenomUnauth = "base_denom_unauth"
	baseDenomSymbol = "base_denom:%s"
)
