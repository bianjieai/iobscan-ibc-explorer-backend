package conf

import (
	"bytes"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/enum"
	"github.com/spf13/viper"
)

type Config struct {
	App   App
	Mongo Mongo
	Mysql Mysql
	Redis Redis
	Log   Log
	Spi   Spi
	Task  Task
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
	Name                 string
	Addr                 string
	Env                  string
	StartTask            bool `mapstructure:"start_task"`
	StartOneOffTask      bool `mapstructure:"start_one_off_task"`
	ApiCacheAliveSeconds int  `mapstructure:"api_cache_alive_seconds"`
	Version              string
	Prometheus           string `mapstructure:"prometheus_port"`
}

type Redis struct {
	Addrs    string
	User     string
	Password string `json:"-"`
	Mode     enum.RedisMode
	Db       int
}

type Log struct {
	LogLevel           string `mapstructure:"log_level"`
	LogFileName        string `mapstructure:"log_file_name"`
	LogPath            string `mapstructure:"log_path"`
	LogMaxAgeDay       int    `mapstructure:"log_max_age_day"`
	LogRotationTimeDay int    `mapstructure:"log_rotation_time_day"`
	LogOutput          string `mapstructure:"log_output"`
}

type Task struct {
	CronJobRelayerAddr     string `mapstructure:"cron_job_relayer_addr"`
	CronTimeChainTask      int    `mapstructure:"cron_time_chain_task"`
	CronTimeChannelTask    int    `mapstructure:"cron_time_channel_task"`
	CronTimeRelayerTask    int    `mapstructure:"cron_time_relayer_task"`
	CronTimeTokenTask      int    `mapstructure:"cron_time_token_task"`
	CronTimeTokenPriceTask int    `mapstructure:"cron_time_token_price_task"`
	RedisLockExpireTime    int    `mapstructure:"redis_lock_expire_time"`
}

type Spi struct {
	CoingeckoPriceUrl string `mapstructure:"coingecko_price_url"`
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
