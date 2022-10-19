package service

import (
	"encoding/json"
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/qiniu/qmgo"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type ITransferService interface {
	TransferTxsCount(req *vo.TranaferTxsReq) (int64, errors.Error)
	TransferTxs(req *vo.TranaferTxsReq) (vo.TranaferTxsResp, errors.Error)
	TransferTxDetail(hash string) (vo.TranaferTxDetailResp, errors.Error)
	TransferTxDetailNew(hash string) (*vo.TranaferTxDetailNewResp, errors.Error)
	TraceSource(hash string, req *vo.TraceSourceReq) (vo.TraceSourceResp, errors.Error)
}

var _ ITransferService = new(TransferService)

type TransferService struct {
	dto       vo.IbcTxDto
	detailDto vo.IbcTxDetailDto
}

func createIbcTxQuery(req *vo.TranaferTxsReq) (dto.IbcTxQuery, error) {
	var (
		query dto.IbcTxQuery
		err   error
	)
	if req.ChainId != "" {
		query.ChainId = strings.Split(req.ChainId, ",")
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
			query.BaseDenomChainId = req.BaseDenomChainId
		}
	} else if req.Denom != "" {
		query.Denom = req.Denom
	}
	return query, nil
}
func (t TransferService) TransferTxsCount(req *vo.TranaferTxsReq) (int64, errors.Error) {
	query, err := createIbcTxQuery(req)
	if err != nil {
		return 0, errors.Wrap(err)
	}
	if len(query.ChainId) > 2 {
		return 0, nil
	}

	checkDisplayMax := func(count int64) int64 {
		if count > constant.DisplayIbcRecordMax {
			return constant.DisplayIbcRecordMax
		}
		return count
	}
	//default cond
	if len(query.ChainId) == 0 && len(query.Status) == 4 && query.StartTime == 0 && len(query.BaseDenom) == 0 && query.Denom == "" {
		data, err := statisticRepo.FindOne(constant.TxLatestAllStatisticName)
		if err != nil {
			return 0, errors.Wrap(err)
		}
		return checkDisplayMax(data.Count), nil
	}
	count, err := ibcTxRepo.CountTransferTxs(query)
	if err != nil {
		return 0, errors.Wrap(err)
	}
	return checkDisplayMax(count), nil
}

func (t TransferService) TransferTxs(req *vo.TranaferTxsReq) (vo.TranaferTxsResp, errors.Error) {
	var resp vo.TranaferTxsResp
	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	query, err := createIbcTxQuery(req)
	if err != nil {
		return resp, errors.Wrap(err)
	}
	if len(query.ChainId) > 2 {
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
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}

func (t TransferService) TransferTxDetail(hash string) (vo.TranaferTxDetailResp, errors.Error) {
	var resp vo.TranaferTxDetailResp
	ibcTxs, err := ibcTxRepo.TxDetail(hash, false)
	if err != nil && err != qmgo.ErrNoSuchDocuments {
		return resp, errors.Wrap(err)
	}
	if len(ibcTxs) == 0 {
		ibcTxs, err = ibcTxRepo.TxDetail(hash, true)
		if err != nil && err != qmgo.ErrNoSuchDocuments {
			return resp, errors.Wrap(err)
		}
	}
	setMap := make(map[string]struct{}, len(ibcTxs))
	for _, val := range ibcTxs {
		packetId := fmt.Sprintf("%s%s%s%s%s", val.ScPort, val.ScChannel, val.DcPort, val.DcChannel, val.Sequence)
		if _, exist := setMap[val.RecordId]; exist {
			continue
		}
		setMap[val.RecordId] = struct{}{}

		item := t.detailDto.LoadDto(val)
		if val.ScChainId != "" && val.ScTxInfo != nil && val.ScTxInfo.Hash != "" {
			item.ScConnect, item.ScSigners = getScTxInfo(val.ScChainId, val.ScTxInfo.Hash, packetId)
		}
		if val.DcChainId != "" && val.DcTxInfo != nil && val.DcTxInfo.Hash != "" {
			item.DcConnect, item.Ack, item.DcSigners = getDcTxInfo(val.DcChainId, val.DcTxInfo.Hash, packetId)
		}
		resp.Items = append(resp.Items, item)
	}
	if len(resp.Items) > 1 {
		// detail api page no support more than one return
		resp.Items = []vo.IbcTxDetailDto{}
	}
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}

func getMsgIndex(tx entity.Tx, msgType string, packetId string) int {
	for i, val := range tx.DocTxMsgs {
		if val.Type == msgType && val.CommonMsg().PacketId == packetId {
			return i
		}
	}
	return -1
}

func getScTxInfo(chainId string, txHash string, packetId string) (scConnect string, signers []string) {
	tx, err := txRepo.GetTxByHash(chainId, txHash)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	signers = tx.Signers
	scConnect = getConnectByTransferEventNews(tx.EventsNew, getMsgIndex(tx, constant.MsgTypeTransfer, packetId))
	return
}

func getDcTxInfo(chainId string, txHash string, packetId string) (dcConnect string, ack string, signers []string) {
	tx, err := txRepo.GetTxByHash(chainId, txHash)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	signers = tx.Signers
	dcConnect, ack = getConnectByRecvPacketEventsNews(tx.EventsNew, getMsgIndex(tx, constant.MsgTypeRecvPacket, packetId))
	return
}

func getConnectByTransferEventNews(eventNews []entity.EventNew, msgIndex int) string {
	var connect string
	for _, item := range eventNews {
		if item.MsgIndex == uint32(msgIndex) {
			for _, val := range item.Events {
				if val.Type == "send_packet" {
					for _, attribute := range val.Attributes {
						switch attribute.Key {
						case "packet_connection":
							connect = attribute.Value
						}
					}
				}
			}
		}
	}
	return connect
}

func getConnectByRecvPacketEventsNews(eventNews []entity.EventNew, msgIndex int) (string, string) {
	var connect, ackData string
	for _, item := range eventNews {
		if item.MsgIndex == uint32(msgIndex) {
			for _, val := range item.Events {
				if val.Type == "write_acknowledgement" || val.Type == "recv_packet" {
					for _, attribute := range val.Attributes {
						switch attribute.Key {
						case "packet_connection":
							connect = attribute.Value
						case "packet_ack":
							ackData = attribute.Value
						}
					}
				}
			}
		}
	}
	return connect, ackData
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
	relayerCfgMap, err := getRelayerCfgMap()
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
		matchInfo := strings.Join([]string{val.ScChainId, val.ScChannel, scRelayerAddr}, ":")
		if cfg, ok := relayerCfgMap[matchInfo]; ok {
			relayerInfo.ScRelayer.RelayerName = cfg.RelayerName
			relayerInfo.ScRelayer.Icon = cfg.Icon
		}
	}
	return &relayerInfo, nil
}
func getRelayerCfgMap() (map[string]entity.IBCRelayerConfig, error) {
	relayerCfgs, err := relayerCfgRepo.FindAll()
	if err != nil {
		return nil, err
	}
	relayerCfgMap := make(map[string]entity.IBCRelayerConfig, len(relayerCfgs))
	for _, val := range relayerCfgs {
		srcChainInfo := strings.Join([]string{val.ChainA, val.ChannelA, val.ChainAAddress}, ":")
		dcChainInfo := strings.Join([]string{val.ChainB, val.ChannelB, val.ChainBAddress}, ":")
		relayerCfgMap[srcChainInfo] = *val
		relayerCfgMap[dcChainInfo] = *val
	}
	return relayerCfgMap, nil
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
		denom, err := denomRepo.FindByDenomChainId(ibcTx.Denoms.ScDenom, ibcTx.ScChainId)
		if err != nil && err != qmgo.ErrNoSuchDocuments {
			return nil, err
		}
		if denom != nil {
			sendToken.DenomPath = strings.Join([]string{denom.DenomPath, denom.RootDenom}, "/")
		}
	}
	if strings.HasPrefix(ibcTx.Denoms.DcDenom, "ibc/") {
		denom, err := denomRepo.FindByDenomChainId(ibcTx.Denoms.DcDenom, ibcTx.DcChainId)
		if err != nil && err != qmgo.ErrNoSuchDocuments {
			return nil, err
		}
		if denom != nil {
			recvToken.DenomPath = strings.Join([]string{denom.DenomPath, denom.RootDenom}, "/")
		}
	}
	return &vo.TokenInfo{
		BaseDenom:        ibcTx.BaseDenom,
		BaseDenomChainId: ibcTx.BaseDenomChainId,
		Amount:           ibcTx.ScTxInfo.MsgAmount.Amount,
		SendToken:        sendToken,
		RecvToken:        recvToken,
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
	ibcBaseDenoms, err := baseDenomRepo.FindAll()
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
		return resp, errors.Wrapf("invalid hash")
	}

	msgType, ok := supportMsgType[req.MsgType]
	if !ok {
		return resp, errors.Wrapf("only support transfer,recv_packet,acknowledge_packet,timeout_packet")
	}

	value, err := lcdTxDataCache.Get(req.ChainId, hash)
	if err == nil {
		utils.UnmarshalJsonIgnoreErr([]byte(value), &resp)
		return resp, nil
	}
	return getMsgAndTxData(msgType, req.ChainId, hash)
}

func getMsgAndTxData(msgType, chainId, hash string) (vo.TraceSourceResp, errors.Error) {
	var resp vo.TraceSourceResp
	chainCfgData, err := chainCfgRepo.FindOne(chainId)
	if err != nil {
		return resp, errors.Wrap(err)
	}
	if !strings.Contains(chainCfgData.Lcd, "://") {
		return resp, errors.Wrap(fmt.Errorf("lcd from chain_config invalid for no include ://"))
	}

	values := strings.Split(chainCfgData.Lcd, "://")
	if len(values) < 2 || len(values[1]) == 0 {
		return resp, errors.Wrap(fmt.Errorf("lcd from chain_config invalid"))
	}

	lcdTxData, err := GetTxDataFromChain(chainCfgData.Lcd, hash)
	if err != nil {
		return resp, errors.Wrap(err)
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
		_ = lcdTxDataCache.Set(chainId, hash, string(utils.MarshalJsonIgnoreErr(resp)))
	}
	return resp, nil
}
func GetTxDataFromChain(lcdUri string, hash string) (LcdTxData, error) {
	url := fmt.Sprintf("%s/cosmos/tx/v1beta1/txs/%s", lcdUri, hash)
	data, err := utils.HttpGet(url)
	if err != nil {
		return LcdTxData{}, err
	}
	var txData LcdTxData
	if err := json.Unmarshal(data, &txData); err != nil {
		return LcdTxData{}, err
	}
	return txData, nil
}
