package dto

import "github.com/shopspring/decimal"

type CountBaseDenomTxsDTO struct {
	BaseDenom string `bson:"_id"`
	Count     int64  `bson:"count"`
}

type GetDenomGroupByChainIdDTO struct {
	ChainId string   `bson:"_id"`
	Denom   []string `bson:"denom"`
}

type GetBaseDenomFromIbcDenomDTO struct {
	BaseDenom string `bson:"_id"`
}

type CountIBCTokenRecvTxsDTO struct {
	BaseDenom string `bson:"base_denom"`
	Denom     string `bson:"denom"`
	ChainId   string `bson:"chain_id"`
	Count     int64  `bson:"count"`
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
	ScChainId      string `bson:"sc_chain_id"`
	ScChannel      string `bson:"sc_channel"`
	DcChainId      string `bson:"dc_chain_id"`
	DcChannel      string `bson:"dc_channel"`
	BaseDenom      string `bson:"base_denom"`
	Count          int64  `bson:"count"`
}

func (dto *CountRelayerPacketTxsCntDTO) Valid() bool {
	return dto.DcChainId != "" && dto.DcChannel != "" && dto.ScChainId != "" && dto.ScChannel != "" && dto.DcChainAddress != ""
}

type CountRelayerPacketAmountDTO struct {
	DcChainAddress string  `bson:"dc_chain_address"`
	DcChainId      string  `bson:"dc_chain_id"`
	DcChannel      string  `bson:"dc_channel"`
	ScChainId      string  `bson:"sc_chain_id"`
	ScChannel      string  `bson:"sc_channel"`
	BaseDenom      string  `bson:"base_denom"`
	Amount         float64 `bson:"amount"`
	Count          int64   `bson:"count"`
}

func (dto *CountRelayerPacketAmountDTO) Valid() bool {
	return dto.DcChainId != "" && dto.DcChannel != "" && dto.ScChainId != "" && dto.ScChannel != "" && dto.BaseDenom != "" && dto.DcChainAddress != ""
}

type CountRelayerBaseDenomAmtDTO struct {
	StatisticId string  `bson:"statistic_id"`
	Address     string  `bson:"address"`
	Amount      float64 `bson:"amount"`
	BaseDenom   string  `bson:"base_denom"`
}

type AggRelayerTxsDTO struct {
	Address         string `bson:"address"`
	StatisticId     string `bson:"statistic_id"`
	SuccessTotalTxs int64  `bson:"success_total_txs"`
	TotalTxs        int64  `bson:"total_txs"`
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

type TokenTraceStatisticsDTO struct {
	Denom      string `bson:"denom"`
	ChainId    string `bson:"chain_id"`
	ReceiveTxs int64  `bson:"receive_txs"`
}

type Aggr24hActiveChannelTxsDTO struct {
	ScChainId string `bson:"sc_chain_id"`
	DcChainId string `bson:"dc_chain_id"`
	ScChannel string `bson:"sc_channel"`
	DcChannel string `bson:"dc_channel"`
}

type DenomSimpleDTO struct {
	Denom   string
	ChainId string
}

type PacketIdDTO struct {
	DcChainId     string
	TimeoutHeight int64
	PacketId      string
	TimeOutTime   int64
}

type HeightTimeDTO struct {
	Height int64
	Time   int64
}

type IbcTxQuery struct {
	StartTime        int64
	EndTime          int64
	ChainId          []string
	Status           []int
	Token            []string
	BaseDenomChainId string
}
