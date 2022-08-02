package entity

const TaskNameFmt = "sync_%s_transfer"

type TaskRecordStatus string

var (
	TaskRecordStatusOpen  TaskRecordStatus = "open"
	TaskRecordStatusClose TaskRecordStatus = "close"
)

type IbcTaskRecord struct {
	TaskName string           `bson:"task_name"`
	Height   int64            `bson:"height"`
	Status   TaskRecordStatus `bson:"status"`
	CreateAt int64            `bson:"create_at"`
	UpdateAt int64            `bson:"update_at"`
}

func (t IbcTaskRecord) CollectionName() string {
	return "ibc_task_record"
}
