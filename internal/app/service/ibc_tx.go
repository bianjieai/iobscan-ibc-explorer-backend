package service

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/qiniu/qmgo"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type IbcTxServerI interface {
	ListFailTxs(req *vo.FailTxsListReq) (*vo.FailTxsListResp, errors.Error)
	ListRelayerTxFees(req *vo.RelayerTxFeesReq) (*vo.RelayerTxFeesResp, errors.Error)
}

type IbcTxService struct {
}

func (svc *IbcTxService) ListFailTxs(req *vo.FailTxsListReq) (*vo.FailTxsListResp, errors.Error) {
	queryStats := []entity.IbcTxStatus{entity.IbcTxStatusFailed, entity.IbcTxStatusRefunded}

	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	//search ex_ibc_tx_latest collection
	ibcTxs, err := ibcTxRepo.FindAllByStatus(queryStats, skip, limit)
	if err != nil && err != qmgo.ErrNoSuchDocuments {
		return nil, errors.Wrap(err)
	}

	//search ex_ibc_tx collection
	if len(ibcTxs) == 0 {
		ibcTxs, err = ibcTxRepo.FindAllHistoryByStatus(queryStats, req.PageNum, req.PageSize)
		if err != nil {
			return nil, errors.Wrap(err)
		}
	}

	items := make([]vo.FailTxsListDto, 0, len(ibcTxs))
	for _, ibcTx := range ibcTxs {
		item := vo.FailTxsListDto{
			SendChain: ibcTx.ScChainId,
			RecvChain: ibcTx.DcChainId,
		}
		switch ibcTx.Status {
		case entity.IbcTxStatusFailed:
			if ibcTx.ScTxInfo.Status == entity.TxStatusSuccess {
				item.ChainId = ibcTx.DcChainId
				item.TxHash = ibcTx.DcTxInfo.Hash
				item.TxErrorLog = getDcTxLog(ibcTx.DcChainId, ibcTx.DcTxInfo)
			} else {
				item.ChainId = ibcTx.ScChainId
				item.TxHash = ibcTx.ScTxInfo.Hash
				item.TxErrorLog = ibcTx.Log.ScLog
			}
			break
		case entity.IbcTxStatusRefunded:
			item.ChainId = ibcTx.ScChainId
			item.TxHash = ibcTx.RefundedTxInfo.Hash
			break
		}
		items = append(items, item)
	}
	totalItem := int64(len(items))
	page := vo.BuildPageInfo(totalItem, req.PageNum, req.PageSize)
	return &vo.FailTxsListResp{
		Items:     items,
		PageInfo:  page,
		TimeStamp: time.Now().Unix(),
	}, nil
}

func (svc *IbcTxService) ListRelayerTxFees(req *vo.RelayerTxFeesReq) (*vo.RelayerTxFeesResp, errors.Error) {
	//search sync_{chain_id}_tx collection
	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	var (
		rets []vo.RelayerTxFeeDto
		err  error
	)
	if req.TxHash != "" && req.ChainId == "" {
		return nil, errors.Wrap(fmt.Errorf("chainId is invalid"))
	}
	getChainRelayerTxs := func(chainId, txHash string) ([]vo.RelayerTxFeeDto, error) {
		var (
			txs []entity.Tx
		)

		if chainId != "" && txHash == "" {
			txs, err = txRepo.GetRelayerTxs(chainId, skip, limit)
			if err != nil && err != qmgo.ErrNoSuchDocuments {
				return nil, errors.Wrap(err)
			}
		} else if chainId != "" && txHash != "" {
			tx, err := txRepo.GetTxByHash(chainId, txHash)
			if err != nil && err != qmgo.ErrNoSuchDocuments {
				return nil, errors.Wrap(err)
			}
			txs = append(txs, tx)
		}

		items := make([]vo.RelayerTxFeeDto, 0, len(txs))
		for _, tx := range txs {
			item := vo.RelayerTxFeeDto{
				ChainId:     chainId,
				Fee:         tx.Fee,
				TxHash:      tx.TxHash,
				RelayerAddr: tx.Signers[0],
			}
			items = append(items, item)
		}
		return items, nil
	}
	var chainIds []string
	if req.ChainId != "" {
		chainIds = strings.Split(req.ChainId, ",")
	} else {
		chainCfg, err := chainConfigRepo.FindAllChainIds()
		if err != nil {
			return nil, errors.Wrap(err)
		}
		for _, val := range chainCfg {
			chainIds = append(chainIds, val.ChainId)
		}
	}

	for _, chainId := range chainIds {
		items, err := getChainRelayerTxs(chainId, req.TxHash)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		rets = append(rets, items...)
	}

	totalItem := int64(len(rets))
	page := vo.BuildPageInfo(totalItem, req.PageNum, req.PageSize)
	return &vo.RelayerTxFeesResp{
		Items:     rets,
		PageInfo:  page,
		TimeStamp: time.Now().Unix(),
	}, nil
}

func getDcTxLog(chainId string, txInfo *entity.TxInfo) string {
	if txInfo == nil {
		return ""
	}

	value, _ := logCacheRepo.Get(chainId, txInfo.Hash)
	if len(value) > 0 {
		return value
	}
	tx, err := txRepo.GetTxByHash(chainId, txInfo.Hash)
	if err != nil {
		logrus.Error("GetTxByHash " + err.Error())
		return ""
	}
	if txInfo.Status == entity.TxStatusSuccess {
		txErrorLog := "ack error"
		if len(tx.EventsNew) > 0 {
			txErrorLog = getConnectByRecvPacketEventsNews(tx.EventsNew, getMsgIndex(tx, constant.MsgTypeRecvPacket, txInfo.Msg.CommonMsg().PacketId))
		}
		logCacheRepo.Set(chainId, txInfo.Hash, txErrorLog)
		return txErrorLog
	} else if txInfo.Hash != "" {
		//find dc_chain_id sync_tx for log
		logCacheRepo.Set(chainId, txInfo.Hash, tx.Log)
		return tx.Log
	}
	return ""
}

func getConnectByRecvPacketEventsNews(eventNews []entity.EventNew, msgIndex int) string {
	var ackData string
	for _, item := range eventNews {
		if item.MsgIndex == uint32(msgIndex) {
			for _, val := range item.Events {
				if val.Type == "write_acknowledgement" || val.Type == "recv_packet" {
					for _, attribute := range val.Attributes {
						switch attribute.Key {
						case "packet_ack":
							ackData = attribute.Value
						}
					}
				}
			}
		}
	}
	return ackData
}

func getMsgIndex(tx entity.Tx, msgType string, packetId string) int {
	for i, val := range tx.DocTxMsgs {
		if val.Type == msgType && val.CommonMsg().PacketId == packetId {
			return i
		}
	}
	return -1
}
