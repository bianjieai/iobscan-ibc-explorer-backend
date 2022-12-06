package service

import (
	"encoding/json"
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/qiniu/qmgo"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	mainPageUrl           = "https://github.com/irisnet/iob-registry/tree/main/relayers"
	xpathRowHeader        = "//div[@role=\"rowheader\"]"
	rawInformationJsonUrl = "https://raw.githubusercontent.com/irisnet/iob-registry/main/relayers/%s/relayer_info.json"
	iconUrl               = "https://iobscan.io/resources/ibc-relayer/%s.png"
	parentPath            = ". ."
)

type RelayerHandler struct {
	chainIdNameMap map[string]string
}

func (h *RelayerHandler) Collect(filepath string) {
	logrus.Infof("RelayerHandler collect %s", filepath)
	st := time.Now().Unix()

	chainIdNameMapData, err := repository.GetChainIdNameMap()
	if err != nil {
		logrus.Errorf("getChainIdNameMap err, %v", err)
		return
	}
	h.chainIdNameMap = chainIdNameMapData

	if filepath == "" {
		h.xPathMainPage()
	} else {
		h.fetchSave(filepath)
	}

	logrus.Infof("RelayerHandler collect %s end, time use: %d(s)", filepath, time.Now().Unix()-st)
}

func (h *RelayerHandler) xPathMainPage() {
	doc, err := htmlquery.LoadURL(mainPageUrl)
	if err != nil {
		logrus.Errorf("RelayerHandler xPathMainPage err, %v", err)
		return
	}

	nodes, err := htmlquery.QueryAll(doc, xpathRowHeader)
	if err != nil {
		logrus.Errorf("RelayerHandler xPathMainPage xpath err, %v", err)
	}

	for _, v := range nodes {
		filepath := strings.TrimSpace(htmlquery.InnerText(v))
		if strings.HasSuffix(filepath, ".json") || strings.HasSuffix(filepath, ".md") || filepath == parentPath {
			continue
		}

		h.fetchSave(filepath)
	}
}

func (h *RelayerHandler) fetchSave(filepath string) {
	logrus.Infof("RelayerHandler fetchSave %s", filepath)
	relayerInfoResp, err := h.queryInfoJson(filepath)
	if err != nil {
		logrus.Infof("RelayerHandler queryInfoJson %s err, %v", filepath, err)
		return
	}

	if relayerInfoResp.TeamName == "" {
		logrus.Warningf("RelayerHandler relayer %s team name is blank", filepath)
		return
	}

	var distRelayerIds []string
	for _, addrMap := range relayerInfoResp.Addresses {
		index := 0
		var chainA, chainAAddress, chainB, chainBAddress string
		for k, v := range addrMap {
			chainName := k
			if len(h.chainIdNameMap) > 0 {
				if name, ok := h.chainIdNameMap[k]; ok {
					chainName = name
				}
			}
			if index == 0 {
				chainA = chainName
				chainAAddress = v
			} else {
				chainB = chainName
				chainBAddress = v
			}

			index++
		}
		distRelayerIds = append(distRelayerIds, entity.GenerateDistRelayerId(chainA, chainAAddress, chainB, chainBAddress))
	}

	if err = h.removeDumpChannelPairs(distRelayerIds); err != nil {
		logrus.Errorf("RelayerHandler removeDumpChannelPairs %s err, %v", relayerInfoResp.TeamName, err)
		return
	}

	if err = h.saveRegistryRelayer(relayerInfoResp.TeamName, distRelayerIds); err != nil {
		logrus.Errorf("RelayerHandler saveRegistryRelayer %s err, %v", relayerInfoResp.TeamName, err)
		return
	}
}

func (h *RelayerHandler) saveRegistryRelayer(relayerName string, nowDistRelayerIds []string) error {
	relayer, err := relayerRepo.FindOneByRelayerName(relayerName)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return h.insertNewRelayer(relayerName, nowDistRelayerIds)
		}
		return err
	}

	return h.updateRelayer(relayer, nowDistRelayerIds)
}

func (h *RelayerHandler) removeDumpChannelPairs(nowDistRelayerIds []string) error {
	for _, v := range nowDistRelayerIds {
		_, addressA, _, addressB := entity.ParseDistRelayerId(v)
		relayerList, err := relayerRepo.FindUnknownByAddrPair(addressA, addressB)
		if err != nil {
			return err
		}

		if len(relayerList) == 0 {
			continue
		}

		removeRelayerIds := make([]string, 0, len(relayerList))
		for _, re := range relayerList {
			removeRelayerIds = append(removeRelayerIds, re.RelayerId)
		}

		if err = relayerRepo.RemoveDumpData(removeRelayerIds); err != nil {
			return err
		}
	}
	return nil
}

func (h *RelayerHandler) insertNewRelayer(relayerName string, nowDistRelayerIds []string) error {
	nowChannelPairInfo := make([]entity.ChannelPairInfo, 0, len(nowDistRelayerIds))
	servedChainSet := utils.NewStringSet()
	for _, v := range nowDistRelayerIds {
		chainA, addressA, chainB, addressB := entity.ParseDistRelayerId(v)
		pairs, err := repository.GetChannelPairInfoByAddressPair(chainA, addressA, chainB, addressB)
		if err != nil {
			return err
		}
		nowChannelPairInfo = append(nowChannelPairInfo, pairs...)
		servedChainSet.AddAll(chainA, chainB)
	}

	err := relayerRepo.InsertOne(&entity.IBCRelayerNew{
		RelayerId:       primitive.NewObjectID().Hex(),
		RelayerName:     relayerName,
		RelayerIcon:     fmt.Sprintf(iconUrl, strings.ReplaceAll(relayerName, " ", "_")),
		ServedChains:    int64(servedChainSet.Len()),
		ChannelPairInfo: nowChannelPairInfo,
		CreateAt:        time.Now().Unix(),
		UpdateAt:        time.Now().Unix(),
	})

	if err != nil {
		logrus.Errorf("RelayerHandler insert relayer %s err, %v", relayerName, err)
		return err
	}
	logrus.Infof("RelayerHandler insert relayer %s succeed", relayerName)
	return nil
}

func (h *RelayerHandler) updateRelayer(relayer *entity.IBCRelayerNew, nowDistRelayerIds []string) error {
	var existedDistRelayerIds []string
	nowChannelPairInfoMap := make(map[string]entity.ChannelPairInfo, len(relayer.ChannelPairInfo))
	var needUpdate bool
	for _, v := range relayer.ChannelPairInfo {
		distRelayerId := entity.GenerateDistRelayerId(v.ChainA, v.ChainAAddress, v.ChainB, v.ChainBAddress)
		existedDistRelayerIds = append(existedDistRelayerIds, distRelayerId)
		if utils.InArray(nowDistRelayerIds, distRelayerId) {
			nowChannelPairInfoMap[v.PairId] = v
		} else { // pair removed
			needUpdate = true
		}
	}

	for _, v := range nowDistRelayerIds {
		if utils.InArray(existedDistRelayerIds, v) {
			continue
		}

		// 新增pair
		needUpdate = true
		chainA, addressA, chainB, addressB := entity.ParseDistRelayerId(v)
		pairs, err := repository.GetChannelPairInfoByAddressPair(chainA, addressA, chainB, addressB)
		if err != nil {
			return err
		}

		for _, p := range pairs {
			nowChannelPairInfoMap[p.PairId] = p
		}
	}

	if !needUpdate {
		logrus.Infof("RelayerHandler relayer %s don't chang", relayer.RelayerName)
		return nil
	}

	pairs := make([]entity.ChannelPairInfo, 0, len(nowChannelPairInfoMap))
	for _, v := range nowChannelPairInfoMap {
		pairs = append(pairs, v)
	}
	if err := relayerRepo.UpdateChannelPairInfo(relayer.RelayerId, pairs); err != nil {
		logrus.Errorf("RelayerHandler update relayer %s err, %v", relayer.RelayerName, err)
		return err
	} else {
		logrus.Infof("RelayerHandler update relayer %s succeed, %v", relayer.RelayerName, err)
		return nil
	}
}

func (h *RelayerHandler) queryInfoJson(filepath string) (*vo.IobRegistryRelayerInfoResp, error) {
	url := fmt.Sprintf(rawInformationJsonUrl, filepath)
	bz, err := utils.HttpGet(url)
	if err != nil {
		return nil, err
	}

	var res vo.IobRegistryRelayerInfoResp
	if err = json.Unmarshal(bz, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
