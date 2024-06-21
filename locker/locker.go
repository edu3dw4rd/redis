package locker

import (
	"fmt"
	"time"

	"github.com/edu3dw4rd/redis"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis"
)

var (
	lockerManager *LockerManager
)

type LockerManager struct {
	redisSync *redsync.Redsync
}

func NewLockerManager() *LockerManager {
	if lockerManager == nil || lockerManager.redisSync == nil {
		// Get redis client
		rdsClient := redis.GetRedisClient()

		pool := goredis.NewPool(rdsClient)
		lockerManager = &LockerManager{
			redisSync: redsync.New(pool),
		}
	}

	return lockerManager
}

func (l *LockerManager) ObtainLock(key string, retryNum int, ttl, retryDelay time.Duration) (*redsync.Mutex, error) {
	var mutex *redsync.Mutex
	options := l.setRedsyncOptions(retryNum, ttl, retryDelay)

	if len(options) > 0 {
		mutex = l.redisSync.NewMutex(key, options...)
	} else {
		mutex = l.redisSync.NewMutex(key)
	}

	if err := mutex.Lock(); err != nil {
		return nil, fmt.Errorf("failed to obtain lock [%s]: %s", mutex.Name(), err.Error())
	}

	return mutex, nil
}

func (l *LockerManager) ReleaseLock(mutex *redsync.Mutex) error {
	if ok, err := mutex.Unlock(); !ok || err != nil {
		return fmt.Errorf("failed to release lock [%s]: %s", mutex.Name(), err.Error())
	}

	return nil
}

func (l *LockerManager) setRedsyncOptions(retryNum int, ttl, retryDelay time.Duration) (options []redsync.Option) {
	if retryNum > 0 {
		options = append(options, redsync.WithTries(retryNum))
	}

	if ttl > 0 {
		options = append(options, redsync.WithExpiry(ttl))
	}

	if retryDelay > 0 {
		options = append(options, redsync.WithRetryDelay(retryDelay))
	}

	return
}
