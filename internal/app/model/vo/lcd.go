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

type UnbondingResp struct {
	UnbondingResponses []struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Entries          []struct {
			CreationHeight string `json:"creation_height"`
			InitialBalance string `json:"initial_balance"`
			Balance        string `json:"balance"`
		} `json:"entries"`
	} `json:"unbonding_responses"`
	Pagination struct {
		NextKey *string `json:"next_key"`
		Total   string  `json:"total"`
	} `json:"pagination"`
}

type DelegationResp struct {
	DelegationResponses []struct {
		Delegation struct {
			DelegatorAddress string `json:"delegator_address"`
			ValidatorAddress string `json:"validator_address"`
			Shares           string `json:"shares"`
		} `json:"delegation"`
		Balance struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"balance"`
	} `json:"delegation_responses"`
	Pagination struct {
		NextKey *string `json:"next_key"`
		Total   string  `json:"total"`
	} `json:"pagination"`
}

type RewardsResp struct {
	Rewards []struct {
		ValidatorAddress string `json:"validator_address"`
		Reward           []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"reward"`
	} `json:"rewards"`
	Total []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"total"`
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

type ChainRegisterResp struct {
	Schema       string `json:"$schema"`
	ChainName    string `json:"chain_name"`
	Status       string `json:"status"`
	NetworkType  string `json:"network_type"`
	PrettyName   string `json:"pretty_name"`
	ChainId      string `json:"chain_id"`
	Bech32Prefix string `json:"bech32_prefix"`
	DaemonName   string `json:"daemon_name"`
	NodeHome     string `json:"node_home"`
	Apis         struct {
		Rpc []struct {
			Address  string `json:"address"`
			Provider string `json:"provider"`
		} `json:"rpc"`
		Rest []struct {
			Address  string `json:"address"`
			Provider string `json:"provider"`
		} `json:"rest"`
		Grpc []struct {
			Address  string `json:"address"`
			Provider string `json:"provider"`
		} `json:"grpc"`
	} `json:"apis"`
}

type StatusResp struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		NodeInfo struct {
			Other struct {
				TxIndex    string `json:"tx_index"`
				RPCAddress string `json:"rpc_address"`
			} `json:"other"`
		} `json:"node_info"`
		SyncInfo struct {
			LatestBlockHeight   int64 `json:"latest_block_height,string"`
			EarliestBlockHeight int64 `json:"earliest_block_height,string"`
			CatchingUp          bool  `json:"catching_up"`
		} `json:"sync_info"`
	} `json:"result"`
}

type AccountResp struct {
	Account struct {
		Type    string `json:"@type"`
		Address string `json:"address"`
		PubKey  struct {
			Type string `json:"@type"`
			Key  string `json:"key"`
		} `json:"pub_key"`
		AccountNumber string `json:"account_number"`
		Sequence      string `json:"sequence"`
	} `json:"account"`
}
