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

type AggrChainAddrDTO struct {
	Chain   string `bson:"chain"`
	Address string `bson:"address"`
}

type CountRelayerBaseDenomAmtBySegmentDTO struct {
	BaseDenom        string  `bson:"base_denom"`
	BaseDenomChain   string  `bson:"base_denom_chain"`
	SegmentStartTime int64   `bson:"segment_start_time"`
	Amount           float64 `bson:"amount"`
	TotalTxs         int64   `bson:"total_txs"`
}

type AggrChainOutflowTrendDTO struct {
	BaseDenom        string  `bson:"base_denom"`
	BaseDenomChain   string  `bson:"base_denom_chain"`
	SegmentStartTime int64   `bson:"segment_start_time"`
	DenomAmount      float64 `bson:"denom_amount"`
}

type AggrDenomTxsDTO struct {
	BaseDenom      string `bson:"base_denom"`
	BaseDenomChain string `bson:"base_denom_chain"`
	TxsNumber      int64  `bson:"txs_number"`
}

type AggrChainInflowTrendDTO struct {
	BaseDenom        string `bson:"base_denom"`
	BaseDenomChain   string `bson:"base_denom_chain"`
	SegmentStartTime int64  `bson:"segment_start_time"`
	//TxsNumber        int64   `bson:"txs_number"`
	DenomAmount float64 `bson:"denom_amount"`
}

type AggrRelayerTxsAmtDTo struct {
	FeeDenom string  `bson:"fee_denom"`
	Chain    string  `bson:"chain"`
	Amount   float64 `bson:"amount"`
	TotalTxs int64   `bson:"total_txs"`
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
	Status         int64   `bson:"status"`
}

type AggrIBCChainInflowDTO struct {
	BaseDenom      string  `bson:"base_denom"`
	BaseDenomChain string  `bson:"base_denom_chain"`
	Chain          string  `bson:"chain"`
	Status         int64   `bson:"status"`
	TxsNum         int64   `bson:"txs_num"`
	DenomAmount    float64 `bson:"denom_amount"`
}

type AggrIBCChainOutflowDTO struct {
	BaseDenom      string  `bson:"base_denom"`
	BaseDenomChain string  `bson:"base_denom_chain"`
	Chain          string  `bson:"chain"`
	Status         int64   `bson:"status"`
	TxsNum         int64   `bson:"txs_num"`
	DenomAmount    float64 `bson:"denom_amount"`
}

type ChannelStatisticsDTO struct {
	ChannelId      string          `bson:"channel_id"`
	BaseDenom      string          `bson:"base_denom"`
	BaseDenomChain string          `bson:"base_denom_chain"`
	TxsCount       int64           `bson:"count"`
	TxsAmount      decimal.Decimal `bson:"amount"`
	Status         int64           `bson:"status"`
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
	ScChain   string `bson:"sc_chain"`
	DcChain   string `bson:"dc_chain"`
	ScChannel string `bson:"sc_channel"`
	DcChannel string `bson:"dc_channel"`
}

type Aggr24hActiveChainsDTO struct {
	ScChain string `bson:"sc_chain"`
	DcChain string `bson:"dc_chain"`
}

type Aggr24hDenomVolumeDTO struct {
	BaseDenom      string  `bson:"base_denom"`
	BaseDenomChain string  `bson:"base_denom_chain"`
	DenomAmount    float64 `bson:"denom_amount"`
}

type DenomSimpleDTO struct {
	Denom string
	Chain string
}

type PacketIdDTO struct {
	ObjectId      string
	DcChain       string
	TimeoutHeight int64
	PacketId      string
	TimeOutTime   int64
}

type HeightTimeDTO struct {
	Height int64
	Time   int64
}

type IbcTxQuery struct {
	StartTime      int64
	EndTime        int64
	Chain          []string
	Status         []int
	BaseDenom      []string
	BaseDenomChain string
	Denom          string
}

type (
	TxsAmtItem struct {
		Txs        int64
		TxsSuccess int64
		Denom      string
		Chain      string
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

type ChainAddressDTO struct {
	Chain   string `bson:"chain"`
	Address string `bson:"address"`
}
