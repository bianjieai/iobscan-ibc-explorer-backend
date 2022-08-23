package main

import (
	"runtime"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/cmd"
)

// @title VisualizationServer Swagger API
// @version visit /version
// @description visualization-server api document
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 8)
	cmd.Execute()
}
