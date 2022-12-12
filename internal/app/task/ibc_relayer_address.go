package task

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/pkg/lcd"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/sirupsen/logrus"
)

var relayerAddressInitTask IbcRelayerAddressInitTask

type IbcRelayerAddressInitTask struct {
}

func (t *IbcRelayerAddressInitTask) Name() string {
	return "ibc_relayer_address_init_task"
}

func (t *IbcRelayerAddressInitTask) Run() int {
	addrList, err := relayerAddressChannelRepo.DistinctAddr()
	if err != nil {
		logrus.Errorf("task %s DistinctAddr err, %v", t.Name(), err)
		return -1
	}

	chainInfosMap, err := getAllChainInfosMap()
	if err != nil {
		logrus.Errorf("task %s getAllChainInfosMap err, %v", t.Name(), err)
		return -1
	}

	authedRelayerList, err := relayerRepo.FindAuthed()
	if err != nil {
		logrus.Errorf("task %s relayerRepo.FindAuthed err, %v", t.Name(), err)
		return -1
	}

	authedAddrMap := make(map[string]struct{}, len(authedRelayerList))
	for _, relayer := range authedRelayerList {
		for _, pair := range relayer.ChannelPairInfo {
			authedAddrMap[fmt.Sprintf("%s:%s", pair.ChainA, pair.ChainAAddress)] = struct{}{}
			authedAddrMap[fmt.Sprintf("%s:%s", pair.ChainB, pair.ChainBAddress)] = struct{}{}
		}
	}

	addrEntityList := make([]*entity.IBCRelayerAddress, 0, len(addrList))
	nowTime := time.Now().Unix()
	for _, v := range addrList {
		gatherStatus := entity.GatherStatusTODO
		if _, ok := authedAddrMap[fmt.Sprintf("%s:%s", v.Chain, v.Address)]; ok {
			gatherStatus = entity.GatherStatusRegistry
		}

		var pubKey string
		if cf, ok := chainInfosMap[v.Chain]; ok {
			if account, err := lcd.GetAccount(v.Chain, v.Address, cf.GrpcRestGateway, cf.LcdApiPath.AccountsPath); err == nil {
				pubKey = account.Account.PubKey.Key
			}
		}

		addrEntityList = append(addrEntityList, &entity.IBCRelayerAddress{
			Address:      v.Address,
			Chain:        v.Chain,
			PubKey:       pubKey,
			GatherStatus: gatherStatus,
			CreateAt:     nowTime,
			UpdateAt:     nowTime,
		})
	}

	if err = relayerAddressRepo.InsertMany(addrEntityList); err != nil {
		logrus.Errorf("task %s InsertMany err, %v", t.Name(), err)
	}

	return 1
}
