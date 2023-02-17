package integration

import (
	"fmt"
	"time"
)

func (s IntegrationTestSuite) TestChainFeeStatistics() {
	for {
		res, err := feeService.ChainFeeStatistics("cosmoshub_4", 1644768000, 1676390399)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(res)

		time.Sleep(10 * time.Minute)
	}
}
