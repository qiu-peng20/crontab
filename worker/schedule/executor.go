package schedule

import (
	"crontab/common"
	"os/exec"
	"time"
)

// Executor 执行器实现
type Executor struct{}

var G_executor *Executor

func InitExecutor() {
	G_executor = &Executor{}
}

func (e *Executor) ExecuteJob(job *common.JobScheduleExecuting) {
	var (
		cmd *exec.Cmd
		result *common.JobExecutorResult
		err error
	)
	go func() {
		//初始化任务结果
		result = &common.JobExecutorResult{
			JSE: job,
			OutPut: make([]byte,0),
		}
		//首先获取分布式锁
		lock := G_jobMgr.CreateLock(job.Job.Name)

		//执行shell命令
		result.StartTime = time.Now()
		err = lock.TryLock()
		defer lock.RemoveLock()
		if err != nil {
			result.EndTime = time.Now()
			return
		}else {
			result.StartTime = time.Now()

			cmd = exec.CommandContext(job.Ctx, "/bin/bash", "-c", job.Job.Command)
			//获取shell命令
			output, err := cmd.Output()

			result.EndTime = time.Now()
			result.OutPut = output
			result.Err = err
			//将执行结果回传给schedule
		}
		G_JobSchedule.PushResult(result)
	}()
}
