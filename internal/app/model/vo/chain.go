package vo

type ChainListResp struct {
	Items   []ChainItem `json:"items"`
	Comment string      `json:"comment"`
}

type ChainItem struct {
	Chain            string `json:"chain"`
	ChainRegistryUrl string `json:"chain_registry_url"`
}
