package task

import (
	"fmt"
	"sync"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/sirupsen/logrus"
)

type FixBaseDenomChainIdTask struct {
	chainMap map[string]*entity.ChainConfig
	denomMap entity.IBCDenomMap
}

var _ OneOffTask = new(FixBaseDenomChainIdTask)

func (t *FixBaseDenomChainIdTask) Name() string {
	return "fix_base_denom_chain_id_task"
}

func (t *FixBaseDenomChainIdTask) Switch() bool {
	return global.Config.Task.SwitchFixBaseDenomChainIdTask
}

func (t *FixBaseDenomChainIdTask) Run() int {
	if err := t.init(); err != nil {
		return -1
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		t.handle(true)
		logrus.Infof("task %s handle history end", t.Name())
	}()

	go func() {
		defer wg.Done()
		t.handle(false)
		logrus.Infof("task %s handle latest end", t.Name())
	}()

	wg.Wait()
	return 1
}

func (t *FixBaseDenomChainIdTask) init() error {
	chainMap, err := getAllChainMap()
	if err != nil {
		logrus.Errorf("task %s getAllChainMap error, %v", t.Name(), err)
		return err
	}
	t.chainMap = chainMap

	denoms, err := denomRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s denomRepo.FindAll error, %v", t.Name(), err)
		return err
	}
	t.denomMap = denoms.ConvertToMap()
	return nil
}

func (t *FixBaseDenomChainIdTask) handle(isTargetHistory bool) {
	const limit = 1000
	var skip int64 = 0
	for {
		items, err := tokenStatisticsRepo.FindEmptyBaseDenomChainIdItems(skip, limit)
		if err != nil {
			logrus.Errorf("task %s FindEmptyBaseDenomChainIdItems error, %v", t.Name(), err)
			return
		}

		logrus.Infof("task %s find items skip: %d, limit: %d, target: %t", t.Name(), skip, limit, isTargetHistory)
		for _, v := range items {
			t.handleSegment(v.SegmentStartTime, v.SegmentEndTime, v.BaseDenom, isTargetHistory)
		}

		if len(items) < limit {
			break
		}
		skip += limit
	}
}

func (t *FixBaseDenomChainIdTask) handleSegment(startTime, endTime int64, baseDenom string, isTargetHistory bool) {
	ibcTxs, err := ibcTxRepo.FindByBaseDenom(startTime, endTime, baseDenom, "", isTargetHistory)
	if err != nil {
		logrus.Errorf("task %s FindByBaseDenom error, %v", t.Name(), err)
		return
	}

	for _, ibcTx := range ibcTxs {
		if ibcTx.Status == entity.IbcTxStatusFailed {
			continue
		}

		packetId := ibcTx.ScTxInfo.Msg.CommonMsg().PacketId
		tx, err := txRepo.GetTxByHash(ibcTx.ScChainId, ibcTx.ScTxInfo.Hash)
		if err != nil {
			logrus.Errorf("task %s GetTxByHash(hash: %s) error, %v", t.Name(), ibcTx.ScTxInfo.Hash, err)
			continue
		}

		for msgIndex, msg := range tx.DocTxMsgs {
			if msg.Type == constant.MsgTypeTransfer && msg.CommonMsg().PacketId == packetId {
				_, _, denomFullPath, _ := parseTransferTxEvents(msgIndex, &tx)

				ibcDenom := traceDenom(denomFullPath, ibcTx.ScChainId, t.chainMap)
				if err = ibcTxRepo.UpdateBaseDenom(ibcTx.RecordId, ibcDenom.BaseDenom, ibcDenom.BaseDenomChainId, isTargetHistory); err != nil {
					logrus.Errorf("task %s UpdateBaseDenom(recordId: %s) error, %v", t.Name(), ibcTx.RecordId, err)
				}

				if ibcTx.Status == entity.IbcTxStatusProcessing || ibcTx.Status == entity.IbcTxStatusSuccess || ibcTx.Status == entity.IbcTxStatusRefunded {
					t.upsertDenom(ibcDenom)
				}

				if ibcTx.Status == entity.IbcTxStatusSuccess {
					dcDenomFullPath, isCrossBack := calculateNextDenomPath(ibcTx.DcTxInfo.Msg.RecvPacketMsg().Packet)
					if !isCrossBack { // transfer to next chain
						dcDenomPath, rootDenom := splitFullPath(dcDenomFullPath)
						dcDenom := &entity.IBCDenom{
							Symbol:           "",
							ChainId:          ibcTx.DcChainId,
							Denom:            ibcTx.Denoms.DcDenom,
							PrevDenom:        ibcTx.Denoms.ScDenom,
							PrevChainId:      ibcTx.ScChainId,
							BaseDenom:        ibcDenom.BaseDenom,
							BaseDenomChainId: ibcDenom.BaseDenomChainId,
							DenomPath:        dcDenomPath,
							RootDenom:        rootDenom,
							IsBaseDenom:      false,
							CreateAt:         time.Now().Unix(),
							UpdateAt:         time.Now().Unix(),
						}

						t.upsertDenom(dcDenom)
					}
				}
			}

			break
		}
	}
}

func (t *FixBaseDenomChainIdTask) upsertDenom(ibcDenom *entity.IBCDenom) {
	_, ok := t.denomMap[fmt.Sprintf("%s%s", ibcDenom.ChainId, ibcDenom.Denom)]
	if !ok {
		if err := denomRepo.Insert(ibcDenom); err != nil {
			logrus.Errorf("task %s denomRepo.Insert error, chain_id: %s, denom: %s, %v", t.Name(), ibcDenom.ChainId, ibcDenom.Denom, err)
		}
	} else {
		if err := denomRepo.UpdateDenom(ibcDenom); err != nil {
			logrus.Errorf("task %s denomRepo.UpdateDenom error, chain_id: %s, denom: %s, %v", t.Name(), ibcDenom.ChainId, ibcDenom.Denom, err)
		}
	}
}
