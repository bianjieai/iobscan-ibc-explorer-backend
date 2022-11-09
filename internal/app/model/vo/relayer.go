package vo

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
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
		ServedChainsNumber   int64             `json:"served_chains_number"`
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
		chainSetMap := GetChainInfoFromChannelPair(relayer.ChannelPairInfo)
		retData := make([]ServedChainInfo, 0, len(chainSetMap))
		for _, info := range chainSetMap {
			retData = append(retData, info)
		}
		return retData
	}

	return RelayerDto{
		RelayerId:            relayer.RelayerId,
		RelayerName:          relayer.RelayerName,
		RelayerIcon:          relayer.RelayerIcon,
		ServedChainsNumber:   relayer.ServedChains,
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
	TeamLogo string `json:"team logo"`
	Contact  struct {
		Website string `json:"website"`
		Github  string `json:"github"`
		Twitter string `json:"twitter"`
		Discord string `json:"discord"`
	} `json:"contact"`
	Introduction []string            `json:"introduction"`
	Addresses    []map[string]string `json:"addresses"`
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
		ServedChains         []string             `json:"served_chains"`
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

	getChannelPairInfo := func() ([]ChannelPairInfoDto, []string) {

		setMap := make(map[string]ChannelPairInfoDto, len(relayer.ChannelPairInfo))
		servedChains := make([]string, 0, 10)

		for _, val := range relayer.ChannelPairInfo {
			servedChains = append(servedChains, val.ChainA, val.ChainB)
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
		servedChains = utils.DistinctSliceStr(servedChains)
		return retData, servedChains
	}
	channelPairInfo, servedChains := getChannelPairInfo()

	return RelayerDetailResp{
		RelayerId:            relayer.RelayerId,
		RelayerName:          relayer.RelayerName,
		RelayerIcon:          relayer.RelayerIcon,
		ServedChains:         servedChains,
		ChannelPairInfo:      channelPairInfo,
		UpdateTime:           relayer.UpdateTime,
		RelayedTotalTxs:      relayer.RelayedTotalTxs,
		RelayedSuccessTxs:    relayer.RelayedSuccessTxs,
		RelayedTotalTxsValue: relayer.RelayedTotalTxsValue,
		TotalFeeValue:        relayer.TotalFeeValue,
	}
}

type DetailRelayerTxsReq struct {
	Page
	Chain       string `json:"chain" form:"chain" binding:"required"`
	TxTimeStart int64  `json:"tx_time_start" form:"tx_time_start"`
	TxTimeEnd   int64  `json:"tx_time_end" form:"tx_time_end"`
	UseCount    bool   `json:"use_count" form:"use_count"`
}

type (
	RelayerTxsDto struct {
		TxHash    string          `json:"tx_hash"`
		TxType    string          `json:"tx_type"`
		Chain     string          `json:"chain"`
		DenomInfo DenomInfo       `json:"denom_info"`
		FeeInfo   CommonInfo      `json:"fee_info"`
		TxStatus  entity.TxStatus `json:"tx_status"`
		Signer    string          `json:"signer"`
		TxTime    int64           `json:"tx_time"`
	}
	CommonInfo struct {
		Denom      string `json:"denom"`
		Amount     string `json:"amount"`
		DenomChain string `json:"denom_chain"`
	}
	DenomInfo struct {
		CommonInfo
		BaseDenom      string `json:"base_denom"`
		BaseDenomChain string `json:"base_denom_chain"`
	}
	DetailRelayerTxsResp struct {
		Items     []RelayerTxsDto `json:"items"`
		PageInfo  PageInfo        `json:"page_info"`
		TimeStamp int64           `json:"time_stamp"`
	}
)

func LoadRelayerTxsDto(tx *entity.Tx, chain string) RelayerTxsDto {
	supportTypes := []string{constant.MsgTypeRecvPacket, constant.MsgTypeAcknowledgement, constant.MsgTypeTimeoutPacket}
	getTxType := func() string {
		if utils.InArray(supportTypes, tx.Type) {
			return tx.Type
		}
		for _, val := range tx.Types {
			if utils.InArray(supportTypes, val) {
				return val
			}
		}
		return tx.Type
	}

	getFeeInfo := func() CommonInfo {
		if tx.Fee != nil {
			return CommonInfo{
				Denom:      tx.Fee.Amount[0].Denom,
				Amount:     tx.Fee.Amount[0].Amount,
				DenomChain: chain,
			}
		}
		return CommonInfo{}
	}

	getSigner := func() string {
		if len(tx.Signers) > 0 {
			return tx.Signers[0]
		}
		return ""
	}

	return RelayerTxsDto{
		TxHash:   tx.TxHash,
		TxTime:   tx.Time,
		TxType:   getTxType(),
		Chain:    chain,
		TxStatus: tx.Status,
		FeeInfo:  getFeeInfo(),
		Signer:   getSigner(),
	}
}
func GetChainInfoFromChannelPair(channelPairInfo []entity.ChannelPairInfo) map[string]ServedChainInfo {
	chainSetMap := make(map[string]ServedChainInfo, len(channelPairInfo))
	for _, val := range channelPairInfo {

		if cacheValue, ok := chainSetMap[val.ChainA]; ok {
			if val.ChainAAddress != "" {
				cacheValue.Addresses = append(cacheValue.Addresses, val.ChainAAddress)
			}
			cacheValue.Addresses = utils.DistinctSliceStr(cacheValue.Addresses)
			chainSetMap[val.ChainA] = cacheValue

		} else {
			item := ServedChainInfo{
				Chain: val.ChainA,
			}
			if val.ChainAAddress != "" {
				item.Addresses = []string{val.ChainAAddress}
			}
			chainSetMap[val.ChainA] = item
		}

		if cacheValue, ok := chainSetMap[val.ChainB]; ok {
			if val.ChainBAddress != "" {
				cacheValue.Addresses = append(cacheValue.Addresses, val.ChainBAddress)
			}
			cacheValue.Addresses = utils.DistinctSliceStr(cacheValue.Addresses)
			chainSetMap[val.ChainB] = cacheValue
		} else {
			item := ServedChainInfo{
				Chain: val.ChainB,
			}
			if val.ChainBAddress != "" {
				item.Addresses = []string{val.ChainBAddress}
			}
			chainSetMap[val.ChainB] = item
		}

	}
	return chainSetMap
}

type RelayerTrendReq struct {
	Days int `json:"days" form:"days"`
}

type (
	RelayerTrendResp []RelayerTrendDto
	RelayerTrendDto  struct {
		Date     string `json:"date"`
		Txs      int64  `json:"txs"`
		TxsValue string `json:"txs_value"`
	}
	DaySegment struct {
		Date      string
		StartTime int64
		EndTime   int64
	}
)
