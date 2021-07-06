package schedule

import (
	"crontab/common"
	"time"
)

// JobSchedule 任务调度
type JobSchedule struct {
	ScheduleChan chan *common.JobEvent //ETCD任务队列
	SchedulePlan map[string]*common.JobSchedulePlan  //任务调度计划
}

var G_JobSchedule *JobSchedule

//取出任务进行执行
func (j *JobSchedule) scheduleLoop() {
	var (
		jobEvent *common.JobEvent
	)
	for {
		select {
		case jobEvent = <-j.ScheduleChan:
			j.handleSchedule(jobEvent)
		}
	}
}

//处理调度任务
func (j *JobSchedule) handleSchedule(jobE *common.JobEvent) {
	var (
		err          error
		schedulePlan *common.JobSchedulePlan
		jobBool bool
	)
	switch jobE.EventType {
	case common.SaveJob:
		schedulePlan, err = common.BuildJobSchedulePlan(jobE.Job)
		if err != nil {
			return
		}
		j.SchedulePlan[schedulePlan.Job.Name] = schedulePlan
	case common.DeleteJob:
		if schedulePlan,jobBool = j.SchedulePlan[jobE.Job.Name];jobBool {
			delete(j.SchedulePlan,schedulePlan.Job.Name)
		}
	}
}

// CheckSchedule 遍历需要执行的任务
func (j *JobSchedule) CheckSchedule() (scheduleAfter time.Duration) {
	var (
		jobSchedulePlan *common.JobSchedulePlan
	)
	//1. 遍历所有的执行任务
	for _,jobSchedulePlan = range j.SchedulePlan {
		
	}
	//2. 到期的任务立即执行

	//3. 统计最近的任务要过期的时间
}

func (j *JobSchedule) PushSchedule(jobEvent *common.JobEvent) {
	j.ScheduleChan <- jobEvent
}

// InitSchedule 初始化一个任务调度器
func InitSchedule() {
	G_JobSchedule = &JobSchedule{
		ScheduleChan: make(chan *common.JobEvent, 1000),
		SchedulePlan: make(map[string]*common.JobSchedulePlan),
	}
	go G_JobSchedule.scheduleLoop()
}
