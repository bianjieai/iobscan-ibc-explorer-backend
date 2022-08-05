package service

import (
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

	if req.Symbol != "" {
		if req.Symbol == constant.UnAuth {
			tokens, err := getUnAuthToken()
			if err != nil {
				return query, err
			}
			query.Token = tokens

		} else {
			baseDenom, err := baseDenomRepo.FindBySymbol(req.Symbol)
			if err != nil {
				return query, err
			}
			query.Token = []string{baseDenom.Denom}
		}
	} else if req.Denom != "" {
		query.Token = []string{req.Denom}
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
	count, err := ibcTxRepo.CountTransferTxs(query)
	if err != nil {
		return 0, errors.Wrap(err)
	}
	return count, nil
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
	total := int64(len(items))
	page := vo.BuildPageInfo(total, req.PageNum, req.PageSize)
	resp.PageInfo = page
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}

func (t TransferService) TransferTxDetail(hash string) (vo.TranaferTxDetailResp, errors.Error) {
	var resp vo.TranaferTxDetailResp
	ibcTxs, err := ibcTxRepo.TxDetail(hash)
	if err != nil && err != qmgo.ErrNoSuchDocuments {
		return resp, errors.Wrap(err)
	}
	if len(ibcTxs) == 0 {
		ibcTxs, err = ibcTxRepo.HistoryTxDetail(hash)
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
			item.ScConnection, item.TimeoutTimestamp, item.ScSigners = getScTxInfo(val.ScChainId, val.ScTxInfo.Hash, packetId)
		}
		if val.DcChainId != "" && val.DcTxInfo != nil && val.DcTxInfo.Hash != "" {
			item.DcConnection, item.Ack, item.DcSigners = getDcTxInfo(val.DcChainId, val.DcTxInfo.Hash, packetId)
		}
		resp.Items = append(resp.Items, item)
	}
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

func getScTxInfo(chainId string, txHash string, packetId string) (scConnect string, timeOutTimestamp string, signers []string) {
	tx, err := txRepo.GetTxByHash(chainId, txHash)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	signers = tx.Signers
	scConnect, timeOutTimestamp = getConnectByTransferEventNews(tx.EventsNew, getMsgIndex(tx, constant.MsgTypeTransfer, packetId))
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

func getConnectByTransferEventNews(eventNews []entity.EventNew, msgIndex int) (string, string) {
	var connect, timeoutTimestamp string
	for _, item := range eventNews {
		if item.MsgIndex == uint32(msgIndex) {
			for _, val := range item.Events {
				if val.Type == "send_packet" {
					for _, attribute := range val.Attributes {
						switch attribute.Key {
						case "packet_connection":
							connect = attribute.Value
						case "packet_timeout_timestamp":
							timeoutTimestamp = attribute.Value
						}
					}
				}
			}
		}
	}
	return connect, timeoutTimestamp
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
