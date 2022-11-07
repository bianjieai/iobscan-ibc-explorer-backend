package vo

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
)

type RelayerListReq struct {
	Page
	RelayerName    string `json:"relayer_name" form:"relayer_name"`
	RelayerAddress string `json:"relayer_address" form:"relayer_address"`
	UseCount       bool   `json:"use_count" form:"use_count"`
}

type (
	RelayerDto struct {
		RelayerId            string            `json:"relayer_id"`
		RelayerName          string            `json:"relayer_name"`
		RelayerIcon          string            `json:"relayer_icon"`
		ServedChains         int64             `json:"served_chains"`
		ServedChainsInfo     []ServedChainInfo `json:"served_chains_info"`
		UpdateTime           int64             `json:"update_time"`
		RelayedTotalTxs      int64             `json:"relayed_total_txs"`
		RelayedSuccessTxs    int64             `json:"relayed_success_txs"`
		RelayedTotalTxsValue string            `json:"relayed_total_txs_value"`
		TotalFeeValue        string            `json:"total_fee_value"`
	}
	ServedChainInfo struct {
		Chain     string   `json:"chain"`
		Addresses []string `json:"addresses"`
	}
)

type RelayerListResp struct {
	Items     []RelayerDto `json:"items"`
	PageInfo  PageInfo     `json:"page_info"`
	TimeStamp int64        `json:"time_stamp"`
}

func (dto RelayerDto) LoadDto(relayer *entity.IBCRelayerNew) RelayerDto {

	getServedChainInfo := func() []ServedChainInfo {

		setMap := make(map[string]ServedChainInfo, len(relayer.ChannelPairInfo))

		for _, val := range relayer.ChannelPairInfo {

			if cacheValue, ok := setMap[val.ChainA]; ok {
				if val.ChainAAddress != "" {
					cacheValue.Addresses = append(cacheValue.Addresses, val.ChainAAddress)
				}
				cacheValue.Addresses = utils.DistinctSliceStr(cacheValue.Addresses)
				setMap[val.ChainA] = cacheValue

			} else {
				item := ServedChainInfo{
					Chain: val.ChainA,
				}
				if val.ChainAAddress != "" {
					item.Addresses = []string{val.ChainAAddress}
				}
				setMap[val.ChainA] = item
			}

			if cacheValue, ok := setMap[val.ChainB]; ok {
				if val.ChainBAddress != "" {
					cacheValue.Addresses = append(cacheValue.Addresses, val.ChainBAddress)
				}
				cacheValue.Addresses = utils.DistinctSliceStr(cacheValue.Addresses)
				setMap[val.ChainB] = cacheValue
			} else {
				item := ServedChainInfo{
					Chain: val.ChainB,
				}
				if val.ChainBAddress != "" {
					item.Addresses = []string{val.ChainBAddress}
				}
				setMap[val.ChainB] = item
			}

		}
		retData := make([]ServedChainInfo, 0, len(setMap))
		for _, info := range setMap {
			retData = append(retData, info)
		}
		return retData
	}

	return RelayerDto{
		RelayerId:            relayer.RelayerId,
		RelayerName:          relayer.RelayerName,
		RelayerIcon:          relayer.RelayerIcon,
		ServedChains:         relayer.ServedChains,
		ServedChainsInfo:     getServedChainInfo(),
		UpdateTime:           relayer.UpdateTime,
		RelayedTotalTxs:      relayer.RelayedTotalTxs,
		RelayedSuccessTxs:    relayer.RelayedSuccessTxs,
		RelayedTotalTxsValue: relayer.RelayedTotalTxsValue,
		TotalFeeValue:        relayer.TotalFeeValue,
	}
}

type IobRegistryRelayerInfoResp struct {
	TeamName string `json:"team name"`
	Contact  struct {
		Website  string `json:"website"`
		Github   string `json:"github"`
		Telegram string `json:"telegram"`
		Twitter  string `json:"twitter"`
		Medium   string `json:"medium"`
		Discord  string `json:"discord"`
	} `json:"contact"`
	Introduction []string `json:"introduction"`
}

type IobRegistryRelayerPairResp struct {
	Chain1 struct {
		Address   string `json:"address"`
		ChainId   string `json:"chain-id"`
		ChannelId string `json:"channel-id"`
		Version   string `json:"version"`
	} `json:"chain-1"`
	Chain2 struct {
		Address   string `json:"address"`
		ChainId   string `json:"chain-id"`
		ChannelId string `json:"channel-id"`
		Version   string `json:"version"`
	} `json:"chain-2"`
}

type (
	ChannelPairInfoDto struct {
		ChainA            string   `json:"chain_a"`
		ChainB            string   `json:"chain_b"`
		ChannelA          string   `json:"channel_a"`
		ChannelB          string   `json:"channel_b"`
		ChainAAddresses   []string `json:"chain_a_addresses"`
		ChainBAddresses   []string `json:"chain_b_addresses"`
		ChannelPairStatus int      `json:"channel_pair_status"`
	}
	RelayerDetailResp struct {
		RelayerId            string               `json:"relayer_id"`
		RelayerName          string               `json:"relayer_name"`
		RelayerIcon          string               `json:"relayer_icon"`
		ServedChains         int64                `json:"served_chains"`
		ChannelPairInfo      []ChannelPairInfoDto `json:"channel_pair_info"`
		UpdateTime           int64                `json:"update_time"`
		RelayedTotalTxs      int64                `json:"relayed_total_txs"`
		RelayedSuccessTxs    int64                `json:"relayed_success_txs"`
		RelayedTotalTxsValue string               `json:"relayed_total_txs_value"`
		TotalFeeValue        string               `json:"total_fee_value"`

		TimeStamp int64 `json:"time_stamp"`
	}
)

func LoadRelayerDetailDto(relayer *entity.IBCRelayerNew, statusMap map[string]int) RelayerDetailResp {

	getChannelPairInfo := func() []ChannelPairInfoDto {

		setMap := make(map[string]ChannelPairInfoDto, len(relayer.ChannelPairInfo))

		for _, val := range relayer.ChannelPairInfo {
			key := val.ChainA + val.ChannelA + val.ChainB + val.ChannelB
			if cacheValue, ok := setMap[key]; ok {
				if val.ChainAAddress != "" {
					cacheValue.ChainAAddresses = append(cacheValue.ChainAAddresses, val.ChainAAddress)
				}
				if val.ChainBAddress != "" {
					cacheValue.ChainBAddresses = append(cacheValue.ChainBAddresses, val.ChainBAddress)
				}

				cacheValue.ChainAAddresses = utils.DistinctSliceStr(cacheValue.ChainAAddresses)
				cacheValue.ChainBAddresses = utils.DistinctSliceStr(cacheValue.ChainBAddresses)

				setMap[key] = cacheValue

			} else {
				item := ChannelPairInfoDto{
					ChainA:   val.ChainA,
					ChainB:   val.ChainB,
					ChannelA: val.ChannelA,
					ChannelB: val.ChannelB,
				}

				if val.ChainAAddress != "" {
					item.ChainAAddresses = []string{val.ChainAAddress}
				}
				if val.ChainBAddress != "" {
					item.ChainBAddresses = []string{val.ChainBAddress}
				}

				setMap[key] = item
			}

		}
		retData := make([]ChannelPairInfoDto, 0, len(setMap))
		for key, info := range setMap {
			info.ChannelPairStatus = statusMap[key]
			retData = append(retData, info)
		}
		return retData
	}

	return RelayerDetailResp{
		RelayerId:            relayer.RelayerId,
		RelayerName:          relayer.RelayerName,
		RelayerIcon:          relayer.RelayerIcon,
		ServedChains:         relayer.ServedChains,
		ChannelPairInfo:      getChannelPairInfo(),
		UpdateTime:           relayer.UpdateTime,
		RelayedTotalTxs:      relayer.RelayedTotalTxs,
		RelayedSuccessTxs:    relayer.RelayedSuccessTxs,
		RelayedTotalTxsValue: relayer.RelayedTotalTxsValue,
		TotalFeeValue:        relayer.TotalFeeValue,
	}
}
