package worker

import (
	"crontab/worker/config"
	"github.com/coreos/etcd/clientv3"
	"time"
)

type Register struct {
	Client  *clientv3.Client
	Kv      clientv3.KV
	Lease   clientv3.Lease
	LocalIp string
}

var G_register *Register

func InitRegister() (err error) {
	var (
		client *clientv3.Client
	)
	c := clientv3.Config{
		Endpoints:   config.G_Config.EtcdEndPoints,
		DialTimeout: time.Duration(config.G_Config.EtcdDialTimeOut) * time.Millisecond,
	}
	client, err = clientv3.New(c)
	if err != nil {
		return
	}
	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)
	G_register = &Register{
		Client: client,
		Kv:     kv,
		Lease:  lease,
	}
	return
}
