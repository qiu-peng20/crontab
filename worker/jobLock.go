package worker

import (
	"context"
	"crontab/common"
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

// JobLock 获取分布式锁
type JobLock struct {
	Kv    clientv3.KV
	Lease clientv3.Lease

	JobName string
	LeaseId clientv3.LeaseID
	Cancel  context.CancelFunc
	LockBool bool
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
		txn       clientv3.Txn
		lockName  string
		commit    *clientv3.TxnResponse
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
		for {
			select {
			case alive := <-keepAlive: //自动应答
				if alive == nil { //如果续租失败则为空
					goto END
				} else {
					fmt.Println("自动续租成功", alive.ID)
				}
			}
		}
	END:
		fmt.Println("续租失败")
	}()
	//3.创建事务 TXN
	txn = jl.Kv.Txn(context.TODO())
	//设置锁路径
	lockName = common.JobLockUrl + jl.JobName
	//事务抢锁
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockName), "=", 0)).
		Then(clientv3.OpPut(lockName, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockName))

	// 提交事物
	if commit, err = txn.Commit(); err != nil {
		goto FALL
	}

	//走then -> success，else -> !success
	if !commit.Succeeded {
		err = common.Lock_failure
		goto FALL
	}
	jl.LeaseId = leaseId
	jl.Cancel = cancel
	jl.LockBool = true
	return
FALL:
	cancel()
	jl.Lease.Revoke(context.TODO(), leaseId)
	return
}

func (jl *JobLock)RemoveLock()  {
	if jl.LockBool {
		fmt.Println("取消锁")
		jl.Cancel()
		jl.Lease.Revoke(context.TODO(),jl.LeaseId)
	}
}
