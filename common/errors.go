package common

import "github.com/pkg/errors"

var (
	Lock_failure = errors.New("抢锁失败")
)
