package lock

import (
	"XcxcVideo/common/helper"
	"XcxcVideo/common/models"
	"context"
	"fmt"
	"os"
	"time"
)

type RedisLock struct {
	LockId string
	Key    string
}

var unlockScript string
var maxRetryTime = 1 * time.Second

func init() {
	result, err := os.ReadFile("common/lock/unlock.lua")
	if err != nil {
		panic(err)
	}
	unlockScript = string(result)
}
func (lock *RedisLock) TryLock(duration time.Duration) bool {
	lock.LockId = helper.GetUUID()
	result, _ := models.RDb.SetNX(context.Background(), lock.Key, lock.LockId, duration).Result()
	return result
}
func (lock *RedisLock) TryLockWithRetry(duration time.Duration) bool {
	lock.LockId = helper.GetUUID()
	resultCh := make(chan bool)
	timeoutCh := time.After(maxRetryTime) // 设置超时通道
	go func() {
		for {
			// 尝试获取锁
			result, _ := models.RDb.SetNX(context.Background(), lock.Key, lock.LockId, duration).Result()
			if result {
				resultCh <- true
				return
			}
			// 重试
			time.Sleep(50 * time.Millisecond)
		}
	}()

	select {
	case result := <-resultCh:
		return result
	case <-timeoutCh:
		return false
	}
}

func (lock *RedisLock) UnLock() {
	//if models.RDb.Get(context.Background(), lock.Key).Val() == lock.LockId {
	//	models.RDb.Del(context.Background(), lock.Key)
	//}
	result, err := models.RDb.Eval(context.Background(), string(unlockScript), []string{lock.Key}, lock.LockId).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else if result == int64(1) {
		fmt.Println("Lock released successfully")
	} else {
		fmt.Println("Failed to release lock (either not owned or not present)")
	}

}
