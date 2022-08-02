package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/qiniu/qmgo"
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

	//search ex_ibc_tx_latest collection
	ibcTxs, err := ibcTxRepo.FindAllByStatus(queryStats, req.PageNum, req.PageSize)
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
		case int(entity.IbcTxStatusFailed):
			if ibcTx.ScTxInfo.Status == int(entity.TxStatusSuccess) {
				item.ChainId = ibcTx.DcChainId
				item.TxHash = ibcTx.DcTxInfo.Hash
				if ibcTx.DcTxInfo.Status == int(entity.TxStatusSuccess) {
					item.TxErrorLog = "ack error"
				} else if ibcTx.DcTxInfo.Hash != "" {
					//find dc_chain_id sync_tx for log
					errLog, err := logCacheRepo.GetLogByHash(item.RecvChain, ibcTx.DcTxInfo.Hash)
					if err != nil {
						return nil, errors.Wrap(err)
					}
					item.TxErrorLog = errLog
				}
			} else {
				item.ChainId = ibcTx.ScChainId
				item.TxHash = ibcTx.ScTxInfo.Hash
				item.TxErrorLog = ibcTx.Log.ScLog
			}
			break
		case int(entity.IbcTxStatusRefunded):
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
	//search ex_ibc_tx_latest collection
	ibcTxs, err := ibcTxRepo.FindRelayerTxs(req.PageNum, req.PageSize)
	if err != nil && err != qmgo.ErrNoSuchDocuments {
		return nil, errors.Wrap(err)
	}

	//search ex_ibc_tx collection
	if len(ibcTxs) == 0 {
		ibcTxs, err = ibcTxRepo.FindHistoryRelayerTxs(req.PageNum, req.PageSize)
		if err != nil {
			return nil, errors.Wrap(err)
		}
	}
	items := make([]vo.RelayerTxFeeDto, 0, len(ibcTxs))
	for _, ibcTx := range ibcTxs {

		switch ibcTx.Status {
		case int(entity.IbcTxStatusFailed), int(entity.TxStatusSuccess):
			item := vo.RelayerTxFeeDto{
				ChainId:     ibcTx.DcChainId,
				Fee:         ibcTx.DcTxInfo.Fee,
				TxHash:      ibcTx.DcTxInfo.Hash,
				RelayerAddr: ibcTx.DcTxInfo.Msg.Msg.Signer,
			}
			items = append(items, item)
			break
		case int(entity.IbcTxStatusRefunded):
			item := vo.RelayerTxFeeDto{
				ChainId:     ibcTx.ScChainId,
				Fee:         ibcTx.RefundedTxInfo.Fee,
				TxHash:      ibcTx.RefundedTxInfo.Hash,
				RelayerAddr: ibcTx.RefundedTxInfo.Msg.Msg.Signer,
			}
			items = append(items, item)
			break
		}
	}
	totalItem := int64(len(items))
	page := vo.BuildPageInfo(totalItem, req.PageNum, req.PageSize)
	return &vo.RelayerTxFeesResp{
		Items:     items,
		PageInfo:  page,
		TimeStamp: time.Now().Unix(),
	}, nil
}
