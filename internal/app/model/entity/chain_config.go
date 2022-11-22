package entity

type ChainStatus int

type (
	ChainConfig struct {
		CurrentChainId  string      `bson:"current_chain_id"`
		Icon            string      `bson:"icon"`
		ChainName       string      `bson:"chain_name"`
		PrettyName      string      `bson:"pretty_name"`
		LcdApiPath      ApiPath     `bson:"lcd_api_path"`
		GrpcRestGateway string      `bson:"grpc_rest_gateway"`
		AddrPrefix      string      `bson:"addr_prefix"`
		IbcInfo         []*IbcInfo  `bson:"ibc_info"`
		IbcInfoHashLcd  string      `bson:"ibc_info_hash_lcd"`
		Status          ChainStatus `bson:"status"`
	}
	ApiPath struct {
		ChannelsPath    string `bson:"channels_path"`
		ClientStatePath string `bson:"client_state_path"`
		SupplyPath      string `bson:"supply_path"`
		BalancesPath    string `bson:"balances_path"`
		ParamsPath      string `bson:"params_path"`
	}
	IbcInfo struct {
		Chain string         `bson:"chain"`
		Paths []*ChannelPath `bson:"paths"`
	}
	ChannelPath struct {
		State        string       `bson:"state"`
		PortId       string       `bson:"port_id"`
		ChannelId    string       `bson:"channel_id"`
		Chain        string       `bson:"chain"`
		ScChain      string       `bson:"sc_chain"`
		ClientId     string       `bson:"client_id"`
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
