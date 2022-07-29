package vo

type SupplyResp struct {
	Supply []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"supply"`
	Pagination struct {
		NextKey *string `json:"next_key"`
		Total   string  `json:"total"`
	} `json:"pagination"`
}

type BalancesResp struct {
	Balances []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"balances"`
	Pagination struct {
		NextKey *string `json:"next_key"`
		Total   string  `json:"total"`
	} `json:"pagination"`
}

type StakeParams struct {
	Params struct {
		UnbondingTime string `json:"unbonding_time"`
	} `json:"params"`
}

type Connections struct {
	ConnectionPaths []string `json:"connection_paths"`
}

type (
	ConnectionChannels struct {
		Channels []LcdChannel `json:"channels"`
	}
	LcdChannel struct {
		State        string `json:"state"`
		Counterparty struct {
			PortId    string `json:"port_id"`
			ChannelId string `json:"channel_id"`
		} `json:"counterparty"`
		ConnectionHops []string `json:"connection_hops"`
		PortId         string   `json:"port_id"`
		ChannelId      string   `json:"channel_id"`
	}
)
