package jobManager

import (
	"context"
	"crontab/common"
	config2 "crontab/worker/config"
	"github.com/coreos/etcd/clientv3"
	"time"
)

type JobMgr struct {
	Client *clientv3.Client
	Kv     clientv3.KV
	Lease  clientv3.Lease
}

var G_jobMgr *JobMgr

func InitJobMgr() (err error) {
	var (
		client *clientv3.Client
		lease  clientv3.Lease
		kv     clientv3.KV
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

	G_jobMgr = &JobMgr{
		Client: client,
		Kv:     kv,
		Lease:  lease,
	}
	return
}

func (j JobMgr) WatchJob() (err error) {
	var (
		getResponse *clientv3.GetResponse
		job *common.Job
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
	}

	//监听后续的变化
	go func() {

	}()
	return
}
