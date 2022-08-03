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

type IbcChannelsResp struct {
	Channels []struct {
		State        string `json:"state"`
		Ordering     string `json:"ordering"`
		Counterparty struct {
			PortId    string `json:"port_id"`
			ChannelId string `json:"channel_id"`
		} `json:"counterparty"`
		ConnectionHops []string `json:"connection_hops"`
		Version        string   `json:"version"`
		PortId         string   `json:"port_id"`
		ChannelId      string   `json:"channel_id"`
	} `json:"channels"`
	Pagination struct {
		NextKey *string `json:"next_key"`
		Total   string  `json:"total"`
	} `json:"pagination"`
	Height struct {
		RevisionNumber string `json:"revision_number"`
		RevisionHeight string `json:"revision_height"`
	} `json:"height"`
}

type ClientStateResp struct {
	IdentifiedClientState struct {
		ClientId    string `json:"client_id"`
		ClientState struct {
			Type       string `json:"@type"`
			ChainId    string `json:"chain_id"`
			TrustLevel struct {
				Numerator   string `json:"numerator"`
				Denominator string `json:"denominator"`
			} `json:"trust_level"`
			TrustingPeriod  string `json:"trusting_period"`
			UnbondingPeriod string `json:"unbonding_period"`
			MaxClockDrift   string `json:"max_clock_drift"`
			FrozenHeight    struct {
				RevisionNumber string `json:"revision_number"`
				RevisionHeight string `json:"revision_height"`
			} `json:"frozen_height"`
			LatestHeight struct {
				RevisionNumber string `json:"revision_number"`
				RevisionHeight string `json:"revision_height"`
			} `json:"latest_height"`
			ProofSpecs []struct {
				LeafSpec struct {
					Hash         string `json:"hash"`
					PrehashKey   string `json:"prehash_key"`
					PrehashValue string `json:"prehash_value"`
					Length       string `json:"length"`
					Prefix       string `json:"prefix"`
				} `json:"leaf_spec"`
				InnerSpec struct {
					ChildOrder      []int       `json:"child_order"`
					ChildSize       int         `json:"child_size"`
					MinPrefixLength int         `json:"min_prefix_length"`
					MaxPrefixLength int         `json:"max_prefix_length"`
					EmptyChild      interface{} `json:"empty_child"`
					Hash            string      `json:"hash"`
				} `json:"inner_spec"`
				MaxDepth int `json:"max_depth"`
				MinDepth int `json:"min_depth"`
			} `json:"proof_specs"`
			UpgradePath                  []string `json:"upgrade_path"`
			AllowUpdateAfterExpiry       bool     `json:"allow_update_after_expiry"`
			AllowUpdateAfterMisbehaviour bool     `json:"allow_update_after_misbehaviour"`
		} `json:"client_state"`
	} `json:"identified_client_state"`
	Proof       interface{} `json:"proof"`
	ProofHeight struct {
		RevisionNumber string `json:"revision_number"`
		RevisionHeight string `json:"revision_height"`
	} `json:"proof_height"`
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