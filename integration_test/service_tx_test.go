package integration

import (
	"fmt"
	"time"
)

func (s IntegrationTestSuite) TestFlowInfoStatistics() {
	for {
		res, err := txService.FlowInfoStatistics("cosmoshub", 0, 99999999999999)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(res)

		time.Sleep(10 * time.Minute)
	}
}
