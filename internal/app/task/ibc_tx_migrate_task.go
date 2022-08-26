package task

import (
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/sirupsen/logrus"
)

type IbcTxMigrateTask struct {
}

var _ Task = new(IbcTxMigrateTask)

func (t *IbcTxMigrateTask) Name() string {
	return "ibc_tx_migrate_task"
}

func (t *IbcTxMigrateTask) Switch() bool {
	return global.Config.Task.SwitchIbcTxMigrateTask
}

func (t *IbcTxMigrateTask) Cron() int {
	if taskConf.CronTimeIbcTxMigrateTask > 0 {
		return taskConf.CronTimeIbcTxMigrateTask
	}
	return EveryHour
}

func (t *IbcTxMigrateTask) Run() int {
	if !t.Switch() {
		logrus.Infof("task %s closed", t.Name())
		return 1
	}

	err1 := t.migrateSetting()
	err2 := t.migrateNormal()

	if err1 != nil || err2 != nil {
		return -1
	}
	return 1
}

func (t *IbcTxMigrateTask) migrateSetting() error {
	const limit = 1000
	status := []entity.IbcTxStatus{entity.IbcTxStatusSetting}
	totalMigrate := 0
	for {
		txList, err := ibcTxRepo.FindByStatus(status, limit)
		if err != nil {
			logrus.Errorf("task %s find setting txs error, %v", t.Name(), err)
			return err
		}

		if err = ibcTxRepo.Migrate(txList); err != nil {
			logrus.Errorf("task %s migrate setting txs error, %v", t.Name(), err)
			return err
		}

		totalMigrate += len(txList)
		if len(txList) < limit {
			break
		} else {
			time.Sleep(200 * time.Millisecond) // avoid master-slave delay problem
		}
	}

	logrus.Infof("task %s migrate %d setting txs", t.Name(), totalMigrate)
	return nil
}

func (t *IbcTxMigrateTask) migrateNormal() error {
	const limit = 1000
	status := entity.IbcTxUsefulStatus
	count, err := ibcTxRepo.CountByStatus(status)
	if err != nil {
		logrus.Errorf("task %s find count normal txs error, %v", t.Name(), err)
		return err
	}

	batch := (count - ibcTxCount) / limit
	if batch <= 0 {
		return nil
	}

	totalMigrate := 0
	for ; batch > 0; batch-- {
		txList, err := ibcTxRepo.FindByStatus(status, limit)
		if err != nil {
			logrus.Errorf("task %s find mormal txs error, %v", t.Name(), err)
			return err
		}

		if err = ibcTxRepo.Migrate(txList); err != nil {
			logrus.Errorf("task %s migrate mormal txs error, %v", t.Name(), err)
			return err
		}

		totalMigrate += len(txList)
		if len(txList) < limit {
			break
		} else {
			time.Sleep(200 * time.Millisecond) // avoid master-slave delay problem
		}
	}

	logrus.Infof("task %s migrate %d mormal txs", t.Name(), totalMigrate)
	return nil
}
