package service

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"strings"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
)

type IRelayerService interface {
	List(req *vo.RelayerListReq) (vo.RelayerListResp, errors.Error)
	ListCount(req *vo.RelayerListReq) (int64, errors.Error)
	Collect(OperatorFile string) errors.Error
	Detail(relayerId string) (vo.RelayerDetailResp, errors.Error)
	DetailRelayerTxsCount(req *vo.DetailRelayerTxsReq) (int64, errors.Error)
	DetailRelayerTxs(req *vo.DetailRelayerTxsReq) (vo.DetailRelayerTxsResp, errors.Error)
	CheckRelayerParams(relayerId, chain string) (vo.ServedChainInfo, errors.Error)
}

type RelayerService struct {
	dto            vo.RelayerDto
	relayerHandler RelayerHandler
}

var _ IRelayerService = new(RelayerService)

func (svc *RelayerService) List(req *vo.RelayerListReq) (vo.RelayerListResp, errors.Error) {
	var resp vo.RelayerListResp
	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	rets, total, err := relayerRepo.FindAllBycond(req.RelayerName, req.RelayerAddress, skip, limit, req.UseCount)
	if err != nil {
		return resp, errors.Wrap(err)
	}
	for _, val := range rets {
		item := svc.dto.LoadDto(val)
		resp.Items = append(resp.Items, item)
	}
	page := vo.BuildPageInfo(total, req.PageNum, req.PageSize)
	resp.PageInfo = page
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}

func (svc *RelayerService) ListCount(req *vo.RelayerListReq) (int64, errors.Error) {
	total, err := relayerRepo.CountBycond(req.RelayerName, req.RelayerAddress)
	if err != nil {
		return 0, errors.Wrap(err)
	}

	return total, nil
}

func (svc *RelayerService) Collect(OperatorFile string) errors.Error {
	go svc.relayerHandler.Collect(OperatorFile)
	return nil
}

func (svc *RelayerService) Detail(relayerId string) (vo.RelayerDetailResp, errors.Error) {
	var resp vo.RelayerDetailResp
	one, err := relayerRepo.FindOneByRelayerId(relayerId)
	if err != nil {
		return resp, errors.Wrap(err)
	}

	channelPairs, err := channelRepo.FindAll()
	if err != nil {
		return resp, errors.Wrap(err)
	}
	channelPairStatusMap := make(map[string]int, len(channelPairs))
	for _, val := range channelPairs {
		channelPairStatusMap[val.ChainA+val.ChannelA+val.ChainB+val.ChannelB] = int(val.Status)
	}

	resp = vo.LoadRelayerDetailDto(one, channelPairStatusMap)

	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}

func (svc *RelayerService) DetailRelayerTxs(req *vo.DetailRelayerTxsReq) (vo.DetailRelayerTxsResp, errors.Error) {
	var resp vo.DetailRelayerTxsResp
	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	txs, err := txRepo.GetRelayerTxs(req.Chain, req.Addresses, constant.RelayerDetailTxsType, req.TxTimeStart, req.TxTimeEnd, skip, limit)
	if err != nil {
		return resp, errors.Wrap(err)
	}
	denomMap := make(map[string]*entity.IBCDenom, len(txs))
	ibcDenoms, err := denomRepo.FindByChainId(req.Chain)
	if err != nil {
		return resp, errors.Wrap(err)
	}
	for _, val := range ibcDenoms {
		denomMap[val.Denom] = val
	}

	items := make([]vo.RelayerTxsDto, 0, len(txs))
	for _, tx := range txs {
		item := vo.LoadRelayerTxsDto(tx, req.Chain)
		item.DenomInfo = getTxDenomInfo(tx, req.Chain, denomMap)
		items = append(items, item)
	}
	resp.Items = items
	page := vo.BuildPageInfo(int64(len(items)), req.PageNum, req.PageSize)
	resp.PageInfo = page
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}

func (svc *RelayerService) DetailRelayerTxsCount(req *vo.DetailRelayerTxsReq) (int64, errors.Error) {
	count, err := txRepo.CountRelayerTxs(req.Chain, req.Addresses, constant.RelayerDetailTxsType, req.TxTimeStart, req.TxTimeEnd)
	if err != nil {
		return 0, errors.Wrap(err)
	}
	return count, nil
}

func getRelayerChainsInfo(relayerId string) (map[string]vo.ServedChainInfo, error) {
	one, err := relayerRepo.FindOneByRelayerId(relayerId)
	if err != nil {
		return nil, err
	}
	servedChainsInfoMap := vo.GetChainInfoFromChannelPair(one.ChannelPairInfo)
	return servedChainsInfoMap, nil
}

func getMsgAmtDenom(msg *model.TxMsg) string {
	var denom string
	switch msg.Type {
	case constant.MsgTypeAcknowledgement:
		packetData := msg.AckPacketMsg().Packet.Data
		if strings.Contains(packetData.Denom, "/") {
			denom = utils.IbcHash(packetData.Denom)
		} else {
			denom = packetData.Denom
		}
	case constant.MsgTypeTimeoutPacket:
		packetData := msg.TimeoutPacketMsg().Packet.Data
		if strings.Contains(packetData.Denom, "/") {
			denom = utils.IbcHash(packetData.Denom)
		} else {
			denom = packetData.Denom
		}
	case constant.MsgTypeRecvPacket:
		packet := msg.RecvPacketMsg().Packet
		dcPrefix := fmt.Sprintf("%s/%s", packet.DestinationPort, packet.DestinationChannel)
		if strings.HasPrefix(packet.Data.Denom, dcPrefix) {
			arrs := strings.Split(packet.Data.Denom, "/")
			if len(arrs) == 3 {
				denom = arrs[2]
			} else {
				denomFullPath := strings.Join(arrs[2:], "/")
				denom = utils.IbcHash(denomFullPath)
			}
		} else {
			denomFullPath := strings.Join([]string{packet.DestinationPort, packet.DestinationChannel, packet.Data.Denom}, "/")
			denom = utils.IbcHash(denomFullPath)
		}
	}
	return denom
}

func getTxDenomInfo(tx *entity.Tx, chain string, denomMap map[string]*entity.IBCDenom) vo.DenomInfo {
	denomInfosMap := make(map[string]vo.DenomInfo, len(tx.DocTxMsgs))
	for _, msg := range tx.DocTxMsgs {
		if utils.InArray(constant.RelayerDetailTxsType, msg.Type) {
			var denomInfo vo.DenomInfo
			packetData := msg.AckPacketMsg().Packet.Data
			denomInfo.Amount = fmt.Sprint(packetData.Amount)
			denomInfo.Denom = getMsgAmtDenom(msg)
			denomInfo.DenomChain = chain
			if data, ok := denomMap[denomInfo.Denom]; ok {
				denomInfo.BaseDenomChain = data.BaseDenomChainId
				denomInfo.BaseDenom = data.BaseDenom
			}

			//处理多msg数据
			if value, ok := denomInfosMap[denomInfo.Denom]; ok {
				amount, err := utils.AddByDecimal(value.Amount, denomInfo.Amount)
				if err != nil {
					continue
				} else {
					value.Amount = amount
					denomInfosMap[denomInfo.Denom] = value
				}
			} else {
				denomInfosMap[denomInfo.Denom] = denomInfo
			}
		}
	}

	for _, val := range denomInfosMap {
		//todo only support one denom
		return val
	}
	return vo.DenomInfo{}
}

func (svc *RelayerService) CheckRelayerParams(relayerId, chain string) (vo.ServedChainInfo, errors.Error) {
	servedChainsInfoMap, err := getRelayerChainsInfo(relayerId)
	if err != nil {
		return vo.ServedChainInfo{}, errors.Wrap(err)
	}
	chainInfo, ok := servedChainsInfoMap[chain]
	if !ok {
		return vo.ServedChainInfo{}, errors.WrapBadRequest(fmt.Errorf("%s is not served chain of relayer_id %s",
			chain, relayerId))
	}
	return chainInfo, nil
}
