package entity

import (
	"fmt"
	"strings"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
)

type IBCRelayerConfig struct {
	RelayerPairId string `bson:"relayer_pair_id"`
	ChainA        string `bson:"chain_a"`
	ChainB        string `bson:"chain_b"`
	ChannelA      string `bson:"channel_a"`
	ChannelB      string `bson:"channel_b"`
	ChainAAddress string `bson:"chain_a_address"`
	ChainBAddress string `bson:"chain_b_address"`
	RelayerName   string `bson:"relayer_name"`
	Icon          string `bson:"icon"`
}

func (i IBCRelayerConfig) CollectionName() string {
	return "ibc_relayer_config"
}

func GenerateRelayerConfigEntity(chain1, channel1, chain1Address, chain2, channel2, chain2Address string) *IBCRelayerConfig {
	chainA, _ := ConfirmRelayerPair(chain1, chain2)
	var res IBCRelayerConfig
	if chainA == chain1 {
		res = IBCRelayerConfig{
			ChainA:        chain1,
			ChainB:        chain2,
			ChannelA:      channel1,
			ChannelB:      channel2,
			ChainAAddress: chain1Address,
			ChainBAddress: chain2Address,
		}
	} else {
		res = IBCRelayerConfig{
			ChainA:        chain2,
			ChainB:        chain1,
			ChannelA:      channel2,
			ChannelB:      channel1,
			ChainAAddress: chain2Address,
			ChainBAddress: chain1Address,
		}
	}

	pairStr := fmt.Sprintf("%s%s%s%s%s%s", res.ChainA, res.ChannelA, res.ChainAAddress, res.ChainB, res.ChannelB, res.ChainBAddress)
	res.RelayerPairId = utils.Md5(pairStr)
	return &res
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
