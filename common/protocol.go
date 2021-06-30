package common

//定时任务
type Job struct {
	Name     string `json:"name"`     //任务名字
	Command  string `json:"command"`  //shell命令
	CronExpr string `json:"cronExpr"` // 任务执行时间
}


