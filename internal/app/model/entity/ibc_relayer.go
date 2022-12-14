package entity

import (
	"fmt"
	"strings"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
)

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
	return i.ChainA != "" && i.ChainB != "" && i.ChannelA != "" && i.ChannelB != "" && i.ChainBAddress != "" && i.ChainAAddress != ""
}

type ChannelPairInfoList []ChannelPairInfo

func (c ChannelPairInfoList) GetChainAddrCombs() []string {
	combMap := make(map[string]struct{})
	for _, v := range c {
		if v.ChainA != "" && v.ChainAAddress != "" {
			combMap[GenerateChainAddressComb(v.ChainA, v.ChainAAddress)] = struct{}{}
		}

		if v.ChainB != "" && v.ChainBAddress != "" {
			combMap[GenerateChainAddressComb(v.ChainB, v.ChainBAddress)] = struct{}{}
		}
	}

	combs := make([]string, 0, len(combMap))
	for k, _ := range combMap {
		combs = append(combs, k)
	}
	return combs
}

func (c ChannelPairInfoList) GetChains() []string {
	chainMap := make(map[string]struct{})
	for _, v := range c {
		if v.ChainA != "" {
			chainMap[v.ChainA] = struct{}{}
		}

		if v.ChainB != "" {
			chainMap[v.ChainB] = struct{}{}
		}
	}

	chains := make([]string, 0, len(chainMap))
	for k, _ := range chainMap {
		chains = append(chains, k)
	}
	return chains
}

func GenerateSingleSideChannelPairInfo(chain1, channel1, chain1Address string) ChannelPairInfo {
	res := ChannelPairInfo{
		ChainA:        chain1,
		ChainB:        "",
		ChannelA:      channel1,
		ChannelB:      "",
		ChainAAddress: chain1Address,
		ChainBAddress: "",
	}

	pairId := GenerateRelayerPairId(chain1, channel1, chain1Address, "", "", "")
	res.PairId = pairId
	return res
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
	if chainA == "" {
		return chainB, chainA
	}

	if chainB == "" {
		return chainA, chainB
	}

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
