package entity

import (
	"fmt"
	"time"
)

type SyncTaskStatus string

const (
	SyncTaskStatusUnderway  = "underway"
	SyncTaskStatusInvalid   = "invalid"
	SyncTaskStatusCompleted = "completed"
)

type SyncTask struct {
	Startheight    int64       `bson:"start_height"`
	EndHeight      int64       `bson:"end_height"`
	CurrentHeight  int64       `bson:"current_height"`
	Status         string      `bson:"status"`
	WorkerId       string      `bson:"worker_id"`
	WorkerLogs     []WorkerLog `bson:"worker_logs"`
	LastUpdateTime int64       `bson:"last_update_time"`
}

type WorkerLog struct {
	WorkerId  string    `bson:"worker_id"`
	BeginTime time.Time `bson:"begin_time"`
}

func (s SyncTask) CollectionName(chain string) string {
	return fmt.Sprintf("sync_%s_task", chain)
}
