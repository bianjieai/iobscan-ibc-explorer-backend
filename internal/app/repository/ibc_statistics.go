package repository

import (
	"context"
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	officialOpts "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type IStatisticRepo interface {
	FindOne(statisticName string) (entity.IbcStatistic, error)
	UpdateOne(statisticName string, count int64) error
	UpdateOneIncre(statistic entity.IbcStatistic) error
	FindBatchName(statisticNames []string) ([]*entity.IbcStatistic, error)
	Save(data entity.IbcStatistic) error
	CreateNew() error
	SwitchColl() error
	InsertCountNew(statisticName string, count int64) error
	InsertNew(data entity.IbcStatistic) error
}

var _ IStatisticRepo = new(IbcStatisticRepo)

type IbcStatisticRepo struct {
}

func (repo *IbcStatisticRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCStatisticsCollName)
}
func (repo *IbcStatisticRepo) collNew() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCStatisticsNewCollName)
}

func (repo *IbcStatisticRepo) UpdateOne(statisticName string, count int64) error {
	err := repo.coll().UpdateOne(context.Background(), bson.M{"statistics_name": statisticName}, bson.M{
		"$set": bson.M{
			"count":     count,
			"update_at": time.Now().Unix(),
		}})
	if err == qmgo.ErrNoSuchDocuments {
		data := entity.IbcStatistic{
			StatisticsName: statisticName,
			Count:          count,
			CreateAt:       time.Now().Unix(),
			UpdateAt:       time.Now().Unix(),
		}
		return repo.Save(data)
	}
	return err
}

func (repo *IbcStatisticRepo) FindOne(statisticName string) (entity.IbcStatistic, error) {
	var res entity.IbcStatistic
	err := repo.coll().Find(context.Background(), bson.M{"statistics_name": statisticName}).One(&res)
	return res, err
}

func (repo *IbcStatisticRepo) FindBatchName(statisticNames []string) ([]*entity.IbcStatistic, error) {
	var res []*entity.IbcStatistic
	err := repo.coll().Find(context.Background(), bson.M{"statistics_name": bson.M{
		"$in": statisticNames,
	}}).All(&res)
	return res, err
}

func (repo *IbcStatisticRepo) Save(data entity.IbcStatistic) error {
	_, err := repo.coll().InsertOne(context.Background(), data)
	return err
}

func (repo *IbcStatisticRepo) UpdateOneIncre(statistic entity.IbcStatistic) error {
	err := repo.coll().UpdateOne(context.Background(), bson.M{"statistics_name": statistic.StatisticsName}, bson.M{
		"$set": bson.M{
			"count":           statistic.Count,
			"count_latest":    statistic.CountLatest,
			"statistics_info": statistic.StatisticsInfo,
			"update_at":       time.Now().Unix(),
		}})
	if err == qmgo.ErrNoSuchDocuments {
		return repo.Save(statistic)
	}
	return err
}

func (repo *IbcStatisticRepo) CreateNew() error {
	indexOpts := officialOpts.Index().SetUnique(true).SetName("statistics_unique")
	key := []string{"statistics_name"}
	return repo.collNew().CreateOneIndex(context.Background(), opts.IndexModel{Key: key, IndexOptions: indexOpts})
}

func (repo *IbcStatisticRepo) SwitchColl() error {
	command := bson.D{{"renameCollection", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCStatisticsNewCollName)},
		{"to", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCStatisticsCollName)},
		{"dropTarget", true}}
	return mgo.Database(adminDatabase).RunCommand(context.Background(), command).Err()
}

func (repo *IbcStatisticRepo) InsertCountNew(statisticName string, count int64) error {

	data := entity.IbcStatistic{
		StatisticsName: statisticName,
		Count:          count,
		CreateAt:       time.Now().Unix(),
		UpdateAt:       time.Now().Unix(),
	}
	return repo.InsertNew(data)
}
func (repo *IbcStatisticRepo) InsertNew(data entity.IbcStatistic) error {
	_, err := repo.collNew().InsertOne(context.Background(), data)
	return err
}
