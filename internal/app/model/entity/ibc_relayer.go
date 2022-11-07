package entity

import (
	"fmt"
	"strings"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
)

type RelayerStatus int

const (
	RelayerRunning RelayerStatus = 1
	RelayerStop    RelayerStatus = 2

	RelayerStopStr    = "Unknown"
	RelayerRunningStr = "Running"
)

type IBCRelayer struct {
	RelayerId             string        `bson:"relayer_id"`
	ChainA                string        `bson:"chain_a"`
	ChainB                string        `bson:"chain_b"`
	ChannelA              string        `bson:"channel_a"`
	ChannelB              string        `bson:"channel_b"`
	ChainAAddress         string        `bson:"chain_a_address"`
	ChainAAllAddress      []string      `bson:"chain_a_all_address"`
	ChainBAddress         string        `bson:"chain_b_address"`
	TimePeriod            int64         `bson:"time_period"`
	Status                RelayerStatus `bson:"status"`
	UpdateTime            int64         `bson:"update_time"`
	TransferTotalTxs      int64         `bson:"transfer_total_txs"`
	TransferSuccessTxs    int64         `bson:"transfer_success_txs"`
	TransferTotalTxsValue string        `bson:"transfer_total_txs_value"`
	CreateAt              int64         `bson:"create_at"`
	UpdateAt              int64         `bson:"update_at"`
}

func (i IBCRelayer) CollectionName() string {
	return "ibc_relayer"
}

func (i IBCRelayer) Valid() bool {
	return i.ChainA != "" && i.ChainB != "" && i.ChannelA != "" && i.ChannelB != "" && (i.ChainBAddress != "" || i.ChainAAddress != "")
}

type IBCRelayerNew struct {
	RelayerId            string            `bson:"relayer_id"`
	RelayerName          string            `bson:"relayer_name"`
	RelayerIcon          string            `bson:"relayer_icon"`
	ServedChains         int64             `bson:"served_chains"`
	ChannelPairInfo      []ChannelPairInfo `bson:"channel_pair_info"`
	UpdateTime           int64             `bson:"update_time"`
	RelayedTotalTxs      int64             `bson:"relayed_total_txs"`
	RelayedSuccessTxs    int64             `bson:"relayed_success_txs"`
	RelayedTotalTxsValue string            `bson:"relayed_total_txs_value"`
	TotalFeeValue        string            `bson:"total_fee_value"`
	CreateAt             int64             `bson:"create_at"`
	UpdateAt             int64             `bson:"update_at"`
}

type ChannelPairInfo struct {
	PairId        string `bson:"pair_id"`
	ChainA        string `bson:"chain_a"`
	ChainB        string `bson:"chain_b"`
	ChannelA      string `bson:"channel_a"`
	ChannelB      string `bson:"channel_b"`
	ChainAAddress string `bson:"chain_a_address"`
	ChainBAddress string `bson:"chain_b_address"`
}

func (i IBCRelayerNew) CollectionName() string {
	return "ibc_relayer"
}

func (i ChannelPairInfo) Valid() bool {
	return i.ChainA != "" && i.ChainB != "" && i.ChannelA != "" && i.ChannelB != "" && (i.ChainBAddress != "" || i.ChainAAddress != "")
}

func GenerateChannelPairInfo(chain1, channel1, chain1Address, chain2, channel2, chain2Address string) ChannelPairInfo {
	chainA, _ := ConfirmRelayerPair(chain1, chain2)
	var res ChannelPairInfo
	if chainA == chain1 {
		res = ChannelPairInfo{
			ChainA:        chain1,
			ChainB:        chain2,
			ChannelA:      channel1,
			ChannelB:      channel2,
			ChainAAddress: chain1Address,
			ChainBAddress: chain2Address,
		}
	} else {
		res = ChannelPairInfo{
			ChainA:        chain2,
			ChainB:        chain1,
			ChannelA:      channel2,
			ChannelB:      channel1,
			ChainAAddress: chain2Address,
			ChainBAddress: chain1Address,
		}
	}

	pairId := GenerateRelayerPairId(chain1, channel1, chain1Address, chain2, channel2, chain2Address)
	res.PairId = pairId
	return res
}

func GenerateRelayerPairId(chain1, channel1, chain1Address, chain2, channel2, chain2Address string) string {
	chainA, _ := ConfirmRelayerPair(chain1, chain2)
	var pairStr string
	if chainA == chain1 {
		pairStr = fmt.Sprintf("%s%s%s%s%s%s", chain1, channel1, chain1Address, chain2, channel2, chain2Address)
	} else {
		pairStr = fmt.Sprintf("%s%s%s%s%s%s", chain2, channel2, chain2Address, chain1, channel1, chain1Address)
	}

	return utils.Md5(pairStr)
}

func ConfirmRelayerPair(chainA, chainB string) (string, string) {
	if strings.HasPrefix(strings.ToLower(chainA), constant.Cosmos) {
		return chainA, chainB
	}

	if strings.HasPrefix(strings.ToLower(chainB), constant.Cosmos) {
		return chainB, chainA
	}

	if strings.HasPrefix(strings.ToLower(chainA), constant.Iris) {
		return chainA, chainB
	}

	if strings.HasPrefix(strings.ToLower(chainB), constant.Iris) {
		return chainB, chainA
	}

	compare := strings.Compare(strings.ToLower(chainA), strings.ToLower(chainB))
	if compare < 0 {
		return chainA, chainB
	} else {
		return chainB, chainA
	}
}

func GenerateDistRelayerId(chain1, chain1Address, chain2, chain2Address string) string {
	chainA, _ := ConfirmRelayerPair(chain1, chain2)
	var pairStr string
	if chainA == chain1 {
		pairStr = fmt.Sprintf("%s:%s:%s:%s", chain1, chain1Address, chain2, chain2Address)
	} else {
		pairStr = fmt.Sprintf("%s:%s:%s:%s", chain2, chain2Address, chain1, chain1Address)
	}

	return pairStr
}

func ParseDistRelayerId(DistRelayerId string) (chain1, chain1Address, chain2, chain2Address string) {
	split := strings.Split(DistRelayerId, ":")
	chain1 = split[0]
	chain1Address = split[1]
	chain2 = split[2]
	chain2Address = split[3]
	return
}
