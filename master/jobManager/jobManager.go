package jobManager

import (
	"context"
	"crontab/common"
	config2 "crontab/master/config"
	"encoding/json"
	"fmt"
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

func (j JobMgr) SaveJob(job common.Job) (oldJob *common.Job, err error) {
	var (
		jobKey   string
		jobValue []byte
		oldJobObj common.Job
		putResponse *clientv3.PutResponse
	)

	jobKey = "/cron/job/"
	//任务信息的json，序列化
	jobValue, err = json.Marshal(job)
	if err != nil {
		return nil, err
	}
	putResponse, err = j.Kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV())
	if err != nil {
		return nil, err
	}
	//如果更新成功，则返回旧的值
	if putResponse.PrevKv != nil {
		//对旧值进行一个反序列化
		err := json.Unmarshal(putResponse.PrevKv.Value, &oldJobObj)
		if err != nil {
			fmt.Print(err)
			err = nil
			return nil, err
		}
		oldJob = &oldJobObj
	}
	return
}
