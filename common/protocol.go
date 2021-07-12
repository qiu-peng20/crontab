package common

import (
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

// Job 定时任务
type Job struct {
	Name     string `json:"name"`     //任务名字
	Command  string `json:"command"`  //shell命令
	CronExpr string `json:"cronExpr"` // 任务执行时间
}

// Response 接口返回策略
type Response struct {
	ErrCode int         `json:"errCode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// JobSchedulePlan JobSchedule schedule的任务调度计划
type JobSchedulePlan struct {
	Job      *Job
	Expr     *cronexpr.Expression
	NextTime time.Time
}

func NewResponse(errCode int, message string, data interface{}) ([]byte, error) {
	var response Response
	response.ErrCode = errCode
	response.Message = message
	response.Data = data

	marshal, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

func JsonParseJob(data []byte) (*Job, error) {
	var job Job
	err := json.Unmarshal(data, &job)
	if err != nil {
		return &job, err
	}
	return &job, nil
}

func FindKey(s string) (str string) {
	str = strings.TrimPrefix(s, JobSaveUrl)
	return
}

type JobEvent struct {
	EventType int
	Job       *Job
}

func BuildJobEvent(Event int, j *Job) *JobEvent {
	jobE := &JobEvent{
		EventType: Event,
		Job:       j,
	}
	return jobE
}

// BuildJobSchedulePlan 生成调度计划
func BuildJobSchedulePlan(j *Job) (jsp *JobSchedulePlan, err error) {
	var (
		parse *cronexpr.Expression
	)
	parse, err = cronexpr.Parse(j.CronExpr)
	if err != nil {
		return &JobSchedulePlan{}, err
	}
	jsp = &JobSchedulePlan{
		Job: j,
		Expr: parse,
		NextTime: parse.Next(time.Now()),
	}
	return
}

type JobScheduleExecuting struct {
	Job *Job
	PlanTime time.Time
	ExecutingTime time.Time
}

func BuildJobExecuting(jp *JobSchedulePlan) (je *JobScheduleExecuting)  {
	return &JobScheduleExecuting{
		Job: jp.Job,
		PlanTime: jp.NextTime,
		ExecutingTime: time.Now(),
	}
}

type JobExecutorResult struct {
	JSE *JobScheduleExecuting
	OutPut []byte //shell执行命令
	Err error
	StartTime time.Time //开始时间
	EndTime time.Time //结束时间
}
