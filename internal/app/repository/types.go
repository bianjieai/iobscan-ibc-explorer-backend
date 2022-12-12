package repository

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	qmgooptions "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ordered            = false
	insertIgnoreErrOpt = qmgooptions.InsertManyOptions{
		InsertHook: nil,
		InsertManyOptions: &options.InsertManyOptions{
			BypassDocumentValidation: nil,
			Ordered:                  &ordered,
		},
	}
)

// getChannelPairInfoByAddressPair 获取一对地址上的所有channel pair
func GetChannelPairInfoByAddressPair(chainA, addressA, chainB, addressB string) ([]entity.ChannelPairInfo, bool, error) {
	if chainA == "" || chainB == "" || addressA == "" || addressB == "" {
		pair := entity.GenerateChannelPairInfo(chainA, "", addressA, chainB, "", addressB)
		return []entity.ChannelPairInfo{pair}, false, nil
	}

	addrChannels, err := new(RelayerAddressChannelRepo).FindChannels([]string{addressA, addressB})
	if err != nil {
		return nil, false, err
	}

	chainAChannelMap := make(map[string]string)
	chainBChannelMap := make(map[string]string)
	for _, c := range addrChannels {
		if c.RelayerAddress == addressA && c.Chain == chainA {
			chainAChannelMap[c.Channel] = c.CounterPartyChannel
		} else if c.RelayerAddress == addressB && c.Chain == chainB {
			chainBChannelMap[c.Channel] = c.CounterPartyChannel
		}
	}

	var res []entity.ChannelPairInfo
	var channelMatched bool
	for ch, cpch := range chainAChannelMap {
		if ch2, _ := chainBChannelMap[cpch]; ch == ch2 { // channel match success
			pairInfo := entity.GenerateChannelPairInfo(chainA, ch, addressA, chainB, cpch, addressB)
			channelMatched = true
			res = append(res, pairInfo)
		}
	}

	if !channelMatched {
		pairInfo := entity.GenerateChannelPairInfo(chainA, "", addressA, chainB, "", addressB)
		res = append(res, pairInfo)
	}

	return res, channelMatched, nil
}

// getChainIdNameMap, map key: chain id, value: chain name
func GetChainIdNameMap() (map[string]string, error) {
	allChainList, err := new(ChainConfigRepo).FindAllChainInfos()
	if err != nil {
		return nil, err
	}

	allChainMap := make(map[string]string)
	for _, v := range allChainList {
		allChainMap[v.CurrentChainId] = v.ChainName
	}

	return allChainMap, err
}
