package common

import (
	"context"
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

func FindKillKey(s string) (str string) {
	str = strings.TrimPrefix(s, JobKillUrl)
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
		Job:      j,
		Expr:     parse,
		NextTime: parse.Next(time.Now()),
	}
	return
}

type JobScheduleExecuting struct {
	Job           *Job
	PlanTime      time.Time
	ExecutingTime time.Time
	Ctx           context.Context
	Cancel        context.CancelFunc //用于取消任务的时候用
}

func BuildJobExecuting(jp *JobSchedulePlan) (je *JobScheduleExecuting) {
	ctx, cancel := context.WithCancel(context.TODO())
	return &JobScheduleExecuting{
		Job:           jp.Job,
		PlanTime:      jp.NextTime,
		ExecutingTime: time.Now(),
		Ctx:           ctx,
		Cancel:        cancel,
	}
}

type JobExecutorResult struct {
	JSE       *JobScheduleExecuting
	OutPut    []byte //shell执行命令
	Err       error
	StartTime time.Time //开始时间
	EndTime   time.Time //结束时间
}

// JobLog MongoDB数据结构
type JobLog struct {
	JobName      string `bson:"job_name"`      //任务执行名字
	Common       string `bson:"common"`        //Shell命令
	Err          string `bson:"err"`           //错误内容
	OutPut       string `bson:"out_put"`       //任务执行结果
	PlanTime     int64  `bson:"plan_time"`     //计划开始时间
	ScheduleTime int64  `bson:"schedule_time"` //实际调度时间
	StartTime    int64  `bson:"start_time"`    //任务开始时间
	EndTime      int64  `bson:"end_time"`      //任务结束时间
}

type LogBash struct {
	Logs [] interface{}
}
