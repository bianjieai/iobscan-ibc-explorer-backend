package conf

import (
	"bytes"
	"github.com/spf13/viper"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/enum"
)

type Config struct {
	App App

	Mongo Mongo `mapstructure:"mongo_sup_tech"`
	Redis Redis
}

type Mysql struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Charset  string
	TimeZone string `mapstructure:"time_zone"`
}

type Mongo struct {
	Url      string
	Database string
}

type App struct {
	ServicePath     string
	ServerName      string
	Addr            string
	Env             string
	LogLevel        string `mapstructure:"log_level"`
	LogFileName     string `mapstructure:"log_file_name"`
	LogPath         string `mapstructure:"log_path"`
	LogmaxAge       int
	LogrotationTime int
	StartTask       bool `mapstructure:"start_task"`
	Version         string
}

type Redis struct {
	Addrs    string
	User     string
	Password string `json:"-"`
	Mode     enum.RedisMode
	Db       int
}

func ReadConfig(data []byte) (*Config, error) {
	v := viper.New()
	v.SetConfigType("toml")
	reader := bytes.NewReader(data)
	err := v.ReadConfig(reader)
	if err != nil {
		return nil, err
	}
	var conf Config
	if err := v.Unmarshal(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
