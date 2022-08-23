package cmd

import (
	"io/ioutil"
	"os"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/spf13/cobra"
)

var (
	ConfigFilePath string
	startCmd       = &cobra.Command{
		Use:   "start",
		Short: "Start Iobscan IBC Explorer Backend.",
		Run: func(cmd *cobra.Command, args []string) {
			online()
		},
	}
	testCmd = &cobra.Command{ // test
		Use:   "test",
		Short: "Start Test Iobscan IBC Explorer Backend.",
		Run: func(cmd *cobra.Command, args []string) {
			test()
		},
	}
)

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.AddCommand(testCmd)
	testCmd.Flags().StringVarP(&ConfigFilePath, "CONFIG", "c", "", "conf path: /opt/local.toml")

}

func test() {
	data, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		panic(err)
	}
	config, err := conf.ReadConfig(data)
	if err != nil {
		panic(err)
	}
	run(config)
}

func online() {
	var config *conf.Config

	filepath, found := os.LookupEnv(constant.EnvNameConfigFilePath)
	if found {
		ConfigFilePath = filepath
	} else {
		panic("not found CONFIG_FILE_PATH")
	}
	data, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		panic(err)
	}
	config, err = conf.ReadConfig(data)
	if err != nil {
		panic(err)
	}

	run(config)
}

func run(cfg *conf.Config) {
	app.Serve(cfg)
}
