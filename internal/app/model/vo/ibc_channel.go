package vo

import "github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"

type ChannelListResp struct {
	Items    []ChannelItem `json:"items"`
	PageInfo PageInfo      `json:"page_info"`
}

type ChannelItem struct {
	ChainA              string               `json:"chain_a"`
	ChannelA            string               `json:"channel_a"`
	ChainB              string               `json:"chain_b"`
	ChannelB            string               `json:"channel_b"`
	OperatingPeriod     int64                `json:"operating_period"`
	Relayers            int                  `json:"relayers"`
	LastUpdated         int64                `json:"last_updated"`
	IbcTransferTxsValue string               `json:"ibc_transfer_txs_value"`
	IbcTransferTxs      int64                `json:"ibc_transfer_txs"`
	Currency            string               `json:"currency"`
	Status              entity.ChannelStatus `json:"status"`
}
