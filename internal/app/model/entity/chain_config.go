package entity

const (
	ApiBalancesPathPlaceholder = "{address}"
)

type (
	ChainConfig struct {
		ChainId        string     `bson:"chain_id"`
		Icon           string     `bson:"icon"`
		ChainName      string     `bson:"chain_name"`
		LcdApiPath     ApiPath    `bson:"lcd_api_path"`
		Lcd            string     `bson:"lcd"`
		AddrPrefix     string     `bson:"addr_prefix"`
		IbcInfo        []*IbcInfo `bson:"ibc_info"`
		IbcInfoHashLcd string     `bson:"ibc_info_hash_lcd"`
	}
	ApiPath struct {
		ChannelsPath    string `bson:"channels_path"`
		ClientStatePath string `bson:"client_state_path"`
		SupplyPath      string `json:"supply_path"`
		BalancesPath    string `json:"balances_path"`
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
		Counterparty CounterParty `bson:"counterparty"`
	}
	CounterParty struct {
		State     string `bson:"state"`
		PortId    string `bson:"port_id"`
		ChannelId string `bson:"channel_id"`
	}
)

func (i ChainConfig) CollectionName() string {
	return "chain_config"
}
