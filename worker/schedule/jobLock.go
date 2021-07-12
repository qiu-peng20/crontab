package schedule

import (
	"context"
	"github.com/coreos/etcd/clientv3"
)

// JobLock 获取分布式锁
type JobLock struct {
	Kv    clientv3.KV
	Lease clientv3.Lease

	JobName string
}

// InitJobLock 初始化一个锁
func InitJobLock(kv clientv3.KV, lease clientv3.Lease, name string) (jl *JobLock) {
	jl = &JobLock{
		Kv:      kv,
		Lease:   lease,
		JobName: name,
	}
	return
}

func (jl *JobLock) TryLock() (err error) {
	var (
		grant     *clientv3.LeaseGrantResponse
		ctx       context.Context
		cancel    context.CancelFunc
		keepAlive <-chan *clientv3.LeaseKeepAliveResponse
	)
	//1.创建租约
	grant, err = jl.Lease.Grant(context.TODO(), 5)
	if err != nil {
		return err
	}
	leaseId := grant.ID
	//2.自动续租
	//2.1  创建context，用于取消续租
	ctx, cancel = context.WithCancel(context.TODO())
	//2.2 续租开始
	keepAlive, err = jl.Lease.KeepAlive(ctx, leaseId)
	if err != nil {
		goto FALL
	}
	//2.3 处理续租的携程
	go func() {
		for  {
			select {
			case alive := <-keepAlive: //自动应答
			if alive == nil {  //如果续租失败则为空
				goto END
			}

			}
		}
		END:
	}()
	//3.创建事务 TXN

	//4.抢锁
	//5.成功则返回，失败释放租约
	FALL:
		cancel()//取消自动续租
	jl.Lease.Revoke(context.TODO(),leaseId) //释放租约
	return
}
