package entity

import "github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"

const (
	ApiBalancesPathPlaceholder  = "{address}"
	ParamsModulePathPlaceholder = "{module}"
	StakeModule                 = "staking"
)

type ChainStatus int

const (
	ChainStatusOpen   ChainStatus = 1
	ChainStatusClosed ChainStatus = 2
)

type (
	ChainConfig struct {
		ChainId        string      `bson:"chain_id"`
		Icon           string      `bson:"icon"`
		ChainName      string      `bson:"chain_name"`
		LcdApiPath     ApiPath     `bson:"lcd_api_path"`
		Lcd            string      `bson:"lcd"`
		AddrPrefix     string      `bson:"addr_prefix"`
		IbcInfo        []*IbcInfo  `bson:"ibc_info"`
		IbcInfoHashLcd string      `bson:"ibc_info_hash_lcd"`
		Status         ChainStatus `bson:"status"`
	}
	ApiPath struct {
		ChannelsPath    string `bson:"channels_path"`
		ClientStatePath string `bson:"client_state_path"`
		SupplyPath      string `bson:"supply_path"`
		BalancesPath    string `bson:"balances_path"`
		ParamsPath      string `bson:"params_path"`
	}
	IbcInfo struct {
		ChainId string         `bson:"chain_id"`
		Paths   []*ChannelPath `bson:"paths"`
	}
	ChannelPath struct {
		State        string       `bson:"state"`
		PortId       string       `bson:"port_id"`
		ChannelId    string       `bson:"channel_id"`
		ChainId      string       `bson:"chain_id"`
		ScChainId    string       `bson:"sc_chain_id"`
		ClientId     string       `bson:"client_id"`
		Counterparty CounterParty `bson:"counterparty"`
	}
	CounterParty struct {
		State     string `bson:"state"`
		PortId    string `bson:"port_id"`
		ChannelId string `bson:"channel_id"`
	}
)

func (c ChainConfig) CollectionName() string {
	return "chain_config"
}

func (c *ChainConfig) GetChannelClient(port, channel string) string {
	if port == "" {
		port = constant.PortTransfer
	}
	if channel == "" {
		return ""
	}
	for _, ibcInfo := range c.IbcInfo {
		for _, path := range ibcInfo.Paths {
			if path.PortId == port && path.ChannelId == channel {
				return path.ClientId
			}
		}
	}

	return ""
}
