package vo

type ChainListResp struct {
	Items   []ChainItem `json:"items"`
	Comment string      `json:"comment"`
}

type ChainItem struct {
	Chain            string `json:"chain"`
	ChainRegistryUrl string `json:"chain_registry_url"`
}

type IbcChainsNumResp struct {
	IbcChainsNumber int64 `json:"ibc_chains_number"`
}
type IbcChainVolume struct {
	ChainName              string `json:"chain_name"`
	IbcTransferVolumeTotal string `json:"ibc_transfer_volume_total"`
	IbcVolumeIn            string `json:"ibc_volume_in"`
	IbcVolumeOut           string `json:"ibc_volume_out"`
}

type IbcChainsVolumeReq struct {
	Chain string `json:"chain" form:"chain"`
}

type IbcChainsVolumeResp struct {
	Chains    []IbcChainVolume `json:"chains"`
	TimeStamp int64            `json:"timestamp"`
}

type IbcChainsActiveResp struct {
	ChainNameList []string `json:"chain_name_list"`
}
