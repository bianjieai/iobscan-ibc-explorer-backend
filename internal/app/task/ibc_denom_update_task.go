package task

type IbcDenomUpdateTask struct {
}

var _ Task = new(IbcDenomUpdateTask)

func (t *IbcDenomUpdateTask) Name() string {
	return "ibc_denom_update_task"
}

func (t *IbcDenomUpdateTask) Cron() int {
	if taskConf.CronTimeDenomUpdateTask > 0 {
		return taskConf.CronTimeDenomUpdateTask
	}
	return EveryMinute
}

func (t *IbcDenomUpdateTask) Run() int {

}
