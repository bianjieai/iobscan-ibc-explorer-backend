package repository

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/qiniu/qmgo"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var mgo *qmgo.Client
var ibcDatabase string

func InitMysqlDB(cfg conf.Mysql) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
		//encodeTimeZone(cfg.TimeZone), use mysql default timezone(UTC)
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("start mysql client failed, db:%s, err:%s", cfg.Database, err.Error())
	}
	return db
}

func GetDB() *gorm.DB {
	return db
}

func DbStatus() bool {
	d, err := db.DB()
	if err != nil {
		return false
	}
	return d.Ping() == nil
}

func encodeTimeZone(timezone string) string {
	if timezone == "" {
		timezone = constant.DefaultTimezone
	}

	return url.QueryEscape(fmt.Sprintf("'%s'", timezone))
}

func CreateTable(db *gorm.DB) {
	_ = db.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate()
}

func InitMgo(cfg conf.Mongo, ctx context.Context) {
	var maxPoolSize uint64 = 4096
	client, err := qmgo.NewClient(ctx, &qmgo.Config{
		Uri: cfg.Url,
		ReadPreference: &qmgo.ReadPref{
			MaxStalenessMS: 90000,
			Mode:           readpref.SecondaryPreferredMode,
		},
		Database:    cfg.Database,
		MaxPoolSize: &maxPoolSize,
	})
	if err != nil {
		logrus.Fatalf("connect mongo failed, uri: %s, err:%s", cfg.Url, err.Error())
	}
	mgo = client
	ibcDatabase = cfg.Database
}
