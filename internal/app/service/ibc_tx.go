package service

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/qiniu/qmgo"
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
			SendChain: ibcTx.ScChain,
			RecvChain: ibcTx.DcChain,
		}
		switch ibcTx.Status {
		case entity.IbcTxStatusFailed:
			if ibcTx.ScTxInfo.Status == entity.TxStatusSuccess {
				item.Chain = ibcTx.DcChain
				item.TxHash = ibcTx.DcTxInfo.Hash
				item.TxErrorLog = ibcTx.DcTxInfo.Log
			} else {
				item.Chain = ibcTx.ScChain
				item.TxHash = ibcTx.ScTxInfo.Hash
				item.TxErrorLog = ibcTx.ScTxInfo.Log
			}
			break
		case entity.IbcTxStatusRefunded:
			item.Chain = ibcTx.ScChain
			item.TxHash = ibcTx.AckTimeoutTxInfo.Hash
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
	//search sync_{chain}_tx collection
	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	var (
		rets []vo.RelayerTxFeeDto
		err  error
	)
	if req.TxHash != "" && req.Chain == "" {
		return nil, errors.Wrap(fmt.Errorf("chainId is invalid"))
	}
	getChainRelayerTxs := func(chain, txHash string) ([]vo.RelayerTxFeeDto, error) {
		var (
			txs []entity.Tx
		)

		if chain != "" && txHash == "" {
			txs, err = txRepo.GetRelayerTxs(chain, skip, limit)
			if err != nil && err != qmgo.ErrNoSuchDocuments {
				return nil, errors.Wrap(err)
			}
		} else if chain != "" && txHash != "" {
			tx, err := txRepo.GetTxByHash(chain, txHash)
			if err != nil && err != qmgo.ErrNoSuchDocuments {
				return nil, errors.Wrap(err)
			}
			txs = append(txs, tx)
		}

		items := make([]vo.RelayerTxFeeDto, 0, len(txs))
		for _, tx := range txs {
			item := vo.RelayerTxFeeDto{
				Chain:       chain,
				Fee:         tx.Fee,
				TxHash:      tx.TxHash,
				RelayerAddr: tx.Signers[0],
			}
			items = append(items, item)
		}
		return items, nil
	}
	var chains []string
	if req.Chain != "" {
		chains = strings.Split(req.Chain, ",")
	} else {
		chainCfg, err := chainConfigRepo.FindAllChains()
		if err != nil {
			return nil, errors.Wrap(err)
		}
		for _, val := range chainCfg {
			chains = append(chains, val.ChainName)
		}
	}

	for _, chain := range chains {
		items, err := getChainRelayerTxs(chain, req.TxHash)
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
