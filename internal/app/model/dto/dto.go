package dto

import (
	"github.com/shopspring/decimal"
	"math"
)

type CountBaseDenomTxsDTO struct {
	BaseDenom      string `bson:"base_denom"`
	BaseDenomChain string `bson:"base_denom_chain"`
	Count          int64  `bson:"count"`
}

type GetDenomGroupByChainDTO struct {
	Chain string   `bson:"_id"`
	Denom []string `bson:"denom"`
}

type GetBaseDenomFromIbcDenomDTO struct {
	BaseDenom string `bson:"_id"`
}

type CountIBCTokenRecvTxsDTO struct {
	Denom string `bson:"denom"`
	Chain string `bson:"chain"`
	Count int64  `bson:"count"`
}

type AggregateIBCChainDTO struct {
	Chain      string  `bson:"_id"`
	DenomValue float64 `bson:"denom_value"`
	Count      int64   `bson:"count"`
}

type RelayerDenomStatisticsDTO struct {
	Signer      string  `bson:"signer"`
	Status      int64   `bson:"status"`
	TxType      string  `bson:"tx_type"`
	Denom       string  `bson:"denom"`
	ScChannel   string  `bson:"sc_channel"`
	DcChannel   string  `bson:"dc_channel"`
	DenomAmount float64 `bson:"denom_amount"`
	TxsCount    int64   `bson:"txs_count"`
}

type RelayerFeeStatisticsDTO struct {
	Signer      string  `bson:"signer"`
	Status      int64   `bson:"status"`
	TxType      string  `bson:"tx_type"`
	Denom       string  `bson:"denom"`
	DenomAmount float64 `bson:"denom_amount"`
	TxsCount    int64   `bson:"txs_count"`
}

type GetRelayerInfoDTO struct {
	ScChainAddress string `bson:"sc_chain_address"`
	DcChainAddress string `bson:"dc_chain_address"`
	ScChain        string `bson:"sc_chain"`
	ScChannel      string `bson:"sc_channel"`
	DcChain        string `bson:"dc_chain"`
	DcChannel      string `bson:"dc_channel"`
}

type CountRelayerPacketTxsCntDTO struct {
	DcChainAddress   string `bson:"dc_chain_address"`
	ScChainId        string `bson:"sc_chain_id"`
	ScChannel        string `bson:"sc_channel"`
	DcChainId        string `bson:"dc_chain_id"`
	DcChannel        string `bson:"dc_channel"`
	BaseDenom        string `bson:"base_denom"`
	BaseDenomChainId string `bson:"base_denom_chain_id"`
	Count            int64  `bson:"count"`
}

func (dto *CountRelayerPacketTxsCntDTO) Valid() bool {
	return dto.DcChainId != "" && dto.DcChannel != "" && dto.ScChainId != "" && dto.ScChannel != "" && dto.DcChainAddress != ""
}

type CountRelayerPacketAmountDTO struct {
	DcChainAddress   string  `bson:"dc_chain_address"`
	DcChainId        string  `bson:"dc_chain_id"`
	DcChannel        string  `bson:"dc_channel"`
	ScChainId        string  `bson:"sc_chain_id"`
	ScChannel        string  `bson:"sc_channel"`
	BaseDenom        string  `bson:"base_denom"`
	BaseDenomChainId string  `bson:"base_denom_chain_id"`
	Amount           float64 `bson:"amount"`
	Count            int64   `bson:"count"`
}

func (dto *CountRelayerPacketAmountDTO) Valid() bool {
	return dto.DcChainId != "" && dto.DcChannel != "" && dto.ScChainId != "" && dto.ScChannel != "" && dto.BaseDenom != "" && dto.DcChainAddress != ""
}

type CountRelayerBaseDenomAmtDTO struct {
	BaseDenom      string  `bson:"base_denom"`
	BaseDenomChain string  `bson:"base_denom_chain"`
	TxStatus       int     `bson:"tx_status"`
	Amount         float64 `bson:"amount"`
	TotalTxs       int64   `bson:"total_txs"`
}

type AggrRelayerTxTypeDTO struct {
	TxType   string `bson:"tx_type"`
	TotalTxs int64  `bson:"total_txs"`
}

type CountRelayerBaseDenomAmtBySegmentDTO struct {
	BaseDenom        string  `bson:"base_denom"`
	BaseDenomChain   string  `bson:"base_denom_chain"`
	SegmentStartTime int64   `bson:"segment_start_time"`
	Amount           float64 `bson:"amount"`
	TotalTxs         int64   `bson:"total_txs"`
}

type AggrRelayerTxsAmtDTo struct {
	FeeDenom string  `bson:"fee_denom"`
	ChainId  string  `bson:"chain_id"`
	Amount   float64 `bson:"amount"`
	TotalTxs int64   `bson:"total_txs"`
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
	BaseDenom      string  `bson:"base_denom"`
	BaseDenomChain string  `bson:"base_denom_chain"`
	ScChain        string  `bson:"sc_chain"`
	DcChain        string  `bson:"dc_chain"`
	ScChannel      string  `bson:"sc_channel"`
	DcChannel      string  `bson:"dc_channel"`
	Count          int64   `bson:"count"`
	Amount         float64 `bson:"amount"`
}

type ChannelStatisticsDTO struct {
	ChannelId      string          `bson:"channel_id"`
	BaseDenom      string          `bson:"base_denom"`
	BaseDenomChain string          `bson:"base_denom_chain"`
	TxsCount       int64           `bson:"count"`
	TxsAmount      decimal.Decimal `bson:"amount"`
}

type ChannelStatisticsAggrDTO struct {
	ChannelId      string  `bson:"channel_id"`
	BaseDenom      string  `bson:"base_denom"`
	BaseDenomChain string  `bson:"base_denom_chain"`
	TxsCount       int64   `bson:"count"`
	TxsAmount      float64 `bson:"amount"`
}

type TokenTraceStatisticsDTO struct {
	Denom      string `bson:"denom"`
	Chain      string `bson:"chain"`
	ReceiveTxs int64  `bson:"receive_txs"`
}

type Aggr24hActiveChannelsDTO struct {
	ScChainId string `bson:"sc_chain_id"`
	DcChainId string `bson:"dc_chain_id"`
	ScChannel string `bson:"sc_channel"`
	DcChannel string `bson:"dc_channel"`
}

type Aggr24hActiveChainsDTO struct {
	ScChainId string `bson:"sc_chain_id"`
	DcChainId string `bson:"dc_chain_id"`
}

type RelayerPairIdDTO struct {
	RelayerPairId string `bson:"relayer_pair_id"`
}

type DenomSimpleDTO struct {
	Denom   string
	ChainId string
}

type PacketIdDTO struct {
	RecordId      string
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
	BaseDenom        []string
	BaseDenomChainId string
	Denom            string
}

type (
	TxsAmtItem struct {
		Txs        int64
		TxsSuccess int64
		Denom      string
		ChainId    string
		Amt        decimal.Decimal
		AmtValue   decimal.Decimal
	}

	CoinItem struct {
		Price float64
		Scale int
	}
)

func CaculateRelayerTotalValue(denomPriceMap map[string]CoinItem, relayerTxsDataMap map[string]TxsAmtItem) decimal.Decimal {
	totalValue := decimal.NewFromFloat(0)

	for key, data := range relayerTxsDataMap {
		baseDenomValue := decimal.NewFromFloat(0)
		decAmt := data.Amt
		if coin, ok := denomPriceMap[key]; ok {
			if coin.Scale > 0 {
				baseDenomValue = decAmt.Div(decimal.NewFromFloat(math.Pow10(coin.Scale))).Mul(decimal.NewFromFloat(coin.Price))
				data.AmtValue = baseDenomValue
				relayerTxsDataMap[key] = data
			}
		}
		totalValue = totalValue.Add(baseDenomValue)
	}
	return totalValue
}

type BaseDenomAmountDTO struct {
	BaseDenom      string  `bson:"base_denom"`
	BaseDenomChain string  `bson:"base_denom_chain"`
	Amount         float64 `bson:"amount"`
}
