package dlock

import (
	"strconv"
	"strings"
	"time"
)

type DistributedLockFactory struct {
	ImplementType string
}

func (factory *DistributedLockFactory) GetDistributedLock(lockId string,
	expireTime int, lockTrueName string) DistributedLock {

	switch strings.ToLower(factory.ImplementType) {
	// TODO：添加更多实现方式的分布式锁（如：Zookeeper，MySQL，Etcd等）
	case "redis":
		return &redisDistributedLock{
			lockId:       lockId + ":" + strconv.FormatInt(time.Now().Unix(), 10),
			expireTime:   expireTime,
			lockTrueName: lockTrueName,
		}
	default:
		factory.ImplementType = "redis"
		return &redisDistributedLock{
			lockId:       lockId + ":" + strconv.FormatInt(time.Now().Unix(), 10),
			expireTime:   expireTime,
			lockTrueName: lockTrueName,
		}
	}
}
