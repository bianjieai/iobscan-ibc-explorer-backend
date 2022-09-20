package metrics

import (
	"fmt"
	"github.com/weichang-bianjie/metric-sdk"
	"github.com/weichang-bianjie/metric-sdk/types"
)

const defaultPort = "9090"

type Monitor interface {
	Report(reports ...func())
}

type client struct {
	metric_sdk.MetricClient
}

func NewMonitor(port string) Monitor {
	if port == "" {
		port = defaultPort
	}
	metricClient := metric_sdk.NewClient(types.Config{
		Address: fmt.Sprintf(":%v", port),
	})

	return client{metricClient}
}

func (c client) Report(reports ...func()) {
	c.MetricClient.Start(func() {
		for _, report := range reports {
			go report()
		}
	})
}
