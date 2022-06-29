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
		UnbondingTime     string `json:"unbonding_time"`
	} `json:"params"`
}