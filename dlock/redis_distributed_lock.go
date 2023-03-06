package dlock

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"reentry-distributed-lock/dlock/redisutil"
)

const lockName = "RedisDistributedLock"

type redisDistributedLock struct {
	lockTrueName string
	lockId       string
	expireTime   int
}

func (lock *redisDistributedLock) Lock() {
	rdb, ctx := redisutil.GetRedisClient()

	for execLockLuaScript(rdb, &ctx, lock.lockId, lock.expireTime, lock.lockTrueName) == 0 {
		time.Sleep(20 * time.Millisecond)
	}

	go func(lock *redisDistributedLock) {
		ticker := time.NewTicker(time.Duration(int64(float64(lock.expireTime)*0.4)) *
			time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if execAutoRenewLuaScript(rdb, &ctx, lock.lockId, lock.expireTime,
				lock.lockTrueName) == 0 {
				break
			}
		}
	}(lock)
}

func (lock *redisDistributedLock) Unlock() {
	rdb, ctx := redisutil.GetRedisClient()

	execUnlockLuaScript(rdb, &ctx, lock.lockId, lock.expireTime, lock.lockTrueName)
}

func execLockLuaScript(rdb *redis.Client, ctx *context.Context,
	lockId string, expireTime int, lockTrueName string) int {

	locked, err := rdb.Eval(*ctx, redisutil.GetLockScript(),
		[]string{lockName + ":" + lockTrueName}, lockId, expireTime).Int()

	if err != nil {
		fmt.Println("加锁发生错误：" + err.Error())
		return 0
	}

	return locked
}

func execUnlockLuaScript(rdb *redis.Client, ctx *context.Context,
	lockId string, expireTime int, lockTrueName string) {

	unlocked, err := rdb.Eval(*ctx, redisutil.GetUnlockScript(),
		[]string{lockName + ":" + lockTrueName}, lockId, expireTime).Int()

	if err != nil {
		fmt.Println("解锁发生错误：" + err.Error())
	}

	if unlocked == 0 {
		fmt.Println("解锁Lua脚本执行错误")
	}
}

func execAutoRenewLuaScript(rdb *redis.Client, ctx *context.Context,
	lockId string, expireTime int, lockTrueName string) int {

	success, err := rdb.Eval(*ctx, redisutil.GetAutoRenewScript(),
		[]string{lockName + ":" + lockTrueName}, lockId, expireTime).Int()

	if err != nil {
		fmt.Println("自动续期发生错误：" + err.Error())
		return 0
	}

	return success
}
