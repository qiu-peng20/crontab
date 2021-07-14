package worker

import (
	"context"
	"crontab/common"
	config2 "crontab/worker/config"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

type JobMgr struct {
	Client  *clientv3.Client
	Kv      clientv3.KV
	Lease   clientv3.Lease
	Watcher clientv3.Watcher
}

var G_jobMgr *JobMgr

func InitJobMgr() (err error) {
	var (
		client  *clientv3.Client
		lease   clientv3.Lease
		kv      clientv3.KV
		watcher clientv3.Watcher
	)
	//配置etcd
	config := clientv3.Config{
		Endpoints:   config2.G_Config.EtcdEndPoints,
		DialTimeout: time.Duration(config2.G_Config.EtcdDialTimeOut) * time.Millisecond,
	}
	client, err = clientv3.New(config)
	if err != nil {
		return err
	}
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)

	G_jobMgr = &JobMgr{
		Client:  client,
		Kv:      kv,
		Lease:   lease,
		Watcher: watcher,
	}
	//启动一个监听携程
	err = G_jobMgr.WatchJob()
	if err != nil {
		return err
	}
	err = G_jobMgr.WatchKillJob()
	if err != nil {
		return err
	}
	return
}

// WatchKillJob 监听强杀任务
func (j JobMgr) WatchKillJob() (err error) {
	var (
		watchChan   clientv3.WatchChan
		eventJob *common.JobEvent
	)
	go func() {
		//监听/cron/job/目录下的所有的事件
		watchChan = j.Watcher.Watch(context.TODO(), common.JobKillUrl, clientv3.WithPrefix())
		//处理监听事件
		for watchResponse := range watchChan {
			for _, watchEvent := range watchResponse.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					//任务保存事件
					str := common.FindKillKey(string(watchEvent.Kv.Value))
					//生成调度计划
					eventJob = common.BuildJobEvent(common.KillJob, &common.Job{
						Name: str,
					})
				case mvccpb.DELETE:
					//任务删除事件
				}
				G_JobSchedule.PushSchedule(eventJob)
				//推给schedule
			}
		}
	}()
	return
}

func (j JobMgr) WatchJob() (err error) {
	var (
		getResponse *clientv3.GetResponse
		job         *common.Job
		watchChan   clientv3.WatchChan
		eventJob    *common.JobEvent
	)
	// get到目前/cron/job/目录下的所有任务，并且获取当前集群的版本号
	getResponse, err = j.Kv.Get(context.TODO(), common.JobSaveUrl, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, value := range getResponse.Kvs {
		//反序列化job
		job, err = common.JsonParseJob(value.Value)
		if err != nil {
			return err
		}
		//TODO：把JOB同步给任务调度携程
		eventJob = common.BuildJobEvent(common.SaveJob, job)
		G_JobSchedule.PushSchedule(eventJob)
	}

	//监听后续的变化
	go func() {
		revision := getResponse.Header.Revision + 1
		//监听/cron/job/目录下的所有的事件
		watchChan = j.Watcher.Watch(context.TODO(), common.JobSaveUrl, clientv3.WithRev(revision), clientv3.WithPrefix())
		//处理监听事件
		for watchResponse := range watchChan {
			for _, watchEvent := range watchResponse.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					job, _ = common.JsonParseJob(watchEvent.Kv.Value)
					//任务保存事件
					eventJob = common.BuildJobEvent(common.SaveJob, job)
					//构建一个event事件，推给schedule
				case mvccpb.DELETE:
					//任务删除事件
					b := common.FindKey(string(watchEvent.Kv.Key))
					eventJob = common.BuildJobEvent(common.DeleteJob, &common.Job{
						Name: b,
					})
				}
				G_JobSchedule.PushSchedule(eventJob)
				//推给schedule
			}
		}
	}()
	return
}

func (j JobMgr) CreateLock(name string) (lock *JobLock) {
	lock = InitJobLock(j.Kv, j.Lease, name)
	return
}
