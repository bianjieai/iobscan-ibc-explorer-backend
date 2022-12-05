package conf

import (
	"bytes"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/enum"
	"github.com/spf13/viper"
)

type Config struct {
	App           App
	Mongo         Mongo
	Mysql         Mysql
	Redis         Redis
	Log           Log
	Spi           Spi
	HintIndexName HintIndexName `mapstructure:"hint_index_name"`
	Task          Task
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

type HintIndexName struct {
	GetRelayerTxsHintIndex           string `mapstructure:"get_relayer_txs_hint_index"`
	CountRelayerTxsHintIndex         string `mapstructure:"count_relayer_txs_hint_index"`
	GetRelayerUpdateTimeHintIndex    string `mapstructure:"get_relayer_update_time_hint_index"`
	GetLatestRecvPacketTimeHintIndex string `mapstructure:"get_latest_recv_packet_time_hint_index"`
}

type App struct {
	Name         string
	Addr         string
	Env          string
	StartMonitor bool `mapstructure:"start_monitor"`
	StartTask    bool `mapstructure:"start_task"`
	//StartOneOffTask      bool  `mapstructure:"start_one_off_task"`
	ApiCacheAliveSeconds int   `mapstructure:"api_cache_alive_seconds"`
	MaxPageSize          int64 `mapstructure:"max_page_size"`
	Version              string
	Prometheus           string `mapstructure:"prometheus_port"`
	EnableSignature      bool   `mapstructure:"enable_signature"`
	EnableRateLimit      bool   `mapstructure:"enable_rate_limit"`
	RateLimitPolicy      string `mapstructure:"rate_limit_policy"`
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
	IBCTxFailLogTask string `mapstructure:"ibc_tx_fail_log_task"`
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
