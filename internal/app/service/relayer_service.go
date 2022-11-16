package service

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type IRelayerService interface {
	List(req *vo.RelayerListReq) (vo.RelayerListResp, errors.Error)
	ListCount(req *vo.RelayerListReq) (int64, errors.Error)
	Collect(operatorFile string) errors.Error
	TransferTypeTxs(relayerId string) (*vo.TransferTypeTxsResp, errors.Error)
	TotalRelayedValue(relayerId string) (*vo.TotalRelayedValueResp, errors.Error)
	TotalFeeCost(relayerId string) (*vo.TotalFeeCostResp, errors.Error)
	Detail(relayerId string) (vo.RelayerDetailResp, errors.Error)
	DetailRelayerTxsCount(relayerId string, req *vo.DetailRelayerTxsReq) (int64, errors.Error)
	DetailRelayerTxs(relayerId string, req *vo.DetailRelayerTxsReq) (vo.DetailRelayerTxsResp, errors.Error)
	RelayerNameList() ([]string, errors.Error)
	RelayerTrend(relayerId string, req *vo.RelayerTrendReq) (vo.RelayerTrendResp, errors.Error)
}

type RelayerService struct {
	dto            vo.RelayerDto
	relayerHandler RelayerHandler
}

var _ IRelayerService = new(RelayerService)

func (svc *RelayerService) List(req *vo.RelayerListReq) (vo.RelayerListResp, errors.Error) {
	var resp vo.RelayerListResp
	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	rets, err := relayerRepo.FindAllBycond(req.RelayerName, req.RelayerAddress, skip, limit)
	if err != nil {
		return resp, errors.Wrap(err)
	}
	for _, val := range rets {
		item := svc.dto.LoadDto(val)
		resp.Items = append(resp.Items, item)
	}
	page := vo.BuildPageInfo(int64(len(resp.Items)), req.PageNum, req.PageSize)
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

func (svc *RelayerService) Collect(operatorFile string) errors.Error {
	go svc.relayerHandler.Collect(operatorFile)
	return nil
}

func (svc *RelayerService) TransferTypeTxs(relayerId string) (*vo.TransferTypeTxsResp, errors.Error) {
	if res, err := relayerDataCache.GetTransferTypeTxs(relayerId); err == nil {
		return res, nil
	}

	servedChainsInfoMap, err := getRelayerChainsInfo(relayerId)
	if err != nil {
		return nil, err
	}

	relayerAddrs := make([]string, 0, len(servedChainsInfoMap))
	for _, val := range servedChainsInfoMap {
		relayerAddrs = append(relayerAddrs, val.Addresses...)
	}

	aggrRes, e := relayerDenomStatisticsRepo.AggrAmtByTxType(relayerAddrs)
	if e != nil {
		return nil, errors.Wrap(e)
	}

	var res vo.TransferTypeTxsResp
	for _, v := range aggrRes {
		txType := entity.TxType(v.TxType)
		switch txType {
		case entity.TxTypeTimeoutPacket:
			res.TimeoutPacketTxs = v.TotalTxs
		case entity.TxTypeRecvPacket:
			res.RecvPacketTxs = v.TotalTxs
		case entity.TxTypeAckPacket:
			res.AcknowledgePacketTxs = v.TotalTxs
		}
	}

	_ = relayerDataCache.SetTransferTypeTxs(relayerId, &res)
	return &res, nil
}

func (svc *RelayerService) TotalRelayedValue(relayerId string) (*vo.TotalRelayedValueResp, errors.Error) {
	_, err := getRelayerChainsInfo(relayerId)
	if err != nil {
		return nil, err
	}

	res, err1 := relayerDataCache.GetTotalRelayedValue(relayerId)
	if err1 != nil {
		return nil, errors.Wrap(err1)
	}
	return res, nil
}

func (svc *RelayerService) TotalFeeCost(relayerId string) (*vo.TotalFeeCostResp, errors.Error) {
	_, err := getRelayerChainsInfo(relayerId)
	if err != nil {
		return nil, err
	}
	res, err1 := relayerDataCache.GetTotalFeeCost(relayerId)
	if err1 != nil {
		return nil, errors.Wrap(err1)
	}
	return res, nil
}

func (svc *RelayerService) Detail(relayerId string) (vo.RelayerDetailResp, errors.Error) {
	var resp vo.RelayerDetailResp
	one, err := relayerRepo.FindOneByRelayerId(relayerId)
	if err != nil {
		return resp, errors.Wrap(err)
	}

	//未注册的relayer的名字为空
	if one.RelayerName == "" {
		return resp, errors.WrapRelayerNoAccessDetailErr(fmt.Errorf("welcome to register this relayer !"))
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

func (svc *RelayerService) DetailRelayerTxs(relayerId string, req *vo.DetailRelayerTxsReq) (vo.DetailRelayerTxsResp, errors.Error) {
	var resp vo.DetailRelayerTxsResp
	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	chainInfo, err0 := svc.checkRelayerParams(relayerId, req.Chain)
	if err0 != nil {
		return resp, err0
	}
	txs, err := txRepo.GetRelayerTxs(req.Chain, chainInfo.Addresses, constant.RelayerDetailTxsType, req.TxTimeStart, req.TxTimeEnd, skip, limit)
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

func (svc *RelayerService) DetailRelayerTxsCount(relayerId string, req *vo.DetailRelayerTxsReq) (int64, errors.Error) {
	//默认情况下chain查询交易
	if req.TxTimeStart == 0 && req.TxTimeEnd == 0 {
		if value, err := relayerDataCache.GetTotalTxs(relayerId, req.Chain); err == nil {
			return value, nil
		}
	}
	chainInfo, err0 := svc.checkRelayerParams(relayerId, req.Chain)
	if err0 != nil {
		return 0, err0
	}
	count, err := txRepo.CountRelayerTxs(req.Chain, chainInfo.Addresses, constant.RelayerDetailTxsType, req.TxTimeStart, req.TxTimeEnd)
	if err != nil {
		return 0, errors.Wrap(err)
	}
	if req.TxTimeStart == 0 && req.TxTimeEnd == 0 {
		_ = relayerDataCache.SetTotalTxs(relayerId, req.Chain, count)
	}
	return count, nil
}

func getRelayerChainsInfo(relayerId string) (map[string]vo.ServedChainInfo, errors.Error) {
	one, err := relayerRepo.FindOneByRelayerId(relayerId)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	//未注册的relayer的名字为空
	if one.RelayerName == "" {
		return nil, errors.WrapRelayerNoAccessDetailErr(fmt.Errorf("welcome to register this relayer !"))
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

func (svc *RelayerService) checkRelayerParams(relayerId, chain string) (vo.ServedChainInfo, errors.Error) {
	servedChainsInfoMap, err := getRelayerChainsInfo(relayerId)
	if err != nil {
		return vo.ServedChainInfo{}, err
	}
	chainInfo, ok := servedChainsInfoMap[chain]
	if !ok {
		return vo.ServedChainInfo{}, errors.WrapBadRequest(fmt.Errorf("%s is not served chain of relayer_id %s",
			chain, relayerId))
	}
	return chainInfo, nil
}

func (svc *RelayerService) RelayerNameList() ([]string, errors.Error) {
	relayers, err := relayerRepo.RelayerNameList()
	if err != nil {
		return nil, errors.Wrap(err)
	}
	res := make([]string, 0, len(relayers))
	for _, val := range relayers {
		res = append(res, val.RelayerName)
	}
	return res, nil
}

func (svc *RelayerService) RelayerTrend(relayerId string, req *vo.RelayerTrendReq) (vo.RelayerTrendResp, errors.Error) {
	servedChainsInfoMap, err := getRelayerChainsInfo(relayerId)
	if err != nil {
		return vo.RelayerTrendResp{}, err
	}
	relayerAddrs := make([]string, 0, 10)
	for _, val := range servedChainsInfoMap {
		relayerAddrs = append(relayerAddrs, val.Addresses...)
	}
	if req.Days <= 0 {
		req.Days = 30
	}

	//从缓存取数据返回
	if value, err := relayerDataCache.GetRelayedTrend(relayerId, strconv.Itoa(req.Days)); err == nil {
		var data vo.RelayerTrendResp
		if err1 := json.Unmarshal([]byte(value), &data); err1 != nil {
			return data, errors.Wrap(err1)
		}
		return data, nil
	}
	segments := svc.getSegmentOfDay(req.Days)
	retData := svc.doHandleDaySegments(relayerAddrs, segments)
	_ = relayerDataCache.SetRelayedTrend(relayerId, strconv.Itoa(req.Days), utils.MustMarshalJsonToStr(retData))

	return retData, nil
}

func (svc *RelayerService) doHandleDaySegments(relayerAddrs []string, segments []*vo.DaySegment) vo.RelayerTrendResp {
	denomPriceMap := cache.TokenPriceMap()
	retData := svc.getDayofRelayerTxsAmt(relayerAddrs, denomPriceMap, segments)
	return retData
}

func (svc *RelayerService) getSegmentOfDay(days int) []*vo.DaySegment {
	end := time.Now()
	endDay := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.Local)
	endUnix := endDay.Unix()
	startDay := endDay.Add(time.Duration(-days+1) * 24 * time.Hour)
	var segments []*vo.DaySegment
	for temp := startDay.Unix(); temp <= endUnix; temp += 24 * 3600 {
		segments = append(segments, &vo.DaySegment{
			Date:      time.Unix(temp, 0).Format(constant.DateFormat),
			StartTime: temp,
			EndTime:   temp + 24*3600 - 1,
		})
	}
	return segments
}

func (svc *RelayerService) getDayofRelayerTxsAmt(relayerAddrs []string, denomPriceMap map[string]dto.CoinItem, segments []*vo.DaySegment) vo.RelayerTrendResp {
	res, err := relayerDenomStatisticsRepo.AggrRelayerAmtAndTxsBySegment(relayerAddrs, segments[0].StartTime, segments[len(segments)-1].EndTime)
	if err != nil {
		logrus.Errorf("aggr  relayer amount and txs by segment  %d-%d  fail,%s", segments[0].StartTime, segments[len(segments)-1].EndTime, err.Error())
		return vo.RelayerTrendResp{}
	}
	segmentTxsValueMap := make(map[string]dto.TxsAmtItem, 20)
	for _, item := range res {

		//计算价值
		baseDenomValue := decimal.NewFromFloat(0)
		decAmt := decimal.NewFromFloat(item.Amount)
		priceKey := item.BaseDenom + item.BaseDenomChain
		if coin, ok := denomPriceMap[priceKey]; ok {
			if coin.Scale > 0 {
				baseDenomValue = decAmt.Div(decimal.NewFromFloat(math.Pow10(coin.Scale))).Mul(decimal.NewFromFloat(coin.Price))
			}
		}

		key := time.Unix(item.SegmentStartTime, 0).Format(constant.DateFormat)
		value, exist := segmentTxsValueMap[key]
		if exist {
			value.AmtValue = value.AmtValue.Add(baseDenomValue)
			value.Txs += item.TotalTxs
			segmentTxsValueMap[key] = value
		} else {
			data := dto.TxsAmtItem{
				ChainId:  item.BaseDenomChain,
				Denom:    item.BaseDenom,
				Txs:      item.TotalTxs,
				AmtValue: baseDenomValue,
			}
			segmentTxsValueMap[key] = data
		}
	}

	retData := make(vo.RelayerTrendResp, 0, len(segments))
	for _, segment := range segments {
		data, ok := segmentTxsValueMap[segment.Date]
		if ok {
			item := vo.RelayerTrendDto{
				Date:     segment.Date,
				Txs:      data.Txs,
				TxsValue: data.AmtValue.String(),
			}
			retData = append(retData, item)
		} else {
			retData = append(retData, vo.RelayerTrendDto{
				Date: segment.Date,
			})
		}
	}
	return retData
}
