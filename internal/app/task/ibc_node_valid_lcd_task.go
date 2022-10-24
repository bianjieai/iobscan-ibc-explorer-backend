package task

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
)

type IbcNodeLcdCronTask struct {
}

func (t *IbcNodeLcdCronTask) Name() string {
	return "ibc_node_valid_lcd_task"
}
func (t *IbcNodeLcdCronTask) Cron() int {
	return OneDay
}
func (t *IbcNodeLcdCronTask) Run() int {
	chainCfgs, err := chainConfigRepo.FindAllChainInfos()
	if err != nil {
		logrus.Errorf("task %s run error, %s", t.Name(), err.Error())
		return -1
	}

	t.doHandleChains(3, chainCfgs)
	return 1

}

func (t *IbcNodeLcdCronTask) RunWithParam(chainId string) int {
	if chainId != "" {
		CheckAndUpdateTraceSourceNode(chainId)
		return 1
	}

	return t.Run()
}

func (t *IbcNodeLcdCronTask) ExpireTime() time.Duration {
	return 24*time.Hour - 1*time.Second
}

func (t *IbcNodeLcdCronTask) doHandleChains(workNum int, chaincfgs []*entity.ChainConfig) {
	if workNum <= 0 {
		return
	}
	st := time.Now().Unix()
	logrus.Infof("task %s worker group start", t.Name())
	defer func() {
		logrus.Infof("task %s worker group end, time use: %d(s)", t.Name(), time.Now().Unix()-st)
	}()
	if len(chaincfgs) == 0 {
		return
	}
	var wg sync.WaitGroup
	wg.Add(workNum)
	for i := 0; i < workNum; i++ {
		num := i
		go func(num int) {
			defer wg.Done()

			for id, v := range chaincfgs {
				if id%workNum != num {
					continue
				}
				logrus.Infof("task %s worker %d chain-id: %s", t.Name(), num, v.ChainId)
				CheckAndUpdateTraceSourceNode(v.ChainId)
			}
		}(num)
	}
	wg.Wait()
}

func getChainRegisterResp(chainId string) (vo.ChainRegisterResp, error) {
	var chainRegisterResp vo.ChainRegisterResp
	chainRegistry, err := chainRegistryRepo.FindOne(chainId)
	if err != nil {
		return chainRegisterResp, fmt.Errorf("find chain_registry by chain-id(%s) error: %s", chainId, err.Error())
	}

	bz, err := utils.HttpGet(chainRegistry.ChainJsonUrl)
	if err != nil {
		return chainRegisterResp, fmt.Errorf("get chain registry json error: %s", err.Error())
	}

	_ = json.Unmarshal(bz, &chainRegisterResp)
	return chainRegisterResp, nil
}

func CheckAndUpdateTraceSourceNode(chainId string) {
	chainRegisterResp, err := getChainRegisterResp(chainId)
	if err != nil {
		logrus.Error(err)
		return
	}
	rpcAddrMap := make(map[string]cache.TraceSourceLcd, len(chainRegisterResp.Apis.Rpc))
	for _, rest := range chainRegisterResp.Apis.Rest {
		var rpcAddress string
		for _, rpc := range chainRegisterResp.Apis.Rpc {
			if rest.Provider == rpc.Provider {
				rpcAddress = rpc.Address
				break
			}
		}

		if rpcAddress == "" { // 没有对应的rpc
			rpcAddrMap[rest.Provider] = cache.TraceSourceLcd{
				FullNode:      false,
				TxIndexEnable: true,
				LcdAddr:       rest.Address,
			}
		} else {
			ok, earliestH := checkNodeTxIndex(rpcAddress)
			//node no reach
			if earliestH < 0 {
				continue
			}
			rpcAddrMap[rest.Provider] = cache.TraceSourceLcd{
				FullNode:      earliestH == 1,
				TxIndexEnable: ok,
				LcdAddr:       rest.Address,
			}
		}
	}

	if len(rpcAddrMap) == 0 {
		logrus.Warnf("CheckAndUpdateTraceSourceNode chain %s addr map is empty", chainId)
		return
	}

	res := make([]cache.TraceSourceLcd, 0, len(rpcAddrMap))
	var needSort bool
	for _, val := range rpcAddrMap {
		//出现全节点且支持交易查询
		if val.FullNode && val.TxIndexEnable {
			needSort = true
		}
		res = append(res, val)
	}
	if needSort {
		//将可用的全节点放在第一个
		for i := range res {
			if res[i].FullNode && res[i].TxIndexEnable {
				res[i], res[0] = res[0], res[i]
			}
		}
	}

	var lcdAddrCache cache.LcdAddrCacheRepo
	chainRegisterResp.ChainId = strings.ReplaceAll(chainRegisterResp.ChainId, "-", "_")
	_ = lcdAddrCache.Set(chainRegisterResp.ChainId, res)

}

// checkNodeTxIndex
//If tx_index is 'on', return true,earliest_height. Else return false,0
//If node is no reach, return false,-1
func checkNodeTxIndex(rpc string) (bool, int64) {
	bz, err := utils.HttpGet(fmt.Sprintf("%s/status", rpc))
	if err != nil {
		logrus.Errorf("checkNodeTxIndex rpc status api error: %v", err)
		return false, -1
	}

	var statusResp vo.StatusResp
	_ = json.Unmarshal(bz, &statusResp)
	if strings.Compare(strings.ToLower(statusResp.Result.NodeInfo.Other.TxIndex), "off") == 0 {
		return false, statusResp.Result.SyncInfo.EarliestBlockHeight
	}

	return true, statusResp.Result.SyncInfo.EarliestBlockHeight
}
