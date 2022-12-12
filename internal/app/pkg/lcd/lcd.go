package lcd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
)

var lcdTxDataCacheRepo cache.LcdTxDataCacheRepo

const (
	replaceHolderAddress = "{address}"
)

func GetAccount(chain, address, lcd, apiPath string) (*vo.AccountResp, error) {
	if state, err := lcdTxDataCacheRepo.GetAccount(chain, address); err == nil {
		return state, nil
	}

	apiPath = strings.ReplaceAll(apiPath, replaceHolderAddress, address)
	url := fmt.Sprintf("%s%s", lcd, apiPath)

	bz, err := utils.HttpGet(url)
	if err != nil {
		return nil, err
	}

	var resp vo.AccountResp
	err = json.Unmarshal(bz, &resp)
	if err != nil {
		return nil, err
	}

	_ = lcdTxDataCacheRepo.SetAccount(chain, address, &resp)
	return &resp, nil
}
