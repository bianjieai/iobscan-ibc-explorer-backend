package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/qiniu/qmgo"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
	"time"
)

type IbcRelayerCronTask struct {
}

func init() {
	RegisterTasks(&IbcRelayerCronTask{})
}

func (t *IbcRelayerCronTask) Name() string {
	return "ibc_relayer_task"
}
func (t *IbcRelayerCronTask) Cron() string {
	return TwentyMinute
}

func (t *IbcRelayerCronTask) Run() {
	group := sync.WaitGroup{}
	group.Add(2)
	go func() {
		defer group.Done()
		t.handleNewRelayer()
	}()
	go func() {
		defer group.Done()
		t.CheckAndChangeStatus()
	}()
	group.Wait()
	t.UpdateIbcChainsRelayer()
}

func (t *IbcRelayerCronTask) handleNewRelayer() {
	relayer, err := relayerRepo.FindLatestOne()
	if err != nil && err != qmgo.ErrNoSuchDocuments {
		logrus.Errorf("findLatestone relayer fail, %s", err.Error())
		return
	}
	latestTxTime := int64(0)
	if relayer != nil {
		latestTxTime = relayer.LatestTxTime
	}
	currentLatestTxTime, _ := ibcTxRepo.GetLatestTxTime()
	relayersData := t.handleIbcTxLatest(latestTxTime)
	if len(relayersData) > 0 && currentLatestTxTime > latestTxTime {
		relayersData[len(relayersData)-1].LatestTxTime = currentLatestTxTime
	}
	relayersHistoryData := t.handleIbcTxHistory(latestTxTime)
	relayersData = append(relayersData, relayersHistoryData...)
	if len(relayersData) > 0 {
		if err := relayerRepo.Insert(relayersData); err != nil {
			logrus.Error("insert  relayer data fail, ", err.Error())
		}
	}
}

func (t *IbcRelayerCronTask) CheckAndChangeStatus() {
	skip := int64(0)
	limit := int64(50)
	for {
		relayers, err := relayerRepo.FindAll(skip, limit)
		if err != nil {
			logrus.Error("find relayer by page fail, ", err.Error())
			return
		}
		for _, val := range relayers {
			timePeriod, updateTime, err := t.getTimePeriodAndupdateTime(val.ChainA, val.ChainB)
			if err != nil {
				logrus.Error(err.Error())
				continue
			}
			if timePeriod == -1 {
				//todo call relayer api (uri为/chain/:id in https://hermes.informal.systems/rest-api.html)
			}
			//Running=>Close: update_client 时间大于relayer基准周期
			if val.TimePeriod > 0 && val.UpdateTime > 0 && val.TimePeriod < updateTime-val.UpdateTime {
				if val.Status == entity.RelayerRunning {
					val.Status = entity.RelayerStop
				}
			}
			paths := t.getChannelsStatus(val.ChainA, val.ChainB)
			status := entity.RelayerUnknow
			//Running=>Close: relayer中继通道只要出现状态不是STATE_OPEN
			if val.Status == entity.RelayerRunning {
				for _, path := range paths {
					if path.ChannelId == val.ChannelB {
						if path.State != constant.ChannelStateOpen {
							status = entity.RelayerStop
							break
						}
					}
					if path.Counterparty.ChannelId == val.ChannelA {
						if path.Counterparty.State != constant.ChannelStateOpen {
							status = entity.RelayerStop
							break
						}
					}
				}
			} else {
				// Close=>Running: relayer的双向通道状态均为STATE_OPEN且update_client 时间小于relayer基准周期
				if val.TimePeriod > updateTime-val.UpdateTime {
					var channelStatus []string
					for _, path := range paths {
						if path.ChannelId == val.ChannelB {
							channelStatus = append(channelStatus, path.State)
						}
						if path.Counterparty.ChannelId == val.ChannelA {
							channelStatus = append(channelStatus, path.Counterparty.State)
						}
					}
					if len(channelStatus) == 2 {
						if channelStatus[0] == channelStatus[1] && channelStatus[0] == constant.ChannelStateOpen {
							status = entity.RelayerRunning
						}
					}
				}
			}
			if err := relayerRepo.Update(val.RelayerId, bson.M{
				"$set": bson.M{
					"status":      status,
					"update_time": updateTime,
					"time_period": timePeriod,
					"update_at":   time.Now().Unix(),
				}}); err != nil {
				logrus.Error("update relayer about time_period and update_time fail, ", err.Error())
			}
		}
		if len(relayers) < int(limit) {
			break
		}
		skip += limit
	}
}

func (t *IbcRelayerCronTask) getChannelsStatus(chainId, dcChainId string) []*entity.ChannelPath {
	// use cache find channels
	var ibcPaths []*entity.ChannelPath
	if paths, _ := ibcInfoCache.Get(chainId, dcChainId); paths != nil {
		data := paths.(string)
		utils.UnmarshalJsonIgnoreErr([]byte(data), &ibcPaths)
	}
	return ibcPaths
}
func (t *IbcRelayerCronTask) CountTxsAndSuccessRate() {

}

func (t *IbcRelayerCronTask) UpdateIbcChainsRelayer() {
	res, err := chainRepo.FindAll()
	if err != nil {
		logrus.Error("find ibc_chains data fail, ", err.Error())
		return
	}

	for _, val := range res {
		relayerCnt, err := relayerRepo.FindRelayersCnt(val.ChainId)
		if err != nil {
			logrus.Error("count relayers of chain fail, ", err.Error())
			continue
		}
		if relayerCnt > 0 {
			if err := chainRepo.UpdateRelayers(val.ChainId, relayerCnt); err != nil {
				logrus.Error("update ibc_chain relayers fail, ", err.Error())
			}
		}
	}
	return
}

func (t *IbcRelayerCronTask) handleIbcTxLatest(latestTxTime int64) []entity.IBCRelayer {
	relayerDtos, err := ibcTxRepo.GetRelayerInfo(latestTxTime)
	if err != nil {
		logrus.Errorf("get relayer info fail, %s", err.Error())
		return nil
	}
	var relayers []entity.IBCRelayer
	for _, dto := range relayerDtos {
		ibcTx, err := ibcTxRepo.GetOneRelayerScTxPacketId(dto)
		if err != nil {
			logrus.Errorf("get ibcTxLatest relayer packetId fail, %s", err.Error())
			continue
		}
		scAddrs, err := txRepo.GetRelayerScChainAddr(ibcTx.ScTxInfo.Msg.Msg.PacketId, dto.ScChainId)
		if err != nil {
			logrus.Errorf("get ibcTxLatest relayer packetId fail, %s", err.Error())
			continue
		}
		relayers = append(relayers, t.creteRelayerData(dto, scAddrs))
	}
	return relayers
}

func (t *IbcRelayerCronTask) handleIbcTxHistory(latestTxTime int64) []entity.IBCRelayer {
	relayerDtos, err := ibcTxRepo.GetHistoryRelayerInfo(latestTxTime)
	if err != nil {
		logrus.Errorf("get relayer info fail, %s", err.Error())
		return nil
	}
	var relayers []entity.IBCRelayer
	for _, dto := range relayerDtos {
		ibcTx, err := ibcTxRepo.GetHistoryOneRelayerScTxPacketId(dto)
		if err != nil {
			logrus.Errorf("get ibcTxLatest relayer packetId fail, %s", err.Error())
			continue
		}
		scAddrs, err := txRepo.GetRelayerScChainAddr(ibcTx.ScTxInfo.Msg.Msg.PacketId, dto.ScChainId)
		if err != nil {
			logrus.Errorf("get ibcTxLatest relayer packetId fail, %s", err.Error())
			continue
		}
		relayers = append(relayers, t.creteRelayerData(dto, scAddrs))
	}
	return relayers
}

func (t *IbcRelayerCronTask) creteRelayerData(dto *dto.GetRelayerInfoDTO,
	scAddrs []*dto.GetRelayerScChainAddreeDTO) entity.IBCRelayer {
	chainAAddress := make([]string, 0, len(scAddrs))
	for _, val := range scAddrs {
		chainAAddress = append(chainAAddress, val.ScChainAddress)
	}
	chainAAddr := ""
	if len(chainAAddress) > 0 {
		chainAAddr = chainAAddress[0]
	}
	relayerId := utils.Md5(dto.ScChannel + dto.DcChannel + dto.ScChainId + dto.DcChainId + chainAAddr + dto.DcChainAddress)
	return entity.IBCRelayer{
		RelayerId:     relayerId,
		ChainA:        dto.ScChainId,
		ChainB:        dto.DcChainId,
		ChannelA:      dto.ScChannel,
		ChannelB:      dto.DcChannel,
		ChainAAddress: chainAAddress,
		ChainBAddress: dto.DcChainAddress,
		CreateAt:      time.Now().Unix(),
		UpdateAt:      time.Now().Unix(),
	}
}

func (t *IbcRelayerCronTask) ExpireTime() time.Duration {
	return 3*time.Minute - 2*time.Second
}

//1: timePeriod
//2: updateTime
//3: error
func (t *IbcRelayerCronTask) getTimePeriodAndupdateTime(chainA, chainB string) (int64, int64, error) {
	updateTimeA, timePeriodA, err := txRepo.GetTimePeriodByUpdateClient(chainA)
	if err != nil {
		return 0, 0, err
	}
	updateTimeB, timePeriodB, err := txRepo.GetTimePeriodByUpdateClient(chainB)
	if err != nil {
		return 0, 0, err
	}
	timePeriod := timePeriodB
	if timePeriodA >= timePeriodB {
		timePeriod = timePeriodA
	}
	updateTime := updateTimeB
	if updateTimeA >= updateTimeB {
		updateTime = updateTimeA
	}
	return timePeriod, updateTime, nil
}
