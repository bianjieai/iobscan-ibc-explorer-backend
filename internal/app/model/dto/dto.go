package dto

import "github.com/shopspring/decimal"

type CountBaseDenomTransferAmountDTO struct {
	BaseDenom string `bson:"base_denom"`
	ScChainId string `bson:"sc_chain_id"`
	DcChainId string `bson:"dc_chain_id"`
	Count     int64  `bson:"count"`
}

type GetDenomGroupByBaseDenomDTO struct {
	BaseDenom string   `bson:"_id"`
	Denom     []string `bson:"denom"`
}

type CountIBCTokenRecvTxsDTO struct {
	Denom string `bson:"_id"`
	Count int64  `bson:"count"`
}

type AggregateIBCChainDTO struct {
	ChainId    string  `bson:"_id"`
	DenomValue float64 `bson:"denom_value"`
	Count      int64   `bson:"count"`
}

type GetRelayerInfoDTO struct {
	DcChainAddress string `bson:"dc_chain_address"`
	ScChainId      string `bson:"sc_chain_id"`
	ScChannel      string `bson:"sc_channel"`
	DcChainId      string `bson:"dc_chain_id"`
	DcChannel      string `bson:"dc_channel"`
}

type GetRelayerScChainAddreeDTO struct {
	ScChainAddress string `bson:"sc_chain_address"`
}

type AggrIBCChannelTxsDTO struct {
	BaseDenom string  `bson:"base_denom"`
	ScChainId string  `bson:"sc_chain_id"`
	DcChainId string  `bson:"dc_chain_id"`
	ScChannel string  `bson:"sc_channel"`
	DcChannel string  `bson:"dc_channel"`
	Count     int64   `bson:"count"`
	Amount    float64 `bson:"amount"`
}

type ChannelStatisticsDTO struct {
	ChannelId       string
	MirrorChannelId string
	BaseDenom       string
	TxsCount        int64
	TxsAmount       decimal.Decimal
}
