package task

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"strings"
)

type (
	RelayerStatisticsTask struct {
	}
	Statistic struct {
		RelayerId  string
		Amounts    decimal.Decimal
		Txs        int64
		TxsSuccess int64
	}
)

func (t *RelayerStatisticsTask) Name() string {
	return "ibc_relayer_statistics_task"
}

func (t *RelayerStatisticsTask) relayerTxsMapKey(chainId, dcChainAddr, dcChannel string) string {
	return fmt.Sprintf("%s:%s:%s", chainId, dcChannel, dcChainAddr)
}

func (t *RelayerStatisticsTask) Run() int {
	historySegments, err := getHistorySegment()
	if err != nil {
		logrus.Errorf("task %s getHistorySegment err, %v", t.Name(), err)
		return -1
	}
	//insert relayer data
	t.handleNewRelayerOnce(historySegments, true)

	logrus.Infof("task %s deal history segment total: %d", t.Name(), len(historySegments))
	if err = t.dealHistory(historySegments); err != nil {
		logrus.Errorf("task %s dealHistory err, %v", t.Name(), err)
		return -1
	}

	segments, err := getSegment()
	if err != nil {
		logrus.Errorf("task %s getSegment err, %v", t.Name(), err)
		return -1
	}
	//insert relayer data
	t.handleNewRelayerOnce(segments, false)

	logrus.Infof("task %s deal segment total: %d", t.Name(), len(segments))
	if err = t.deal(segments, opInsert); err != nil {
		logrus.Errorf("task %s deal err, %v", t.Name(), err)
		return -1
	}

	return 1
}

func (t *RelayerStatisticsTask) saveData(relayerStaticsMap map[string]Statistic, startTime, endTime int64, op int) error {
	var relayerStatics []entity.IBCRelayerStatistics
	for key, value := range relayerStaticsMap {
		if arrs := strings.Split(key, ":"); len(arrs) == 4 {
			chainId, baseDenom, _, channel := arrs[0], arrs[1], arrs[2], arrs[3]
			item := createIBCRelayerStatistics(channel, chainId, value.RelayerId, baseDenom, value.Amounts,
				value.TxsSuccess, value.Txs, startTime, endTime)
			relayerStatics = append(relayerStatics, item)
		}
	}
	if len(relayerStatics) > 0 {
		switch op {
		case opInsert:
			if err := relayerStatisticsRepo.Insert(relayerStatics); err != nil && !qmgo.IsDup(err) {
				return err
			}
		case opUpdate:
			for _, val := range relayerStatics {
				if err := relayerStatisticsRepo.InserOrUpdate(val); err != nil && err != qmgo.ErrNoSuchDocuments {
					logrus.Error("relayer statistic update fail, ", err.Error())
				}
			}

		}

	}
	return nil
}

// dealHistory 处理历史记录，针对ex_ibc_tx
func (t *RelayerStatisticsTask) dealHistory(segments []*segment) error {
	relayerTxs := make([]*dto.CountRelayerPacketTxsCntDTO, 0, 20)
	relayerSuccessTxs := make([]*dto.CountRelayerPacketTxsCntDTO, 0, 20)
	for _, v := range segments {
		relayerTxs = collectTxs(relayerTxs, v.StartTime, v.EndTime, ibcTxRepo.CountHistoryRelayerPacketTxs)
		relayerSuccessTxs = collectTxs(relayerSuccessTxs, v.StartTime, v.EndTime, ibcTxRepo.CountHistoryRelayerSuccessPacketTxs)
		relayerAmounts, err := ibcTxRepo.CountHistoryRelayerPacketAmount(v.StartTime, v.EndTime)
		if err != nil {
			logrus.Error(err.Error())
		}
		aggr := t.aggr(relayerTxs, relayerSuccessTxs, relayerAmounts)
		if err = t.saveData(aggr, v.StartTime, v.EndTime, opInsert); err != nil {
			return err
		}
		logrus.Debugf("dealHistory task %s scan ex_ibc_tx finish segment [%v:%v]", t.Name(), v.StartTime, v.EndTime)
	}
	return nil
}

// deal 处理最新的记录，针对ex_ibc_tx_latest
func (t *RelayerStatisticsTask) deal(segments []*segment, op int) error {
	relayerTxs := make([]*dto.CountRelayerPacketTxsCntDTO, 0, 20)
	relayerSuccessTxs := make([]*dto.CountRelayerPacketTxsCntDTO, 0, 20)
	for _, v := range segments {
		relayerTxs = collectTxs(relayerTxs, v.StartTime, v.EndTime, ibcTxRepo.CountRelayerPacketTxs)
		relayerSuccessTxs = collectTxs(relayerSuccessTxs, v.StartTime, v.EndTime, ibcTxRepo.CountRelayerSuccessPacketTxs)
		relayerAmounts, err := ibcTxRepo.CountRelayerPacketAmount(v.StartTime, v.EndTime)
		if err != nil {
			logrus.Error(err.Error())
		}
		aggr := t.aggr(relayerTxs, relayerSuccessTxs, relayerAmounts)
		if err := t.saveData(aggr, v.StartTime, v.EndTime, op); err != nil {
			return err
		}
		logrus.Debugf("deal task %s scan ex_ibc_tx_latest finish segment [%v:%v]", t.Name(), v.StartTime, v.EndTime)
	}
	return nil
}

func (t *RelayerStatisticsTask) aggr(relayerTxs, relayerSuccessTxs []*dto.CountRelayerPacketTxsCntDTO, relayerAmounts []*dto.CountRelayerPacketAmountDTO) map[string]Statistic {
	relayerTxsMap := make(map[string]TxsItem, 20)
	for _, tx := range relayerTxs {
		key := t.relayerTxsMapKey(tx.DcChainId, tx.DcChainAddress, tx.DcChannel)
		value, exist := relayerTxsMap[key]
		if exist {
			value.Txs += tx.Count
			relayerTxsMap[key] = value
		} else {
			relayerTxsMap[key] = TxsItem{Txs: tx.Count}
		}
	}

	for _, tx := range relayerSuccessTxs {
		key := t.relayerTxsMapKey(tx.DcChainId, tx.DcChainAddress, tx.DcChannel)
		value, exist := relayerTxsMap[key]
		if exist {
			value.TxsSuccess += tx.Count
			relayerTxsMap[key] = value
		} else {
			relayerTxsMap[key] = TxsItem{TxsSuccess: tx.Count}
		}
	}

	getRelayerTxs := func(data *entity.IBCRelayer, relayerTxsMap map[string]TxsItem) (int64, int64) {
		keyA := t.relayerTxsMapKey(data.ChainA, data.ChainAAddress, data.ChannelA)
		keyB := t.relayerTxsMapKey(data.ChainB, data.ChainBAddress, data.ChannelB)
		totalTxsAValue, _ := relayerTxsMap[keyA]
		totalTxsBValue, _ := relayerTxsMap[keyB]
		txsSuccess := totalTxsAValue.TxsSuccess + totalTxsBValue.TxsSuccess
		txs := totalTxsAValue.Txs + totalTxsBValue.Txs
		return txs, txsSuccess
	}

	relayerAmtsMap := createAmounts(relayerAmounts)
	relayerStaticsMap := make(map[string]Statistic, 20)
	for key, val := range relayerAmtsMap {
		if arrs := strings.Split(key, ":"); len(arrs) == 4 {
			chainId, _, relayerAddr, channel := arrs[0], arrs[1], arrs[2], arrs[3]
			relayerData, err := relayerRepo.FindRelayer(chainId, relayerAddr, channel)
			if err != nil {
				if err != qmgo.ErrNoSuchDocuments {
					logrus.Warn(chainId, relayerAddr, channel, "find relayer id fail, ", err.Error())
				}
				continue
			}
			txs, txsSuccess := getRelayerTxs(relayerData, relayerTxsMap)
			relayerStaticsMap[key] = Statistic{
				Amounts:    val,
				RelayerId:  relayerData.RelayerId,
				Txs:        txs,
				TxsSuccess: txsSuccess,
			}
		}
	}

	return relayerStaticsMap
}

func (t *RelayerStatisticsTask) handleNewRelayerOnce(segments []*segment, historyData bool) {
	for _, v := range segments {
		var relayersData []entity.IBCRelayer
		if historyData {
			relayersData = handleIbcTxHistory(v.StartTime, v.EndTime)
		} else {
			relayersData = handleIbcTxLatest(v.StartTime, v.EndTime)
		}
		if len(relayersData) > 0 {
			relayersData = distinctRelayer(relayersData)
			relayersData = filterDbExist(relayersData, historyData)
			if len(relayersData) == 0 {
				continue
			}
			if err := relayerRepo.Insert(relayersData); err != nil && !qmgo.IsDup(err) {
				logrus.Error("insert  relayer data fail, ", err.Error())
			}
		}
		logrus.Debugf("task %s find relayer finish segment [%v:%v], relayers:%v", t.Name(), v.StartTime, v.EndTime, len(relayersData))
	}
}

func handleIbcTxHistory(startTime, endTime int64) []entity.IBCRelayer {
	relayerDtos, err := ibcTxRepo.GetHistoryRelayerInfo(startTime, endTime)
	if err != nil {
		logrus.Errorf("get relayer info fail, %s", err.Error())
		return nil
	}
	var relayers []entity.IBCRelayer
	for _, val := range relayerDtos {
		relayers = append(relayers, createRelayerData(val))
	}
	return relayers
}
