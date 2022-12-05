package app

import (
	"context"
	"os"
	"path"
	"strings"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/pkg/distributiontask"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/task"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

func Serve(cfg *conf.Config) {
	time.Local = time.UTC
	initCore(cfg)
	defer repository.Close()

	if cfg.App.ApiCacheAliveSeconds > 0 {
		api.SetApiCacheAliveTime(cfg.App.ApiCacheAliveSeconds)
	}

	r := gin.Default()
	api.Routers(r)
	//if cfg.App.StartMonitor {
	//	go monitor.Start(cfg.App.Prometheus)
	//}
	if cfg.App.StartTask {
		go startTask(cfg.Redis, cfg.Task)
	}
	logrus.Fatal(r.Run(cfg.App.Addr))
}

func initCore(cfg *conf.Config) {
	global.Config = cfg
	initLogger(&cfg.Log)
	repository.InitMgo(cfg.Mongo, context.Background())
	repository.LoadIndexNameConf(cfg.HintIndexName)
	cache.InitRedisClient(cfg.Redis)
}

func startTask(c conf.Redis, tc conf.Task) {
	task.LoadTaskConf(tc)

	distributionTask, err := distributiontask.NewDistributedTaskWithRedis(c.Addrs, c.User, c.Password, string(c.Mode), c.Db)
	if err != nil {
		logrus.Fatal(err)
	}

	distributionTask.RegisterTasks(new(task.IBCTxFailLogTask))
	distributionTask.Start()
}

func initLogger(logCfg *conf.Log) {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   constant.DefaultTimeFormat,
		DisableHTMLEscape: true,
	})
	if level, err := logrus.ParseLevel(logCfg.LogLevel); err == nil {
		logrus.SetLevel(level)
	}

	if strings.ToUpper(logCfg.LogOutput) == "FILE" {
		if _, err := os.Stat(logCfg.LogPath); os.IsNotExist(err) {
			_ = os.MkdirAll(logCfg.LogPath, os.ModePerm)
		}
		baseLogPath := path.Join(logCfg.LogPath, logCfg.LogFileName)
		writer, err := rotatelogs.New(
			baseLogPath+"_%Y%m%d.log",
			rotatelogs.WithLinkName(baseLogPath),                                               // 生成软链，指向最新日志文件
			rotatelogs.WithMaxAge(time.Duration(logCfg.LogMaxAgeDay*24)*time.Hour),             // 文件最大保存时间
			rotatelogs.WithRotationTime(time.Duration(logCfg.LogRotationTimeDay*24)*time.Hour), // 日志切割时间间隔
		)
		if err != nil {
			logrus.Fatalf("config local file system logger error. %s", err.Error())
		}

		logrus.SetOutput(writer)
	} else {
		logrus.SetOutput(os.Stdout)
	}
}
