package vo

type BaseInfoResp struct {
	Address         string `json:"address"`
	PubKey          string `json:"pub_key"`
	Chain           string `json:"chain"`
	AccountNumber   string `json:"account_number"`
	PubKeyType      string `json:"pub_key_type"`
	PubKeyAlgorithm string `json:"pub_key_algorithm"`
	AccountSequence string `json:"account_sequence"`
}
