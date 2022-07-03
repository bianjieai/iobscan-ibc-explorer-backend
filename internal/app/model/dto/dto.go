package dto

import "github.com/shopspring/decimal"

type CountBaseDenomTransferAmountDTO struct {
	BaseDenom string `bson:"_id"`
	Count     int64  `bson:"count"`
}

type GetDenomGroupByChainIdDTO struct {
	ChainId string   `bson:"_id"`
	Denom   []string `bson:"denom"`
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

type CountRelayerPacketTxsCntDTO struct {
	DcChainAddress string `bson:"dc_chain_address"`
	DcChainId      string `bson:"dc_chain_id"`
	DcChannel      string `bson:"dc_channel"`
	Count          int64  `bson:"count"`
}

type CountRelayerPacketAmountDTO struct {
	DcChainAddress string  `bson:"dc_chain_address"`
	DcChainId      string  `bson:"dc_chain_id"`
	DcChannel      string  `bson:"dc_channel"`
	BaseDenom      string  `bson:"base_denom"`
	Amount         float64 `bson:"amount"`
}
type CountRelayerTotalValueDTO struct {
	RelayerId string  `bson:"relayer_id"`
	ChainId   string  `bson:"chain_id"`
	Channel   string  `bson:"channel"`
	Amount    float64 `bson:"amount"`
}

type CountChannelRelayersDTO struct {
	ChainA   string `bson:"chain_a"`
	ChannelA string `bson:"channel_a"`
	ChainB   string `bson:"chain_b"`
	ChannelB string `bson:"channel_b"`
	Count    int64  `bson:"count"`
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
	ChannelId string          `bson:"channel_id"`
	BaseDenom string          `bson:"base_denom"`
	TxsCount  int64           `bson:"count"`
	TxsAmount decimal.Decimal `bson:"amount"`
}

type ChannelStatisticsAggrDTO struct {
	ChannelId string  `bson:"channel_id"`
	BaseDenom string  `bson:"base_denom"`
	TxsCount  int64   `bson:"count"`
	TxsAmount float64 `bson:"amount"`
}
