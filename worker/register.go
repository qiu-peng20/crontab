package worker

import (
	"crontab/worker/config"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"net"
	"time"
)

type Register struct {
	Client  *clientv3.Client
	Kv      clientv3.KV
	Lease   clientv3.Lease
	LocalIp string
}

var G_register *Register

//获取网卡ip
func getLocalIp() (localIp string, err error) {
	var (
		addrs []net.Addr
	)
	addrs, err = net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		fmt.Println(addr)
	}
	return
}

func InitRegister() (err error) {
	var (
		client *clientv3.Client
		ip string
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
	ip, err = getLocalIp()
	if err != nil {
		return err
	}
	G_register = &Register{
		Client: client,
		Kv:     kv,
		Lease:  lease,
		LocalIp: ip,
	}
	return
}
