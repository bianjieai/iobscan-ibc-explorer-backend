package service

import (
	"encoding/json"
	"fmt"
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
	teamNameMap  map[string]string
	jsonFileList []string
	pairIdMap    map[string]struct{}
	imageMap     map[string]bool
}

func (h *RelayerHandler) Collect(filepath string) {
	logrus.Infof("RelayerHandler collect %s", filepath)
	st := time.Now().Unix()
	h.teamNameMap = make(map[string]string)
	h.pairIdMap = make(map[string]struct{})
	h.imageMap = make(map[string]bool)
	h.jsonFileList = nil

	ids, err := relayerCfgRepo.FindRelayerPairIds()
	if err != nil {
		logrus.Infof("RelayerHandler FindRelayerPairIds err, %v", err)
		return
	}
	for _, v := range ids {
		h.pairIdMap[v.RelayerPairId] = struct{}{}
	}

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

	var distRelayerIds []string
	for _, addrMap := range relayerInfoResp.Addresses {
		index := 0
		var chainA, chainAAddress, chainB, chainBAddress string
		for k, v := range addrMap {
			if index == 0 {
				chainA = k
				chainAAddress = v
			} else {
				chainB = k
				chainBAddress = v
			}

			index++
		}
		distRelayerIds = append(distRelayerIds, entity.GenerateDistRelayerId(chainA, chainAAddress, chainB, chainBAddress))
	}

	if err = h.saveRegistryRelayer(relayerInfoResp.TeamName, distRelayerIds); err != nil {
		logrus.Errorf("RelayerHandler saveRegistryRelayer %s err, %v", relayerInfoResp.TeamName, err)
		return
	}

	if err = h.removeDumpChannelPairs(distRelayerIds); err != nil {
		logrus.Errorf("RelayerHandler removeDumpChannelPairs %s err, %v", relayerInfoResp.TeamName, err)
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
		pairs, err := getChannelPairInfoByAddressPair(chainA, addressA, chainB, addressB)
		if err != nil {
			return err
		}
		nowChannelPairInfo = append(nowChannelPairInfo, pairs...)
		servedChainSet.AddAll(chainA, chainB)
	}

	err := relayerRepo.InsertOne(&entity.IBCRelayerNew{
		RelayerId:       primitive.NewObjectID().Hex(),
		RelayerName:     relayerName,
		RelayerIcon:     strings.ReplaceAll(relayerName, " ", "_"),
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
		pairs, err := getChannelPairInfoByAddressPair(chainA, addressA, chainB, addressB)
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

// getChannelPairInfoByAddressPair 获取一对地址上的所有channel pair
func getChannelPairInfoByAddressPair(chainA, addressA, chainB, addressB string) ([]entity.ChannelPairInfo, error) {
	addrChannels, err := relayerAddrChannelRepo.FindChannels([]string{addressA, addressB})
	if err != nil {
		return nil, err
	}

	chainAChannelMap := make(map[string]string)
	chainBChannelMap := make(map[string]string)
	for _, c := range addrChannels {
		if c.RelayerAddress == addressA && c.Chain == chainA {
			chainAChannelMap[c.Channel] = c.CounterPartyChannel
		} else if c.RelayerAddress == addressB && c.Chain == chainB {
			chainBChannelMap[c.Channel] = c.CounterPartyChannel
		}
	}

	var res []entity.ChannelPairInfo
	var channelMatched bool
	for ch, cpch := range chainAChannelMap {
		if ch2, _ := chainBChannelMap[cpch]; ch == ch2 { // channel match success
			pairInfo := entity.GenerateChannelPairInfo(chainA, ch, addressA, chainB, cpch, addressB)
			channelMatched = true
			res = append(res, pairInfo)
		}
	}

	if !channelMatched {
		pairInfo := entity.GenerateChannelPairInfo(chainA, "", addressA, chainB, "", addressB)
		res = append(res, pairInfo)
	}

	return res, nil
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
