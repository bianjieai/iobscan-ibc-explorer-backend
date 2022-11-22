package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/qiniu/qmgo"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type ITransferService interface {
	TransferTxsCount(req *vo.TranaferTxsReq) (*vo.TransferTxsCountResp, errors.Error)
	TransferTxs(req *vo.TranaferTxsReq) (vo.TranaferTxsResp, errors.Error)
	TransferTxDetailNew(hash string) (*vo.TranaferTxDetailNewResp, errors.Error)
	TraceSource(hash string, req *vo.TraceSourceReq) (vo.TraceSourceResp, errors.Error)
	SearchCondition() (*vo.SearchConditionResp, errors.Error)
}

var _ ITransferService = new(TransferService)

type TransferService struct {
	dto vo.IbcTxDto
}

func createIbcTxQuery(req *vo.TranaferTxsReq) (dto.IbcTxQuery, error) {
	var (
		query dto.IbcTxQuery
		err   error
	)
	if req.Chain != "" {
		query.Chain = strings.Split(req.Chain, ",")
	}
	if req.DateRange != "" {
		dateRange := strings.Split(req.DateRange, ",")
		if len(dateRange) == 2 {
			query.StartTime, err = strconv.ParseInt(dateRange[0], 10, 64)
			if err != nil {
				return query, err
			}
			query.EndTime, err = strconv.ParseInt(dateRange[1], 10, 64)
			if err != nil {
				return query, err
			}
		}
	}
	if req.Status != "" {
		stats := strings.Split(req.Status, ",")
		for _, val := range stats {
			stat, err := strconv.Atoi(val)
			if err != nil {
				return query, err
			}
			query.Status = append(query.Status, stat)
		}
	}

	if req.BaseDenom != "" {
		if strings.ToLower(req.BaseDenom) == constant.OtherDenom {
			tokens, err := getUnAuthToken()
			if err != nil {
				return query, err
			}
			query.BaseDenom = tokens
		} else {
			query.BaseDenom = []string{req.BaseDenom}
			query.BaseDenomChain = req.BaseDenomChain
		}
	} else if req.Denom != "" {
		query.Denom = req.Denom
	}
	return query, nil
}
func (t TransferService) TransferTxsCount(req *vo.TranaferTxsReq) (*vo.TransferTxsCountResp, errors.Error) {
	query, err := createIbcTxQuery(req)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if len(query.Chain) > 2 {
		return nil, errors.WrapBadRequest(fmt.Errorf("invalid chain id"))
	}

	txsCountChan := make(chan *vo.TxsCountChanDTO)
	txsValueChan := make(chan *vo.TxsValueChanDTO)
	// 计算交易总数
	go func() {
		checkDisplayMax := func(count int64) int64 {
			if count > constant.DisplayIbcRecordMax {
				return constant.DisplayIbcRecordMax
			}
			return count
		}
		//default cond
		if len(query.Chain) == 0 && len(query.Status) == 4 && query.StartTime == 0 && len(query.BaseDenom) == 0 && query.Denom == "" {
			data, err2 := statisticRepo.FindOne(constant.TxLatestAllStatisticName)
			if err2 != nil {
				txsCountChan <- &vo.TxsCountChanDTO{Count: 0, Err: err2}
				return
			}
			txsCountChan <- &vo.TxsCountChanDTO{Count: checkDisplayMax(data.Count), Err: nil}
			return
		}
		count, err2 := ibcTxRepo.CountTransferTxs(query)
		if err2 != nil {
			txsCountChan <- &vo.TxsCountChanDTO{Count: 0, Err: err2}
			return
		}
		txsCountChan <- &vo.TxsCountChanDTO{Count: checkDisplayMax(count), Err: nil}
	}()

	// 计算交易价值
	go func() {
		if (len(query.Chain) == 0) ||
			(len(query.Chain) == 1 && query.Chain[0] == constant.AllChain) ||
			(query.Chain[0] == constant.AllChain && query.Chain[1] == constant.AllChain) {
			txsValueChan <- &vo.TxsValueChanDTO{Value: "", Err: nil}
			return
		}

		aggrTxsValue, err2 := ibcTxRepo.AggrTxsValue(query, false)
		if err2 != nil {
			txsValueChan <- &vo.TxsValueChanDTO{Value: "", Err: err2}
			return
		}

		priceMap := cache.TokenPriceMap()
		totalTxsValue := decimal.Zero
		for _, v := range aggrTxsValue {
			txsValue := CalculateDenomValue(priceMap, v.BaseDenom, v.BaseDenomChain, decimal.NewFromFloat(v.Amount))
			totalTxsValue = totalTxsValue.Add(txsValue)
		}

		txsValueChan <- &vo.TxsValueChanDTO{Value: totalTxsValue.String(), Err: err2}
		return
	}()

	txsCountRes := <-txsCountChan
	txsValueRes := <-txsValueChan
	if txsCountRes.Err != nil {
		return nil, errors.Wrap(txsCountRes.Err)
	} else if txsValueRes.Err != nil {
		return nil, errors.Wrap(txsValueRes.Err)
	} else {
		return &vo.TransferTxsCountResp{
			TxsCount: txsCountRes.Count,
			TxsValue: txsValueRes.Value,
		}, nil
	}
}

func (t TransferService) TransferTxs(req *vo.TranaferTxsReq) (vo.TranaferTxsResp, errors.Error) {
	var resp vo.TranaferTxsResp
	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	query, err := createIbcTxQuery(req)
	if err != nil {
		return resp, errors.Wrap(err)
	}
	if len(query.Chain) > 2 {
		return resp, nil
	}
	res, err := ibcTxRepo.FindTransferTxs(query, skip, limit)
	if err != nil {
		return resp, errors.Wrap(err)
	}
	items := make([]vo.IbcTxDto, 0, len(res))
	for _, val := range res {
		item := t.dto.LoadDto(val)
		items = append(items, item)
	}
	resp.Items = items
	page := vo.BuildPageInfo(int64(len(items)), req.PageNum, req.PageSize)
	resp.PageInfo = page
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}

func (t TransferService) TransferTxDetailNew(hash string) (*vo.TranaferTxDetailNewResp, errors.Error) {
	var resp vo.TranaferTxDetailNewResp
	ibcTxs, err := ibcTxRepo.TxDetail(hash, false)
	if err != nil && err != qmgo.ErrNoSuchDocuments {
		return nil, errors.Wrap(err)
	}
	if len(ibcTxs) == 0 {
		ibcTxs, err = ibcTxRepo.TxDetail(hash, true)
		if err != nil && err != qmgo.ErrNoSuchDocuments {
			return nil, errors.Wrap(err)
		}
	}
	if len(ibcTxs) == 0 {
		return nil, nil
	}

	if len(ibcTxs) == 1 {
		resp = vo.LoadTranaferTxDetail(ibcTxs[0])
		resp.RelayerInfo, err = getRelayerInfo(ibcTxs[0])
		if err != nil {
			return nil, errors.Wrap(err)
		}
		resp.TokenInfo, err = getTokenInfo(ibcTxs[0])
		if err != nil {
			return nil, errors.Wrap(err)
		}
	} else if len(ibcTxs) > 1 {
		resp.IsList = true
		for _, val := range ibcTxs {
			item := t.dto.LoadDto(val)
			resp.Items = append(resp.Items, item)
		}
	}
	resp.TimeStamp = time.Now().Unix()
	return &resp, nil
}

func getRelayerInfo(val *entity.ExIbcTx) (*vo.RelayerInfo, error) {
	relayerMap, err := getRelayerMap()
	if err != nil {
		return nil, err
	}
	if val.DcTxInfo == nil && val.RefundedTxInfo == nil {
		return nil, nil
	}
	var relayerInfo vo.RelayerInfo
	if val.DcTxInfo != nil && val.DcTxInfo.Msg != nil {
		dcRelayerAddr := val.DcTxInfo.Msg.CommonMsg().Signer
		relayerInfo.DcRelayer.RelayerAddr = dcRelayerAddr
	}
	if val.RefundedTxInfo != nil && val.RefundedTxInfo.Msg != nil {
		scRelayerAddr := val.RefundedTxInfo.Msg.CommonMsg().Signer
		relayerInfo.ScRelayer.RelayerAddr = scRelayerAddr
	}
	chainA, _ := entity.ConfirmRelayerPair(val.ScChain, val.DcChain)
	matchInfo := strings.Join([]string{relayerInfo.ScRelayer.RelayerAddr, relayerInfo.DcRelayer.RelayerAddr}, ":")
	if chainA != val.ScChain {
		matchInfo = strings.Join([]string{relayerInfo.DcRelayer.RelayerAddr, relayerInfo.ScRelayer.RelayerAddr}, ":")
	}
	if value, ok := relayerMap[matchInfo]; ok {
		relayerInfo.ScRelayer.RelayerName = value.RelayerName
		relayerInfo.ScRelayer.Icon = value.RelayerIcon
	}
	return &relayerInfo, nil
}
func getRelayerMap() (map[string]entity.IBCRelayerNew, error) {
	relayers, err := relayerCache.FindAll()
	if err != nil {
		return nil, err
	}
	relayerMap := make(map[string]entity.IBCRelayerNew, len(relayers))
	for _, val := range relayers {
		data := entity.IBCRelayerNew{
			RelayerIcon: val.RelayerIcon,
			RelayerName: val.RelayerName,
		}
		for _, one := range val.ChannelPairInfo {
			key := strings.Join([]string{one.ChainAAddress, one.ChainBAddress}, ":")
			if _, exist := relayerMap[key]; exist {
				continue
			}
			relayerMap[key] = data
		}
	}
	return relayerMap, nil
}

func getTokenInfo(ibcTx *entity.ExIbcTx) (*vo.TokenInfo, error) {
	var (
		sendToken = vo.DetailToken{
			Denom: ibcTx.Denoms.ScDenom,
		}
		recvToken = vo.DetailToken{
			Denom: ibcTx.Denoms.DcDenom,
		}
	)
	if strings.HasPrefix(ibcTx.Denoms.ScDenom, "ibc/") {
		denom, err := denomRepo.FindByDenomChain(ibcTx.Denoms.ScDenom, ibcTx.ScChain)
		if err != nil && err != qmgo.ErrNoSuchDocuments {
			return nil, err
		}
		if denom != nil {
			sendToken.DenomPath = strings.Join([]string{denom.DenomPath, denom.RootDenom}, "/")
		}
	}
	if strings.HasPrefix(ibcTx.Denoms.DcDenom, "ibc/") {
		denom, err := denomRepo.FindByDenomChain(ibcTx.Denoms.DcDenom, ibcTx.DcChain)
		if err != nil && err != qmgo.ErrNoSuchDocuments {
			return nil, err
		}
		if denom != nil {
			recvToken.DenomPath = strings.Join([]string{denom.DenomPath, denom.RootDenom}, "/")
		}
	}
	return &vo.TokenInfo{
		BaseDenom:      ibcTx.BaseDenom,
		BaseDenomChain: ibcTx.BaseDenomChain,
		Amount:         ibcTx.ScTxInfo.MsgAmount.Amount,
		SendToken:      sendToken,
		RecvToken:      recvToken,
	}, nil
}

func getUnAuthToken() ([]string, error) {
	value, err := cache.GetRedisClient().Get(cache.BaseDenomUnauth)
	if err == nil && len(value) > 0 {
		var unAuthTokens []string
		utils.UnmarshalJsonIgnoreErr([]byte(value), &unAuthTokens)
		return unAuthTokens, nil
	}
	//获取ibc_base_denom表数据
	ibcBaseDenoms, err := authDenomRepo.FindAll()
	if err != nil {
		return nil, err
	}
	authDenomMap := make(map[string]struct{}, len(ibcBaseDenoms))
	for _, val := range ibcBaseDenoms {
		authDenomMap[val.Denom] = struct{}{}
	}

	//聚合ibc_denom表获取symbol为空的base_denom
	baseDenomsNoSymbol, err := denomRepo.GetBaseDenomNoSymbol()
	if err != nil {
		return nil, err
	}

	unAuthTokens := make([]string, 0, len(baseDenomsNoSymbol))
	for _, val := range baseDenomsNoSymbol {
		//移除已配置的base_denom
		if _, auth := authDenomMap[val.BaseDenom]; auth {
			continue
		}
		unAuthTokens = append(unAuthTokens, val.BaseDenom)
	}
	_ = cache.GetRedisClient().Set(cache.BaseDenomUnauth, utils.MarshalJsonIgnoreErr(unAuthTokens), cache.FiveMin)
	return unAuthTokens, nil
}

func (t TransferService) TraceSource(hash string, req *vo.TraceSourceReq) (vo.TraceSourceResp, errors.Error) {
	var resp vo.TraceSourceResp
	var supportMsgType = map[string]string{
		constant.MsgTypeRecvPacket:      "MsgRecvPacket",
		constant.MsgTypeTimeoutPacket:   "MsgTimeout",
		constant.MsgTypeAcknowledgement: "MsgAcknowledgement",
		constant.MsgTypeTransfer:        "MsgTransfer",
	}

	if len(hash) == 0 {
		return resp, errors.WrapBadRequest(fmt.Errorf("invalid hash"))
	}

	msgType, ok := supportMsgType[req.MsgType]
	if !ok {
		return resp, errors.WrapBadRequest(fmt.Errorf("only support transfer,recv_packet,acknowledge_packet,timeout_packet"))
	}

	value, err := lcdTxDataCache.Get(req.Chain, hash)
	if err == nil {
		utils.UnmarshalJsonIgnoreErr([]byte(value), &resp)
		return resp, nil
	}
	return getMsgAndTxData(msgType, req.Chain, hash)
}

func getMsgAndTxData(msgType, chain, hash string) (vo.TraceSourceResp, errors.Error) {
	var resp vo.TraceSourceResp

	lcdTxData, err := GetLcdTxData(chain, hash)
	if err != nil {
		return resp, err
	}
	logMap := make(map[int][]entity.Event)
	for _, val := range lcdTxData.TxResponse.Logs {
		logMap[val.MsgIndex] = val.Events
	}

	for i, val := range lcdTxData.TxResponse.Tx.Body.Messages {
		msgTy := []byte(val.Type)[strings.LastIndex(val.Type, ".")+1:]
		if string(msgTy) == msgType {
			resp.Msg = val
			resp.Events = logMap[i]
		}
	}
	if resp.Msg != nil {
		_ = lcdTxDataCache.Set(chain, hash, string(utils.MarshalJsonIgnoreErr(resp)))
	}
	return resp, nil
}

func GetLcdTxData(chain, hash string) (LcdTxData, errors.Error) {
	lcdAddrs, _ := lcdAddrCache.Get(chain)
	if len(lcdAddrs) > 0 {
		//全节点且支持交易查询
		if lcdAddrs[0].FullNode && lcdAddrs[0].TxIndexEnable {
			return GetTxDataFromChain(lcdAddrs[0].LcdAddr, hash)
		}
		//获取支持交易查询的lcd节点
		var validNodes []cache.TraceSourceLcd
		for _, val := range lcdAddrs {
			if val.TxIndexEnable {
				validNodes = append(validNodes, val)
			}
		}
		//并发处理
		return doHandleTxData(2, validNodes, hash)
	} else {
		cfg, err := chainCfgRepo.FindOne(chain)
		if err != nil {
			return LcdTxData{}, errors.Wrap(fmt.Errorf("invalid chain id"))
		}
		return GetTxDataFromChain(cfg.GrpcRestGateway, hash)
	}
}

func doHandleTxData(workNum int, lcdAddrs []cache.TraceSourceLcd, hash string) (LcdTxData, errors.Error) {
	resData := make([]LcdTxData, len(lcdAddrs))
	var wg sync.WaitGroup
	wg.Add(workNum)
	for i := 0; i < workNum; i++ {
		num := i
		go func(num int) {
			defer wg.Done()
			var err error
			for id, v := range lcdAddrs {
				if id%workNum != num {
					continue
				}
				resData[id], err = GetTxDataFromChain(v.LcdAddr, hash)
				if err == nil {
					break
				} else {
					logrus.Errorf("err:%s lcd:%s hash:%s", err.Error(), v.LcdAddr, hash)
				}
			}
		}(num)
	}
	wg.Wait()
	for i := range resData {
		if len(resData[i].TxResponse.Tx.Body.Messages) > 0 {
			return resData[i], nil
		}
	}
	return LcdTxData{}, errors.WrapLcdNodeErr("no found")
}

func GetTxDataFromChain(lcdUri string, hash string) (LcdTxData, errors.Error) {
	var txData LcdTxData
	url := fmt.Sprintf("%s/cosmos/tx/v1beta1/txs/%s", lcdUri, hash)
	resp, err := http.Get(url)
	if err != nil {
		return txData, errors.Wrap(err)
	}

	defer resp.Body.Close()
	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return txData, errors.Wrap(err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp LcdErrRespond
		if err := json.Unmarshal(bz, &errResp); err != nil {
			return txData, errors.Wrap(err)
		}
		return txData, errors.WrapLcdNodeErr(errResp.Message)
	} else {
		if err := json.Unmarshal(bz, &txData); err != nil {
			return LcdTxData{}, errors.Wrap(err)
		}
	}

	if err := json.Unmarshal(bz, &txData); err != nil {
		return LcdTxData{}, errors.Wrap(err)
	}
	return txData, nil
}

func (t TransferService) SearchCondition() (*vo.SearchConditionResp, errors.Error) {
	txTime, err := ibcTxRepo.GetMinTxTime(false)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &vo.SearchConditionResp{TxTimeMin: txTime}, nil
}
