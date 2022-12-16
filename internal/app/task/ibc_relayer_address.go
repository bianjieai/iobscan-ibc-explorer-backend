package task

import (
	"fmt"
	"strings"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/pkg/lcd"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			if account, err := lcd.GetAccount(v.Chain, v.Address, cf.GrpcRestGateway, cf.LcdApiPath.AccountsPath, false); err == nil {
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

// ==================================================================
// ==================================================================
// ==================================================================

var relayerAddressGatherTask RelayerAddressGatherTask

type RelayerAddressGatherTask struct {
	chainMap map[string]*entity.ChainConfig
}

func (t *RelayerAddressGatherTask) Name() string {
	return "relayer_address_gather_task"
}

func (t *RelayerAddressGatherTask) Run() int {
	chainMap, err := getAllChainInfosMap()
	if err != nil {
		logrus.Errorf("task %s getAllChainInfosMap err, %v", t.Name(), err)
		return -1
	}
	t.chainMap = chainMap

	t.repairEmptyPubKey()
	t.gather()

	return 1
}

// repairEmptyPubKey 修复pub_key 为空的address
func (t *RelayerAddressGatherTask) repairEmptyPubKey() {
	logrus.Infof("task %s repairEmptyPubKey start", t.Name())
	st := time.Now().Unix()
	startTime := st - relayerAddressGatherRangeTime
	addresses, err := relayerAddressRepo.FindNoPubKey(startTime)
	if err != nil {
		logrus.Errorf("task %s relayerAddressRepo.FindNoPubKey err, %v", t.Name(), err)
		return
	}

	fastFailChainMap := make(map[string]string)
	for _, v := range addresses {
		if _, ok := fastFailChainMap[v.Chain]; ok {
			continue
		}

		var apiHost, apiPath string
		if cfg, ok := t.chainMap[v.Chain]; ok {
			apiHost = cfg.GrpcRestGateway
			apiPath = cfg.LcdApiPath.AccountsPath
		} else {
			continue
		}

		account, err := lcd.GetAccount(v.Chain, v.Address, apiHost, apiPath, false)
		if err != nil {
			if t.isFastFailErr(err) {
				fastFailChainMap[v.Chain] = err.Error()
			}

			continue
		}

		if err = relayerAddressRepo.UpdatePubKey(v.Address, v.Chain, account.Account.PubKey.Key); err != nil {
			logrus.Errorf("task %s UpdatePubKey err, %v", t.Name(), err)
		}
	}

	logrus.Infof("task %s repairEmptyPubKey end, time use: %d[s]", t.Name(), time.Now().Unix()-st)
}

func (t *RelayerAddressGatherTask) isFastFailErr(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "connection refused") || strings.Contains(errStr, "i/o timeout") ||
		strings.Contains(errStr, "unsupported protocol scheme") || strings.Contains(errStr, "501") ||
		strings.Contains(errStr, "429")
}

func (t *RelayerAddressGatherTask) gather() {
	logrus.Errorf("task %s gather start", t.Name())
	st := time.Now().Unix()
	startTime := st - relayerAddressGatherRangeTime
	todoList, err := relayerAddressRepo.FindToBeGathered(startTime)
	if err != nil {
		logrus.Errorf("task %s FindToBeGathered err, %v", t.Name(), err)
		return
	}

	for _, v := range todoList {
		// 找到公钥相同的address
		samePubKeyList, err := relayerAddressRepo.FindByPubKey(v.PubKey)
		if err != nil {
			continue
		}

		var registeredRelayer, unknownRelayer *entity.IBCRelayerNew
		for _, s := range samePubKeyList {
			if s.Chain == v.Chain && s.Address == v.Address {
				continue
			}
			switch s.GatherStatus {
			case entity.GatherStatusPubKey:
				if unknownRelayer == nil {
					res, err := relayerRepo.FindByChannelPairChainA(s.Chain, s.Address)
					if err == nil {
						unknownRelayer = res
					} else {
						res, err = relayerRepo.FindByChannelPairChainB(s.Chain, s.Address)
						if err == nil {
							unknownRelayer = res
						}
					}
				}
			case entity.GatherStatusRegistry:
				if registeredRelayer == nil {
					res, err := relayerRepo.FindByChannelPairChainA(s.Chain, s.Address)
					if err == nil {
						registeredRelayer = res
					} else {
						res, err = relayerRepo.FindByChannelPairChainB(s.Chain, s.Address)
						if err == nil {
							registeredRelayer = res
						}
					}
				}
			default:
				continue
			}
		}

		// 将address 归档到与其公钥相同的地址所在的relayer中，如果没有公钥相同的地址，则此地址独立作为一个单项的relayer
		if registeredRelayer != nil {
			t.addRelayerChannelPair(registeredRelayer, v.Address, v.Chain, entity.GatherStatusRegistry)
		} else if unknownRelayer != nil {
			t.addRelayerChannelPair(unknownRelayer, v.Address, v.Chain, entity.GatherStatusPubKey)
		} else {
			t.addRelayer(v.Address, v.Chain)
		}
	}

	logrus.Infof("task %s gather end, time use: %d[s]", t.Name(), time.Now().Unix()-st)
}

func (t *RelayerAddressGatherTask) getNewPair(address, chain string) ([]entity.ChannelPairInfo, error) {
	addressChainList, err := relayerAddressChannelRepo.FindByAddressChain(address, chain)
	if err != nil {
		logrus.Errorf("task %s relayerAddressChannelRepo.FindByAddressChain err, %v", t.Name(), err)
		return nil, err
	}

	if len(addressChainList) == 0 {
		return nil, fmt.Errorf("no address channel")
	}

	newPair := make([]entity.ChannelPairInfo, 0, len(addressChainList))
	for _, v := range addressChainList {
		newPair = append(newPair, entity.GenerateSingleSideChannelPairInfo(chain, v.Channel, address))
	}

	return newPair, nil
}

func (t *RelayerAddressGatherTask) addRelayerChannelPair(relayer *entity.IBCRelayerNew, address, chain string, status entity.GatherStatus) {
	newPair, err := t.getNewPair(address, chain)
	if err != nil {
		return
	}

	relayer.ChannelPairInfo = append(relayer.ChannelPairInfo, newPair...)
	if err = relayerRepo.UpdateChannelPairInfo(relayer.RelayerId, relayer.ChannelPairInfo); err != nil {
		logrus.Errorf("task %s relayerRepo.UpdateChannelPairInfo err, %v", t.Name(), err)
		return
	}

	if err = relayerAddressRepo.UpdateGatherStatus(address, chain, status); err != nil {
		logrus.Errorf("task %s UpdateGatherStatus %s-%s err, %v", t.Name(), address, chain, err)
	}
}

func (t *RelayerAddressGatherTask) addRelayer(address, chain string) {
	newPair, err := t.getNewPair(address, chain)
	if err != nil {
		return
	}

	now := time.Now().Unix()
	newRelayer := &entity.IBCRelayerNew{
		RelayerId:            primitive.NewObjectID().Hex(),
		RelayerName:          "",
		RelayerIcon:          "",
		ServedChains:         1,
		ChannelPairInfo:      newPair,
		UpdateTime:           0,
		RelayedTotalTxs:      0,
		RelayedSuccessTxs:    0,
		RelayedTotalTxsValue: "0",
		TotalFeeValue:        "0",
		CreateAt:             now,
		UpdateAt:             now,
	}

	if err = relayerRepo.InsertOne(newRelayer); err != nil {
		logrus.Errorf("task %s insert relayer err, %v", t.Name(), err)
	}

	if err = relayerAddressRepo.UpdateGatherStatus(address, chain, entity.GatherStatusPubKey); err != nil {
		logrus.Errorf("task %s UpdateGatherStatus %s-%s err, %v", t.Name(), address, chain, err)
	}
}
