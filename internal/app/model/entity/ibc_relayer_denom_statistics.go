package entity

import "fmt"

type TxType string

var ICS20TransferTxTypes = []TxType{TxTypeTransfer, TxTypeRecvPacket, TxTypeTimeoutPacket, TxTypeAckPacket}
var ICS20AllTxTypes = []TxType{TxTypeTransfer, TxTypeRecvPacket, TxTypeTimeoutPacket, TxTypeAckPacket, TxTypeUpdateClient}

const (
	TxTypeTransfer      TxType = "transfer"
	TxTypeRecvPacket    TxType = "recv_packet"
	TxTypeTimeoutPacket TxType = "timeout_packet"
	TxTypeAckPacket     TxType = "acknowledge_packet"
	TxTypeUpdateClient  TxType = "update_client"

	IBCRelayerDenomStatisticsCollName    = "ibc_relayer_denom_statistics"
	IBCRelayerDenomStatisticsNewCollName = "ibc_relayer_denom_statistics_new"
)

type IBCRelayerDenomStatistics struct {
	StatisticChain   string   `bson:"statistics_chain"`
	RelayerAddress   string   `bson:"relayer_address"`
	ChainAddressComb string   `bson:"chain_address_comb"`
	TxStatus         TxStatus `bson:"tx_status"`
	TxType           TxType   `bson:"tx_type"`
	BaseDenom        string   `bson:"base_denom"`
	BaseDenomChain   string   `bson:"base_denom_chain"`
	RelayedAmount    float64  `bson:"relayed_amount"`
	RelayedTxs       int64    `bson:"relayed_txs"`
	SegmentStartTime int64    `bson:"segment_start_time"`
	SegmentEndTime   int64    `bson:"segment_end_time"`
	CreateAt         int64    `bson:"create_at"`
	UpdateAt         int64    `bson:"update_at"`
}

func (i IBCRelayerDenomStatistics) CollectionName(isNew bool) string {
	if isNew {
		return IBCRelayerDenomStatisticsNewCollName
	}
	return IBCRelayerDenomStatisticsCollName
}

func GenerateChainAddressComb(chain, address string) string {
	return fmt.Sprintf("%s|%s", chain, address)
}
