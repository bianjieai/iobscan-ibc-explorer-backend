package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
)

const (
	mainPageUrl           = "https://github.com/irisnet/iob-registry/tree/main/relayers"
	subPageUrl            = "https://github.com/irisnet/iob-registry/tree/main/relayers/%s"
	xpathRowHeader        = "//div[@role=\"rowheader\"]"
	rawJsonUrl            = "https://raw.githubusercontent.com/irisnet/iob-registry/main/relayers/%s/%s"
	rawInformationJsonUrl = "https://raw.githubusercontent.com/irisnet/iob-registry/main/relayers/%s/information.json"
	informationFileName   = "information.json"
	iconUrl               = "https://iobscan.io/resources/ibc-relayer/%s.png"
	concurrentNum         = 3
	parentPath            = ". ."
	imagesFileName        = "images"
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
		h.xPathSubPage(filepath)
	}

	if len(h.jsonFileList) > 0 {
		var wg sync.WaitGroup
		wg.Add(concurrentNum)
		for i := 0; i < concurrentNum; i++ {
			seq := i
			go func() {
				defer wg.Done()
				h.fetchAndSave(seq)
			}()
		}
		wg.Wait()
	}

	_, _ = relayerCfgRepo.(*cache.RelayerConfigCacheRepo).DelCacheFindRelayerPairIds()
	_, _ = relayerCfgRepo.(*cache.RelayerConfigCacheRepo).DelCacheFindAll()
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

		h.xPathSubPage(filepath)
	}
}

func (h *RelayerHandler) xPathSubPage(filepath string) {
	if _, err := h.queryInfoJson(filepath); err != nil {
		return
	}

	doc, err := htmlquery.LoadURL(fmt.Sprintf(subPageUrl, filepath))
	if err != nil {
		logrus.Errorf("RelayerHandler xPathSubPage err, %v", err)
		return
	}

	nodes, err := htmlquery.QueryAll(doc, xpathRowHeader)
	if err != nil {
		logrus.Errorf("RelayerHandler xPathSubPage xpath err, %v", err)
	}

	for _, v := range nodes {
		jsonFile := strings.TrimSpace(htmlquery.InnerText(v))
		if strings.HasSuffix(jsonFile, ".json") && jsonFile != informationFileName {
			h.jsonFileList = append(h.jsonFileList, fmt.Sprintf("%s|%s", filepath, jsonFile))
		}

		if jsonFile == imagesFileName {
			h.imageMap[filepath] = true
		}
	}
}

func (h *RelayerHandler) queryInfoJson(filepath string) (*vo.IobRegistryRelayerInfoResp, error) {
	url := fmt.Sprintf(rawInformationJsonUrl, filepath)
	bz, err := utils.HttpGet(url)
	if err != nil {
		logrus.Errorf("RelayerHandler queryInfoJson err, %v", err)
		return nil, nil
	}

	var res vo.IobRegistryRelayerInfoResp
	_ = json.Unmarshal(bz, &res)
	h.teamNameMap[filepath] = res.TeamName
	return &res, nil
}

func (h *RelayerHandler) queryPairJson(filepath, filename string) (*vo.IobRegistryRelayerPairResp, error) {
	url := fmt.Sprintf(rawJsonUrl, filepath, filename)
	bz, err := utils.HttpGet(url)
	if err != nil {
		logrus.Errorf("RelayerHandler queryChannelPairJson(%s) err, %v", url, err)
		return nil, nil
	}

	var res vo.IobRegistryRelayerPairResp
	_ = json.Unmarshal(bz, &res)
	return &res, nil
}

func (h *RelayerHandler) fetchAndSave(seq int) {
	logrus.Infof("RelayerHandler coroutine-%d start", seq)
	st := time.Now().Unix()
	for i, file := range h.jsonFileList {
		if i%concurrentNum != seq {
			continue
		}

		split := strings.Split(file, "|")
		pairJson, err := h.queryPairJson(split[0], split[1]) // split[0]: filepath, split[1]: json file name
		if err != nil {
			continue
		}

		chain1 := strings.ReplaceAll(pairJson.Chain1.ChainId, "-", "_")
		chain2 := strings.ReplaceAll(pairJson.Chain2.ChainId, "-", "_")
		cfgEntity := entity.GenerateRelayerConfigEntity(chain1, pairJson.Chain1.ChannelId, pairJson.Chain1.Address, chain2, pairJson.Chain2.ChannelId, pairJson.Chain2.Address)
		if _, ok := h.pairIdMap[cfgEntity.RelayerPairId]; !ok {
			cfgEntity.RelayerName = h.teamNameMap[split[0]]
			if h.imageMap[split[0]] {
				iconName := strings.ReplaceAll(cfgEntity.RelayerName, " ", "_")
				cfgEntity.Icon = fmt.Sprintf(iconUrl, iconName)
			}

			if err = relayerCfgRepo.Insert(cfgEntity); err != nil {
				logrus.Errorf("RelayerHandler insert relayer config error, %v", err)
			}
		}
	}
	logrus.Infof("RelayerHandler coroutine-%d end, time use: %d(s)", seq, time.Now().Unix()-st)
}
