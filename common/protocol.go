package common

import (
	"encoding/json"
	"strings"
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
