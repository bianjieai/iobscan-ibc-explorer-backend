package vo

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
)

type ChainListReq struct {
	Page
	UseCount bool `json:"use_count" form:"use_count"`
}
type ChainDto struct {
	ChainId          string `json:"chain_id"`
	ConnectedChains  int64  `json:"connected_chains"`
	Channels         int64  `json:"channels"`
	Relayers         int64  `json:"relayers"`
	IbcTokens        int64  `json:"ibc_tokens"`
	IbcTokensValue   string `json:"ibc_tokens_value"`
	TransferTxs      int64  `json:"transfer_txs"`
	TransferTxsValue string `json:"transfer_txs_value"`
	Currency         string `json:"currency"`
}

type ChainListResp struct {
	Items    []ChainDto `json:"items"`
	PageInfo PageInfo   `json:"page_info"`
}

func (dto ChainDto) LoadDto(chain *entity.IBCChain) ChainDto {
	return ChainDto{
		ChainId:          chain.ChainId,
		ConnectedChains:  chain.ConnectedChains,
		Channels:         chain.Channels,
		Relayers:         chain.Relayers,
		IbcTokens:        chain.IbcTokens,
		IbcTokensValue:   chain.IbcTokensValue,
		TransferTxs:      chain.TransferTxs,
		TransferTxsValue: chain.TransferTxsValue,
		Currency:         constant.DefaultCurrency,
	}
}
