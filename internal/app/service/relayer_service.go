package service

import (
	"encoding/json"
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
)

type IRelayerService interface {
	List(req *vo.RelayerListReq) (vo.RelayerListResp, errors.Error)
	ListCount(req *vo.RelayerListReq) (int64, errors.Error)
	Collect(OperatorFile string) errors.Error
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
		if value, err := relayerTxsCache.Get(relayerId, req.Chain); err == nil {
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
		_ = relayerTxsCache.Set(relayerId, req.Chain, count)
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
	if value, err := relayerTrendCache.Get(relayerId, strconv.Itoa(req.Days)); err == nil {
		var data vo.RelayerTrendResp
		if err1 := json.Unmarshal([]byte(value), &data); err1 != nil {
			return data, errors.Wrap(err1)
		}
		return data, nil
	}
	segments := svc.getSegmentOfDay(req.Days)
	retData := svc.doHandleDaySegments(relayerAddrs, segments)
	_ = relayerTrendCache.Set(relayerId, strconv.Itoa(req.Days), utils.MustMarshalJsonToStr(retData))

	return retData, nil
}

func (svc *RelayerService) doHandleDaySegments(relayerAddrs []string, segments []*vo.DaySegment) vo.RelayerTrendResp {
	retData := make(vo.RelayerTrendResp, len(segments))
	denomPriceMap := cache.TokenPriceMap()
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		num := i
		go func(num int) {
			defer wg.Done()

			for id, v := range segments {
				if id%3 != num {
					continue
				}
				retData[id] = svc.getDayofRelayerTxsAmt(relayerAddrs, denomPriceMap, v)
			}
		}(num)
	}
	wg.Wait()
	return retData
}

func (svc *RelayerService) getSegmentOfDay(days int) []*vo.DaySegment {
	end := time.Now()
	endDay := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.Local)
	endUnix := endDay.Unix()
	startDay := endDay.Add(time.Duration(-days) * 24 * time.Hour)
	var segments []*vo.DaySegment
	for temp := startDay.Unix(); temp < endUnix; temp += 24 * 3600 {
		segments = append(segments, &vo.DaySegment{
			Date:      time.Unix(temp, 0).Format(constant.DateFormat),
			StartTime: temp,
			EndTime:   temp + 24*3600 - 1,
		})
	}
	return segments
}

func (svc *RelayerService) getDayofRelayerTxsAmt(relayerAddrs []string, denomPriceMap map[string]dto.CoinItem, segment *vo.DaySegment) vo.RelayerTrendDto {
	res, err := relayerDenomStatisticsRepo.AggrRelayerAmtAndTxsBySegment(relayerAddrs, segment.StartTime, segment.EndTime)
	if err != nil {
		logrus.Errorf("aggr date[%s] relayer amount and txs by segment[%d:%d] fail,%s", segment.Date, segment.StartTime, segment.EndTime, err.Error())
		return vo.RelayerTrendDto{}
	}
	txsNum := int64(0)
	relayerTxsAmtMap := make(map[string]dto.TxsAmtItem, 20)
	for _, item := range res {
		txsNum += item.TotalTxs
		key := fmt.Sprintf("%s%s", item.BaseDenom, item.BaseDenomChainId)
		value, exist := relayerTxsAmtMap[key]
		if exist {
			value.Amt = value.Amt.Add(decimal.NewFromFloat(item.Amount))
			relayerTxsAmtMap[key] = value
		} else {
			data := dto.TxsAmtItem{
				ChainId: item.BaseDenomChainId,
				Denom:   item.BaseDenom,
				Amt:     decimal.NewFromFloat(item.Amount),
			}
			relayerTxsAmtMap[key] = data
		}
	}

	txsValue := dto.CaculateRelayerTotalValue(denomPriceMap, relayerTxsAmtMap)

	return vo.RelayerTrendDto{
		Date:     segment.Date,
		Txs:      txsNum,
		TxsValue: txsValue.String(),
	}
}
