package schedule

import (
	"crontab/common"
	"fmt"
)

// JobSchedule 任务调度
type JobSchedule struct {
	ScheduleChan chan *common.JobEvent //ETCD任务队列
}

var G_JobSchedule *JobSchedule

func (j JobSchedule) scheduleLoop() {
	for job := range j.ScheduleChan{
		switch job.EventType {
		case common.SaveJob:
			fmt.Print(111)
		case common.DeleteJob:
			fmt.Print(222)
		}
	}
	//var (
	//	jobEvent *common.JobEvent
	//)
	//for {
	//	select {
	//	case jobEvent = <-j.ScheduleChan:
	//		switch jobEvent.EventType {
	//		case common.SaveJob:
	//			fmt.Print(111)
	//		case common.DeleteJob:
	//			fmt.Print(222)
	//		}
	//	}
	//}
}

func (j JobSchedule) PushSchedule(jobEvent *common.JobEvent) {
	j.ScheduleChan <- jobEvent
}

// InitSchedule 初始化一个任务调度器
func InitSchedule() {
	G_JobSchedule = &JobSchedule{
		ScheduleChan: make(chan *common.JobEvent, 1000),
	}
	go G_JobSchedule.scheduleLoop()
}
