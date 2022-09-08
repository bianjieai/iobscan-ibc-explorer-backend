package vo

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
)

type RelayerListReq struct {
	Page
	Chain    string `json:"chain" form:"chain"`
	Status   int    `json:"status" form:"status"`
	UseCount bool   `json:"use_count" form:"use_count"`
}

type RelayerDto struct {
	RelayerId             string `json:"relayer_id"`
	RelayerName           string `json:"relayer_name"`
	RelayerIcon           string `json:"relayer_icon"`
	ChainA                string `json:"chain_a"`
	ChainB                string `json:"chain_b"`
	ChannelA              string `json:"channel_a"`
	ChannelB              string `json:"channel_b"`
	ChainAAddress         string `json:"chain_a_address"`
	ChainBAddress         string `json:"chain_b_address"`
	TimePeriod            int64  `json:"time_period"`
	Status                int    `json:"status"`
	UpdateTime            int64  `json:"update_time"`
	TransferTotalTxs      int64  `json:"transfer_total_txs"`
	TransferSuccessTxs    int64  `json:"transfer_success_txs"`
	TransferTotalTxsValue string `json:"transfer_total_txs_value"`
	Currency              string `json:"currency"`
}

type RelayerListResp struct {
	Items     []RelayerDto `json:"items"`
	PageInfo  PageInfo     `json:"page_info"`
	TimeStamp int64        `json:"time_stamp"`
}

func (dto RelayerDto) LoadDto(relayer *entity.IBCRelayer) RelayerDto {
	return RelayerDto{
		RelayerId:             relayer.RelayerId,
		ChainA:                relayer.ChainA,
		ChainB:                relayer.ChainB,
		ChannelA:              relayer.ChannelA,
		ChannelB:              relayer.ChannelB,
		ChainAAddress:         relayer.ChainAAddress,
		ChainBAddress:         relayer.ChainBAddress,
		TimePeriod:            relayer.TimePeriod,
		Status:                int(relayer.Status),
		UpdateTime:            relayer.UpdateTime,
		TransferTotalTxs:      relayer.TransferTotalTxs,
		TransferSuccessTxs:    relayer.TransferSuccessTxs,
		TransferTotalTxsValue: relayer.TransferTotalTxsValue,
		Currency:              constant.DefaultCurrency,
	}
}

type IobRegistryRelayerInfoResp struct {
	TeamName string `json:"team name"`
	Contact  struct {
		Website  string `json:"website"`
		Github   string `json:"github"`
		Telegram string `json:"telegram"`
		Twitter  string `json:"twitter"`
		Medium   string `json:"medium"`
		Discord  string `json:"discord"`
	} `json:"contact"`
	Introduction []string `json:"introduction"`
}

type IobRegistryRelayerPairResp struct {
	Chain1 struct {
		Address   string `json:"address"`
		ChainId   string `json:"chain-id"`
		ChannelId string `json:"channel-id"`
		Version   string `json:"version"`
	} `json:"chain-1"`
	Chain2 struct {
		Address   string `json:"address"`
		ChainId   string `json:"chain-id"`
		ChannelId string `json:"channel-id"`
		Version   string `json:"version"`
	} `json:"chain-2"`
}
