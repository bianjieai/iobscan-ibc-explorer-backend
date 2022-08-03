package main

import (
	"runtime"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/cmd"
)

// @title Iobscan Ibc Explorer Support API
// @version visit /version
// @description Iobscan Ibc Explorer Support API document
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 8)
	cmd.Execute()
}
