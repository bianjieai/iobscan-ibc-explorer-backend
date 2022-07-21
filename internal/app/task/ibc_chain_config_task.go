package task

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
)

type IbcChainConfigTask struct {
	allChainList    []string                         // all chain id list
	channelStateMap map[string]string                // channel -> state map
	chainUpdateMap  map[string]bool                  // chain 最后是否能被更新map
	chainChannelMap map[string][]*entity.ChannelPath // chain -> chain的所有channel map
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
	for _, chain := range chainConfList {
		channelPathList, err := t.getIbcChannels(chain.ChainId, chain.Lcd, chain.LcdApiPath.ChannelsPath)
		if err != nil {
			t.chainUpdateMap[chain.ChainId] = false // 出错时，此链的信息将不会被更新
		} else {
			t.setChainIdAndCounterpartyState(chain, channelPathList)
			t.chainUpdateMap[chain.ChainId] = true
		}

		t.chainChannelMap[chain.ChainId] = channelPathList
	}

	// 为channel设置counterparty state
	for _, chain := range chainConfList {
		t.setCounterpartyState(chain.ChainId)
	}

	for _, chain := range chainConfList {
		if t.chainUpdateMap[chain.ChainId] {
			_ = t.updateChain(chain)
		}
	}

	return 1
}

func (t *IbcChainConfigTask) init() {
	t.channelStateMap = make(map[string]string)
	t.chainUpdateMap = make(map[string]bool)
	t.chainChannelMap = make(map[string][]*entity.ChannelPath)
}

func (t *IbcChainConfigTask) getChainConf() ([]*entity.ChainConfig, error) {
	chainConfList, err := chainConfigRepo.FindAll()
	if err != nil {
		return nil, err
	}

	allChainList := make([]string, 0, len(chainConfList))
	for _, v := range chainConfList {
		allChainList = append(allChainList, v.ChainId)
	}
	t.allChainList = allChainList

	return chainConfList, nil
}

// getIbcChannels 通过lcd channels_path 接口获取链上存在的所有channel信息
func (t *IbcChainConfigTask) getIbcChannels(chainId, lcd, apiPath string) ([]*entity.ChannelPath, error) {
	limit := 1000
	offset := 0
	var channelPathList []*entity.ChannelPath

	for {
		apiPath = strings.ReplaceAll(apiPath, replaceHolderOffset, strconv.Itoa(offset))
		apiPath = strings.ReplaceAll(apiPath, replaceHolderLimit, strconv.Itoa(limit))
		url := fmt.Sprintf("%s%s", lcd, apiPath)
		bz, err := utils.HttpGet(url)
		if err != nil {
			logrus.Errorf("task %s %s getIbcChannels error, %v", t.Name(), chainId, err)
			return nil, err
		}

		var resp vo.IbcChannelsResp
		err = json.Unmarshal(bz, &resp)
		if err != nil {
			logrus.Errorf("task %s %s getIbcChannels error, %v", t.Name(), chainId, err)
			return nil, err
		}

		for _, v := range resp.Channels {
			channelPathList = append(channelPathList, &entity.ChannelPath{
				State:     v.State,
				PortId:    v.PortId,
				ChannelId: v.ChannelId,
				ChainId:   "",
				ScChainId: chainId,
				Counterparty: entity.CounterParty{
					State:     "",
					PortId:    v.Counterparty.PortId,
					ChannelId: v.Counterparty.ChannelId,
				},
			})
			t.channelStateMap[fmt.Sprintf("%s%s%s%s%s", chainId, v.PortId, v.ChannelId, v.Counterparty.PortId, v.Counterparty.ChannelId)] = v.State
		}

		if len(resp.Channels) < limit {
			break
		}
		offset += limit
	}

	return channelPathList, nil
}

// setChainIdAndCounterpartyState 设置channel path的目标链chain id 和 目标链channel state
// 1. 对于之前已经存在的channel，取之前的值即可;对于新增的channel，需要查询lcd 接口获取
// 2. 对于之前已经存在的channel，目标链channel state，暂取之前的值，后面 setCounterpartyState 方法会进一步处理
func (t *IbcChainConfigTask) setChainIdAndCounterpartyState(chain *entity.ChainConfig, channelPathList []*entity.ChannelPath) {
	existDcChainIdMap := make(map[string]string)
	existDcChannelStateMap := make(map[string]string)
	for _, ibcInfo := range chain.IbcInfo {
		for _, path := range ibcInfo.Paths {
			key := fmt.Sprintf("%s%s%s%s", path.PortId, path.ChannelId, path.Counterparty.PortId, path.Counterparty.ChannelId)
			existDcChainIdMap[key] = path.ChainId
			existDcChannelStateMap[key] = path.Counterparty.State
		}
	}

	lcdConnectionErr := false
	for _, v := range channelPathList {
		key := fmt.Sprintf("%s%s%s%s", v.PortId, v.ChannelId, v.Counterparty.PortId, v.Counterparty.ChannelId)
		dcChainId, ok := existDcChainIdMap[key]
		if ok {
			v.ChainId = dcChainId
		} else {
			if !lcdConnectionErr { // 如果遇到lcd连接问题，则不再请求lcd.
				stateResp, err := t.queryClientState(chain.ChainId, chain.Lcd, chain.LcdApiPath.ClientStatePath, v.PortId, v.ChannelId)
				if err != nil {
					lcdConnectionErr = isConnectionErr(err)
				} else {
					v.ChainId = t.formatChainId(stateResp.IdentifiedClientState.ClientState.ChainId)
				}
			}
		}

		state, ok := existDcChannelStateMap[key]
		if ok {
			v.Counterparty.State = state
		}
	}
}

func (t *IbcChainConfigTask) formatChainId(chainId string) string {
	return strings.ReplaceAll(chainId, "-", "_")
}

// queryClientState 查询lcd client_state_path接口
func (t *IbcChainConfigTask) queryClientState(chainId, lcd, apiPath, port, channel string) (*vo.ClientStateResp, error) {
	apiPath = strings.ReplaceAll(apiPath, replaceHolderChannel, channel)
	apiPath = strings.ReplaceAll(apiPath, replaceHolderPort, port)
	url := fmt.Sprintf("%s%s", lcd, apiPath)
	bz, err := utils.HttpGet(url)
	if err != nil {
		logrus.Errorf("task %s %s queryClientState error, %s, %v", t.Name(), chainId, url, err)
		return nil, err
	}

	var resp vo.ClientStateResp
	err = json.Unmarshal(bz, &resp)
	if err != nil {
		logrus.Errorf("task %s %s queryClientState error, %s, %v", t.Name(), chainId, url, err)
		return nil, err
	}

	return &resp, nil
}

func (t *IbcChainConfigTask) setCounterpartyState(chainId string) {
	for _, v := range t.chainChannelMap[chainId] {
		key := fmt.Sprintf("%s%s%s%s%s", v.ChainId, v.Counterparty.PortId, v.Counterparty.ChannelId, v.PortId, v.ChannelId)
		counterpartyState, ok := t.channelStateMap[key]
		if ok {
			v.Counterparty.State = counterpartyState
		}
	}
}

func (t *IbcChainConfigTask) updateChain(chainConf *entity.ChainConfig) error {
	channelGroupMap := make(map[string][]*entity.ChannelPath)
	for _, v := range t.chainChannelMap[chainConf.ChainId] {
		if !utils.InArray(t.allChainList, v.ChainId) {
			continue
		}

		channelGroupMap[v.ChainId] = append(channelGroupMap[v.ChainId], v)
	}

	ibcInfoList := make([]*entity.IbcInfo, 0, len(channelGroupMap))
	for dcChainid, paths := range channelGroupMap {
		sort.Slice(paths, func(i, j int) bool {
			return paths[i].ChannelId < paths[j].ChannelId
		})
		ibcInfoList = append(ibcInfoList, &entity.IbcInfo{
			ChainId: dcChainid,
			Paths:   paths,
		})
	}

	sort.Slice(ibcInfoList, func(i, j int) bool {
		return ibcInfoList[i].ChainId < ibcInfoList[i].ChainId
	})

	hashCode := utils.Md5(utils.MustMarshalJsonToStr(ibcInfoList))
	if hashCode == chainConf.IbcInfoHashLcd {
		return nil
	}

	chainConf.IbcInfoHashLcd = hashCode
	chainConf.IbcInfo = ibcInfoList
	err := chainConfigRepo.UpdateIbcInfo(chainConf)
	if err != nil {
		logrus.Errorf("task %s %s UpdateIbcInfo error, %v", t.Name(), chainConf.ChainId, err)
	}
	return err
}
