package saleapi

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"

	"reentry-distributed-lock/dlock"
)

var (
	rdb *redis.Client
	ctx context.Context

	distributedLockFactory dlock.DistributedLockFactory
)

const (
	expireTime   = 30
	lockTrueName = "Sale"
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.182.100:6379",
		Password: "redis",
		DB:       0,
	})
	ctx = context.Background()

	distributedLockFactory = dlock.DistributedLockFactory{ImplementType: "redis"}
}

func Sale(c *gin.Context) {
	id, err := uuid.NewV4()
	if err != nil {
		c.JSON(200, gin.H{
			"error": "获取库存失败",
		})

		fmt.Println("获取uuid失败：" + err.Error())
		return
	}

	lock := distributedLockFactory.GetDistributedLock(id.String(), expireTime, lockTrueName)
	lock.Lock()
	storeStr, err := rdb.Get(ctx, "store").Result()
	if err != redis.Nil && err != nil {
		lock.Unlock()
		c.JSON(200, gin.H{
			"error": "获取库存失败：" + err.Error(),
		})
		return
	}

	storeNum, _ := strconv.Atoi(storeStr)
	if storeNum <= 0 {
		lock.Unlock()
		c.JSON(200, gin.H{
			"message": "库存为0",
		})
		return
	}

	storeNum -= 1
	rdb.Set(ctx, "store", storeNum, 0)
	lock.Unlock()

	c.JSON(200, gin.H{
		"message": "购买成功，库存还剩：" + strconv.Itoa(storeNum),
	})
}
