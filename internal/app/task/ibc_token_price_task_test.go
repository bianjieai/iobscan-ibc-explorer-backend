package task

import (
	"testing"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
)

var tokenPriceTask TokenPriceTask

func TestTokenPriceTaskRun(t *testing.T) {
	global.Config = &conf.Config{Spi: conf.Spi{CoingeckoPriceUrl: "https://api.coingecko.com/api/v3/simple/price"}}
	tokenPriceTask.Run()
}
