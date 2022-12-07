package task

import (
	"encoding/json"
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
)

type IbcChainConfigTask struct {
	allChainList    []string  // all chain list
	channelStateMap *sync.Map // channel -> state map
	chainUpdateMap  *sync.Map // map[string]bool chain 最后是否能被更新map
	chainChannelMap *sync.Map // chain -> chain的所有channel map
	chainIdNameMap  map[string]string
}

var _ibcChainConfigTask Task = new(IbcChainConfigTask)

func (t *IbcChainConfigTask) Name() string {
	return "ibc_chain_config_task"
}
func (t *IbcChainConfigTask) Cron() int {
	if taskConf.CronTimeChainConfigTask > 0 {
		return taskConf.CronTimeChainConfigTask
	}
	return EveryMinute
}

func (t *IbcChainConfigTask) Run() int {
	t.init()
	chainConfList, err := t.getChainConf()
	if err != nil {
		logrus.Errorf("task %s getChainConf error, %s", t.Name(), err.Error())
		return -1
	}

	// 获取所有chain的channel信息
	var wg sync.WaitGroup
	wg.Add(len(chainConfList))
	for _, v := range chainConfList {
		chain := v
		go func() {
			defer wg.Done()
			channelPathList, err := t.getIbcChannels(chain.ChainName, chain.GrpcRestGateway, chain.LcdApiPath.ChannelsPath)
			if err != nil {
				t.chainUpdateMap.Store(chain.ChainName, false) // 出错时，此链的信息将不会被更新
			} else {
				t.setChainAndCounterpartyState(chain, channelPathList)
				t.chainUpdateMap.Store(chain.ChainName, true)
			}
			t.chainChannelMap.Store(chain.ChainName, channelPathList)
		}()
	}
	wg.Wait()

	// 为channel设置counterparty state
	for _, chain := range chainConfList {
		t.setCounterpartyState(chain.ChainName)
	}

	for _, chain := range chainConfList {
		enableUpdate, ok := t.chainUpdateMap.Load(chain.ChainName)
		if ok {
			if enableUpdate.(bool) {
				t.updateChain(chain)
			}
		}
	}

	return 1
}

func (t *IbcChainConfigTask) init() {
	t.channelStateMap = new(sync.Map)
	t.chainUpdateMap = new(sync.Map)
	t.chainChannelMap = new(sync.Map)
	mapData, err := repository.GetChainIdNameMap()
	if err != nil {
		logrus.Fatal(err.Error())
	}
	t.chainIdNameMap = mapData
}

func (t *IbcChainConfigTask) getChainConf() ([]*entity.ChainConfig, error) {
	chainConfList, err := chainConfigRepo.FindAll()
	if err != nil {
		return nil, err
	}

	allChainList := make([]string, 0, len(chainConfList))
	for _, v := range chainConfList {
		allChainList = append(allChainList, v.ChainName)
	}
	t.allChainList = allChainList

	return chainConfList, nil
}

// getIbcChannels 通过lcd channels_path 接口获取链上存在的所有channel信息
func (t *IbcChainConfigTask) getIbcChannels(chain, lcd, apiPath string) ([]*entity.ChannelPath, error) {
	if lcd == "" {
		logrus.Errorf("task %s %s getIbcChannels error, lcd error", t.Name(), chain)
		return nil, fmt.Errorf("lcd error")
	}

	limit := 1000
	offset := 0
	var channelPathList []*entity.ChannelPath

	for {
		apiPath = strings.ReplaceAll(apiPath, replaceHolderOffset, strconv.Itoa(offset))
		apiPath = strings.ReplaceAll(apiPath, replaceHolderLimit, strconv.Itoa(limit))
		url := fmt.Sprintf("%s%s", lcd, apiPath)
		bz, err := utils.HttpGet(url)
		if err != nil {
			logrus.Errorf("task %s %s getIbcChannels error, %v", t.Name(), chain, err)
			return nil, err
		}

		var resp vo.IbcChannelsResp
		err = json.Unmarshal(bz, &resp)
		if err != nil {
			logrus.Errorf("task %s %s getIbcChannels error, %v", t.Name(), chain, err)
			return nil, err
		}

		for _, v := range resp.Channels {
			channelPathList = append(channelPathList, &entity.ChannelPath{
				State:     v.State,
				PortId:    v.PortId,
				ChannelId: v.ChannelId,
				Chain:     "",
				ScChain:   chain,
				Counterparty: entity.CounterParty{
					State:     "",
					PortId:    v.Counterparty.PortId,
					ChannelId: v.Counterparty.ChannelId,
				},
			})
			k := fmt.Sprintf("%s%s%s%s%s", chain, v.PortId, v.ChannelId, v.Counterparty.PortId, v.Counterparty.ChannelId)
			t.channelStateMap.Store(k, v.State)
		}

		if len(resp.Channels) < limit {
			break
		}
		offset += limit
	}

	return channelPathList, nil
}

// setChainAndCounterpartyState 设置channel path的目标链chain 和 目标链channel state
// 1. 对于之前已经存在的channel，取之前的值即可;对于新增的channel，需要查询lcd 接口获取
// 2. 对于之前已经存在的channel，目标链channel state，暂取之前的值，后面 setCounterpartyState 方法会进一步处理
func (t *IbcChainConfigTask) setChainAndCounterpartyState(chain *entity.ChainConfig, channelPathList []*entity.ChannelPath) {
	existChannelStateMap := make(map[string]*entity.ChannelPath)
	for _, ibcInfo := range chain.IbcInfo {
		for _, path := range ibcInfo.Paths {
			key := fmt.Sprintf("%s%s%s%s", path.PortId, path.ChannelId, path.Counterparty.PortId, path.Counterparty.ChannelId)
			existChannelStateMap[key] = path
		}
	}

	lcdConnectionErr := false
	for _, v := range channelPathList {
		key := fmt.Sprintf("%s%s%s%s", v.PortId, v.ChannelId, v.Counterparty.PortId, v.Counterparty.ChannelId)
		existChannelState, ok := existChannelStateMap[key]
		if ok {
			v.Counterparty.State = existChannelState.Counterparty.State
		}

		if ok && existChannelState.Chain != "" && existChannelState.ClientId != "" {
			v.Chain = existChannelState.Chain
			v.ClientId = existChannelState.ClientId
		} else {
			if !lcdConnectionErr { // 如果遇到lcd连接问题，则不再请求lcd.
				stateResp, err := queryClientState(chain.GrpcRestGateway, chain.LcdApiPath.ClientStatePath, v.PortId, v.ChannelId)
				if err != nil {
					lcdConnectionErr = isConnectionErr(err)
					logrus.Errorf("task %s %s queryClientState error, %v", t.Name(), chain.ChainName, err)
				} else {
					v.Chain = t.chainIdNameMap[stateResp.IdentifiedClientState.ClientState.ChainId]
					v.ClientId = stateResp.IdentifiedClientState.ClientId
				}
			}
		}
	}
}

func (t *IbcChainConfigTask) setCounterpartyState(chain string) {
	channels, ok := t.chainChannelMap.Load(chain)
	if !ok {
		return
	}

	for _, v := range channels.([]*entity.ChannelPath) {
		key := fmt.Sprintf("%s%s%s%s%s", v.Chain, v.Counterparty.PortId, v.Counterparty.ChannelId, v.PortId, v.ChannelId)
		counterpartyState, ok := t.channelStateMap.Load(key)
		if ok {
			v.Counterparty.State = counterpartyState.(string)
		}
	}
}

func (t *IbcChainConfigTask) updateChain(chainConf *entity.ChainConfig) {
	channelGroupMap := make(map[string][]*entity.ChannelPath)
	channels, ok := t.chainChannelMap.Load(chainConf.ChainName)
	if !ok {
		return
	}

	for _, v := range channels.([]*entity.ChannelPath) {
		if !utils.InArray(t.allChainList, v.Chain) {
			continue
		}

		channelGroupMap[v.Chain] = append(channelGroupMap[v.Chain], v)
	}

	ibcInfoList := make([]*entity.IbcInfo, 0, len(channelGroupMap))
	for dcChain, paths := range channelGroupMap {
		sort.Slice(paths, func(i, j int) bool {
			return paths[i].ChannelId < paths[j].ChannelId
		})
		ibcInfoList = append(ibcInfoList, &entity.IbcInfo{
			Chain: dcChain,
			Paths: paths,
		})
	}

	sort.Slice(ibcInfoList, func(i, j int) bool {
		return ibcInfoList[i].Chain < ibcInfoList[i].Chain
	})

	hashCode := utils.Md5(utils.MustMarshalJsonToStr(ibcInfoList))
	if hashCode == chainConf.IbcInfoHashLcd {
		return
	}

	chainConf.IbcInfoHashLcd = hashCode
	chainConf.IbcInfo = ibcInfoList
	if err := chainConfigRepo.UpdateIbcInfo(chainConf); err != nil {
		logrus.Errorf("task %s %s UpdateIbcInfo error, %v", t.Name(), chainConf.ChainName, err)
	}
}
