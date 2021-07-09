package schedule

import (
	"crontab/common"
	"fmt"
	"time"
)

// JobSchedule 任务调度
type JobSchedule struct {
	ScheduleChan chan *common.JobEvent //ETCD任务队列
	SchedulePlan map[string]*common.JobSchedulePlan  //任务调度计划
	JobExecutingTable map[string]*common.JobScheduleExecuting  //任务执行计划
}

var G_JobSchedule *JobSchedule

//取出任务进行执行
func (j *JobSchedule) scheduleLoop() {
	var (
		jobEvent *common.JobEvent
	)
	//初始化任务调度器
	scheduleAfter :=  j.CheckSchedule()
	timer := time.NewTimer(scheduleAfter)
	for {
		select {
		case jobEvent = <-j.ScheduleChan:
			j.handleSchedule(jobEvent)
		case <-timer.C: //最近任务到期

		}
		scheduleAfter = j.CheckSchedule()
		timer.Reset(scheduleAfter) //重制调度间隔
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
		nearTime *time.Time
	)
	nowTime := time.Now() // 当前时间

	//如果任务列表为空时候，睡眠一秒
	if len(j.SchedulePlan) == 0 {
		scheduleAfter = 1 * time.Second
		return
	}

	//1. 遍历所有的执行任务
	for _,jobSchedulePlan = range j.SchedulePlan {
		//2. 到期的任务立即执行
		if jobSchedulePlan.NextTime.Equal(nowTime) || jobSchedulePlan.NextTime.Before(nowTime) {
			//TODO:尝试执行任务
			j.tryJob(jobSchedulePlan)
			jobSchedulePlan.NextTime = jobSchedulePlan.Expr.Next(nowTime)
		}
		//统计最近的需要执行的时间
		if nearTime == nil || jobSchedulePlan.NextTime.Before(*nearTime) {
			nearTime = &jobSchedulePlan.NextTime
		}
	}
	scheduleAfter = (*nearTime).Sub(nowTime)
	return
}

func (j *JobSchedule) tryJob(jp *common.JobSchedulePlan) {
	//调度和任务执行
	var (
		jobBool bool
		jobScheduleExecuting *common.JobScheduleExecuting
	)
	if jobScheduleExecuting, jobBool = j.JobExecutingTable[jp.Job.Name]; jobBool {
		fmt.Println("任务尚未推出，跳过执行")
		return
	}
	//构建执行任务
	jobScheduleExecuting = common.BuildJobExecuting(jp)
	//保存执行任务
	j.JobExecutingTable[jp.Job.Name] = jobScheduleExecuting
	//TODO:执行任务:
	fmt.Println("执行任务",j.JobExecutingTable[jp.Job.Name].ExecutingTime,j.JobExecutingTable[jp.Job.Name].PlanTime)
}

func (j *JobSchedule) PushSchedule(jobEvent *common.JobEvent) {
	j.ScheduleChan <- jobEvent
}

// InitSchedule 初始化一个任务调度器
func InitSchedule() {
	G_JobSchedule = &JobSchedule{
		ScheduleChan: make(chan *common.JobEvent, 1000),
		SchedulePlan: make(map[string]*common.JobSchedulePlan),
		JobExecutingTable: make(map[string]*common.JobScheduleExecuting),
	}
	go G_JobSchedule.scheduleLoop()
}
