package task

import (
	"context"
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/qiniu/qmgo"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type IModifyChainIdContent interface {
	Run() int
}

func NewModifyChainIdContent(collName string) IModifyChainIdContent {
	chainCfgMap, err := _initChainCfgMap()
	if err != nil {
		return nil
	}

	switch collName {
	case entity.ChainConfig{}.CollectionName():
		return NewModifyChainConfig()
	case entity.IBCRelayerNew{}.CollectionName():
		return NewModifyIbcRelayer(chainCfgMap)
	case entity.IbcTaskRecord{}.CollectionName():
		return NewModifyIbcTaskRecord(chainCfgMap)
	case entity.IBCChannelStatisticsCollName, entity.IBCChannel{}.CollectionName():
		return NewModifyChannalIdTask(chainCfgMap)
	}
	return nil
}

//========================================================================================
//========================================================================================
//========get all chain config data================================
func _initChainCfgMap() (map[string]string, error) {
	mapData, err := chainConfigRepo.GetChainCfgMap()
	if err != nil {
		return nil, err
	}
	return mapData, nil
}

//========chain_id format convert================================
func _formatChainId(chainId string) string {
	return strings.ReplaceAll(chainId, "_", "-")
}

//========================================================================================
//========================================================================================
//============================modify ibc_task_record taskname===========================
type fixIbcTaskRecordTask struct {
	chainCfgMap map[string]string
}

func NewModifyIbcTaskRecord(chainCfgMapData map[string]string) *fixIbcTaskRecordTask {
	return &fixIbcTaskRecordTask{
		chainCfgMap: chainCfgMapData,
	}
}

type _ibcTaskRecord struct {
	Id       primitive.ObjectID `bson:"_id"`
	TaskName string             `bson:"task_name"`
}

func (t *fixIbcTaskRecordTask) Name() string {
	return "fix_ibc_task_record_task"
}

func (t *fixIbcTaskRecordTask) Run() int {
	datas, err := t.GetIbcTaskRecordData()
	if err != nil {
		logrus.Errorf("task[%s] get ibc_task_record faild,err is %s", t.Name(), err.Error())
		return -1
	}

	for _, val := range datas {
		err := t.UpdateIbcTaskRecord(*val)
		if err != nil {
			logrus.Errorf("task[%s] update ibc_task_record [%s] faild,err is %s", t.Name(), val.TaskName, err.Error())
		}
	}
	return 1
}

func (t *fixIbcTaskRecordTask) _historyIbcTaskRecordRepo() *qmgo.Collection {
	return repository.GetDatabase().Collection(entity.IbcTaskRecord{}.CollectionName())
}

func (t *fixIbcTaskRecordTask) GetIbcTaskRecordData() ([]*_ibcTaskRecord, error) {
	var datas []*_ibcTaskRecord
	err := t._historyIbcTaskRecordRepo().Find(context.Background(), bson.M{}).All(&datas)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (t *fixIbcTaskRecordTask) UpdateIbcTaskRecord(record _ibcTaskRecord) error {
	if len(t.chainCfgMap) == 0 {
		return fmt.Errorf("init don't work")
	}
	arrs := strings.Split(record.TaskName, "_")
	chainId := strings.Join(arrs[1:len(arrs)-1], "_")

	chain, ok := t.chainCfgMap[_formatChainId(chainId)]
	if !ok {
		return fmt.Errorf("chain-id[%s] no found in chain_config", _formatChainId(chainId))
	}
	return t._historyIbcTaskRecordRepo().UpdateId(context.Background(), record.Id, bson.M{
		"$set": bson.M{
			"task_name": "sync_" + chain + "_transfer",
		},
	})
}

//========================================================================================
//============================modify ibc_relayer chain_id -> chain========================
//============================modify ibc_relayer channel_pair_id =========================

type fixIbcRelayerTask struct {
	chainCfgMap map[string]string
}

func NewModifyIbcRelayer(chainCfgMapData map[string]string) *fixIbcRelayerTask {
	return &fixIbcRelayerTask{
		chainCfgMap: chainCfgMapData,
	}
}

func (t *fixIbcRelayerTask) Name() string {
	return "fix_ibc_relayer_task"
}

func (t *fixIbcRelayerTask) Run() int {
	skip := int64(0)
	limit := int64(1000)
	for {
		datas, err := t.GetIbcRelayerData(skip, limit)
		if err != nil {
			logrus.Errorf("task[%s] get ibc_relayer faild,err is %s", t.Name(), err.Error())
			return -1
		}

		for _, val := range datas {
			err := t.UpdateIbcRelayerData(*val)
			if err != nil {
				logrus.Errorf("task[%s] update ibc_relayer [relayer_id:%s] faild,err is %s", t.Name(), val.RelayerId, err.Error())
			}
		}

		if len(datas) < int(limit) {
			break
		}
		skip += limit
	}

	return 1
}

func (t *fixIbcRelayerTask) GetIbcRelayerData(skip, limit int64) ([]*entity.IBCRelayerNew, error) {
	datas, err := relayerRepo.FindAll(skip, limit, repository.RelayerAllType)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (t *fixIbcRelayerTask) UpdateIbcRelayerData(relayer entity.IBCRelayerNew) error {
	if len(t.chainCfgMap) == 0 {
		return fmt.Errorf("init don't work")
	}
	channelPairInfos := make([]entity.ChannelPairInfo, 0, len(relayer.ChannelPairInfo))
	for _, val := range relayer.ChannelPairInfo {
		chainA, ok := t.chainCfgMap[_formatChainId(val.ChainA)]
		if !ok {
			return fmt.Errorf("chainA[%s] no found in chain_config", _formatChainId(val.ChainA))
		}
		chainB, ok := t.chainCfgMap[_formatChainId(val.ChainB)]
		if !ok {
			return fmt.Errorf("chainB[%s] no found in chain_config", _formatChainId(val.ChainB))
		}
		val.PairId = entity.GenerateRelayerPairId(chainA, val.ChannelA, val.ChainAAddress,
			chainB, val.ChannelB, val.ChainBAddress)
		val.ChainA = chainA
		val.ChainB = chainB
		channelPairInfos = append(channelPairInfos, val)
	}

	return relayerRepo.UpdateChannelPairInfo(relayer.RelayerId, channelPairInfos)
}

//========================================================================================
//========================================================================================
//================modify ibc_channel base_denom_chain_id -> base_denom_chain==
//================modify ibc_channel channel_id ===============================
//==========modify ibc_channel_statistic channel_id==============================

type fixChannelIdTask struct {
	chainCfgMap map[string]string
}

func NewModifyChannalIdTask(chainCfgMapData map[string]string) *fixChannelIdTask {
	return &fixChannelIdTask{
		chainCfgMap: chainCfgMapData,
	}
}

type _historyIBCChannelStatistics struct {
	Id               primitive.ObjectID `bson:"_id"`
	ChannelId        string             `bson:"channel_id"`
	BaseDenomChainId string             `bson:"base_denom_chain_id"`
}

type _historyIBCChannel struct {
	Id        primitive.ObjectID `bson:"_id"`
	ChannelId string             `bson:"channel_id"`
	ChainA    string             `bson:"chain_a"`
	ChainB    string             `bson:"chain_b"`
	ChannelA  string             `bson:"channel_a"`
	ChannelB  string             `bson:"channel_b"`
}

func (t *fixChannelIdTask) Name() string {
	return "fix_channel_id_task"
}

func (t *fixChannelIdTask) Run() int {

	handleIbcChannel := func() int {
		skip := int64(0)
		limit := int64(1000)
		for {
			datas, err := t.GetIbcChannelData(skip, limit)
			if err != nil {
				logrus.Errorf("task[%s] get ibc_channel faild,err is %s", t.Name(), err.Error())
				return -1
			}
			for _, val := range datas {
				if err := t.UpdateIbcChannel(*val); err != nil {
					logrus.Errorf("task[%s] update ibc_channel [channel_id:%s] faild,err is %s", t.Name(), val.ChannelId, err.Error())
				}
			}

			if len(datas) < int(limit) {
				break
			}
			skip += limit
		}

		return 1
	}

	handleIbcChannelStatistic := func() int {

		skip := int64(0)
		limit := int64(1000)
		for {
			statisticDatas, err := t.GetIbcChannelStatisticData(skip, limit)
			if err != nil {
				return -1
			}
			for _, val := range statisticDatas {
				if err := t.UpdateIbcChannelStatistic(*val); err != nil {
					logrus.Errorf("task[%s] update ibc_channel_statistic [channel_id:%s] faild,err is %s", t.Name(), val.ChannelId, err.Error())
				}
			}

			if len(statisticDatas) < int(limit) {
				break
			}
			skip += limit
		}
		return 1
	}

	if ret := handleIbcChannel(); ret < 0 {
		return ret
	}
	if ret := handleIbcChannelStatistic(); ret < 0 {
		return ret
	}

	return 1
}

func (t *fixChannelIdTask) _historyIbcChannelRepo() *qmgo.Collection {
	return repository.GetDatabase().Collection(entity.IBCChannel{}.CollectionName())
}

func (t *fixChannelIdTask) _historyIbcChannelStatisticRepo() *qmgo.Collection {
	return repository.GetDatabase().Collection(entity.IBCChannelStatisticsCollName)
}

func (t *fixChannelIdTask) GetIbcChannelData(skip, limit int64) ([]*_historyIBCChannel, error) {
	var datas []*_historyIBCChannel
	err := t._historyIbcChannelRepo().Find(context.Background(), bson.M{}).Skip(skip).Limit(limit).All(&datas)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (t *fixChannelIdTask) GetIbcChannelStatisticData(skip, limit int64) ([]*_historyIBCChannelStatistics, error) {
	var datas []*_historyIBCChannelStatistics
	err := t._historyIbcChannelStatisticRepo().Find(context.Background(), bson.M{}).Skip(skip).Limit(limit).All(&datas)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (t *fixChannelIdTask) UpdateIbcChannel(channel _historyIBCChannel) error {
	if len(t.chainCfgMap) == 0 {
		return fmt.Errorf("init don't work")
	}
	scChain, ok := t.chainCfgMap[_formatChainId(channel.ChainA)]
	if !ok {
		return fmt.Errorf("ChainA-id[%s] no found in chain_config", _formatChainId(channel.ChainA))
	}
	dcChain, ok := t.chainCfgMap[_formatChainId(channel.ChainB)]
	if !ok {
		return fmt.Errorf("ChainB-id[%s] no found in chain_config", _formatChainId(channel.ChainB))
	}
	channelId := generateChannelId(scChain, channel.ChannelA, dcChain, channel.ChannelB)
	return t._historyIbcChannelRepo().UpdateId(context.Background(), channel.Id,
		bson.M{
			"$set": bson.M{
				"channel_id": channelId,
				"chain_a":    scChain,
				"chain_b":    dcChain,
			},
		})
}

func (t *fixChannelIdTask) UpdateIbcChannelStatistic(channelStatistic _historyIBCChannelStatistics) error {
	if len(t.chainCfgMap) == 0 {
		return fmt.Errorf("init don't work")
	}
	chainA, channelA, chainB, channelB, err := t.parseChannelId(channelStatistic.ChannelId)
	if err != nil {
		return err
	}

	scChain, ok := t.chainCfgMap[_formatChainId(chainA)]
	if !ok {
		logrus.Warnf("update ibc_channel_statistic fail for [%s] no found in chain_config", _formatChainId(chainA))
		return nil
	}
	dcChain, ok := t.chainCfgMap[_formatChainId(chainB)]
	if !ok {
		logrus.Warnf("update ibc_channel_statistic fail for [%s] no found in chain_config", _formatChainId(chainB))
		return nil
	}
	baseDenomChain := t.chainCfgMap[_formatChainId(channelStatistic.BaseDenomChainId)]
	channelId := generateChannelId(scChain, channelA, dcChain, channelB)
	return t._historyIbcChannelStatisticRepo().UpdateId(context.Background(), channelStatistic.Id,
		bson.M{
			"$set": bson.M{
				"channel_id":       channelId,
				"base_denom_chain": baseDenomChain,
			},
			"$unset": bson.M{
				"base_denom_chain_id": 1,
			},
		})

}

func (t *fixChannelIdTask) parseChannelId(channelId string) (chainA, channelA, chainB, channelB string, err error) {
	split := strings.Split(channelId, "|")
	if len(split) != 4 {
		logrus.Errorf("task %s parseChannelId error, %v", t.Name(), err)
		return "", "", "", "", fmt.Errorf("channel id format error")
	}
	return split[0], split[1], split[2], split[3], nil
}

//========================================================================================
//========================================================================================
//==================modify chain_config chain_id -> current_chain_id==============
//==================modify chain_config chain_id -> chain=========================
type modifyChainConfigTask struct {
	chainCfgMap map[string]string
}

func NewModifyChainConfig() *modifyChainConfigTask {
	return &modifyChainConfigTask{}
}

type (
	_chainConfig struct {
		Id         primitive.ObjectID `bson:"_id"`
		ChainId    string             `bson:"chain_id"`
		ChainName  string             `bson:"chain_name"`
		LcdApiPath entity.ApiPath     `bson:"lcd_api_path"`
		Lcd        string             `bson:"lcd"`
		IbcInfo    []*_ibcInfo        `bson:"ibc_info"`
	}
	_ibcInfo struct {
		ChainId string          `bson:"chain_id"`
		Paths   []*_channelPath `bson:"paths"`
	}
	_channelPath struct {
		State        string              `bson:"state"`
		PortId       string              `bson:"port_id"`
		ChannelId    string              `bson:"channel_id"`
		ChainId      string              `bson:"chain_id"`
		ScChainId    string              `bson:"sc_chain_id"`
		ClientId     string              `bson:"client_id"`
		Counterparty entity.CounterParty `bson:"counterparty"`
	}
	CounterParty struct {
		State     string `bson:"state"`
		PortId    string `bson:"port_id"`
		ChannelId string `bson:"channel_id"`
	}
)

func (t *modifyChainConfigTask) Name() string {
	return "fix_chain_config_task"
}

func (t *modifyChainConfigTask) Run() int {

	datas, err := t.GetAllChainConigs()
	if err != nil {
		logrus.Errorf("task[%s] get chain_config faild,err is %s", t.Name(), err.Error())
		return -1
	}
	for _, val := range datas {
		if err := t.UpdateChainConfig(*val); err != nil {
			logrus.Errorf("task[%s] update chain_config [chain_id:%s] faild,err is %s", t.Name(), val.ChainId, err.Error())
		}
	}
	return 1
}

func (t *modifyChainConfigTask) _historyChainCfgRepo() *qmgo.Collection {
	return repository.GetDatabase().Collection(entity.ChainConfig{}.CollectionName())
}

func (t *modifyChainConfigTask) GetAllChainConigs() ([]*_chainConfig, error) {
	var datas []*_chainConfig
	err := t._historyChainCfgRepo().Find(context.Background(), bson.M{"chain_id": bson.M{"$exists": true}}).All(&datas)
	if err != nil {
		return nil, err
	}

	chainCfgMap := make(map[string]string, 10)
	for _, val := range datas {
		chainCfgMap[_formatChainId(val.ChainId)] = val.ChainName
	}
	t.chainCfgMap = chainCfgMap

	return datas, nil
}

func (t *modifyChainConfigTask) UpdateChainConfig(config _chainConfig) error {
	if len(t.chainCfgMap) == 0 {
		return fmt.Errorf("init don't work")
	}
	loadChannelPath := func(path *_channelPath) *entity.ChannelPath {
		return &entity.ChannelPath{
			State:        path.State,
			PortId:       path.PortId,
			ChannelId:    path.ChannelId,
			ClientId:     path.ClientId,
			Counterparty: path.Counterparty,
		}
	}

	ibcInfos := make([]*entity.IbcInfo, 0, len(config.IbcInfo))
	for _, val := range config.IbcInfo {
		paths := make([]*entity.ChannelPath, 0, len(val.Paths))
		for _, path := range val.Paths {
			item := loadChannelPath(path)
			chain, ok := t.chainCfgMap[_formatChainId(path.ChainId)]
			if !ok {
				return fmt.Errorf("Chain[%s] no found in chain_config", _formatChainId(path.ChainId))
			}
			scChain, ok := t.chainCfgMap[_formatChainId(path.ScChainId)]
			if !ok {
				return fmt.Errorf("ScChain[%s] no found in chain_config", _formatChainId(path.ScChainId))
			}
			item.ScChain = scChain
			item.Chain = chain
			paths = append(paths, item)
		}

		chain, ok := t.chainCfgMap[_formatChainId(val.ChainId)]
		if !ok {
			return fmt.Errorf("Chain[%s] no found in chain_config", _formatChainId(val.ChainId))
		}
		ibcInfo := &entity.IbcInfo{
			Chain: chain,
			Paths: paths,
		}
		ibcInfos = append(ibcInfos, ibcInfo)
	}

	return t._historyChainCfgRepo().UpdateId(context.Background(), config.Id,
		bson.M{
			"$set": bson.M{
				repository.ChainConfigFieldCurrentChainId:  _formatChainId(config.ChainId),
				repository.ChainConfigFieldGrpcRestGateway: config.Lcd,
				repository.ChainConfigFieldIbcInfo:         ibcInfos,
			},
			"$unset": bson.M{
				"chain_id": 1,
				"lcd":      1,
			},
		})

}
