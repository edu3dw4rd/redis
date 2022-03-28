package redis

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/joho/godotenv/autoload"
)

type redisUtil struct {
	rds *redis.Client
}

var redisOnce sync.Once
var redisInstance *redisUtil

func newRedisClient() *redis.Client {
	redisOnce.Do(func() {
		redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
		redisDb, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
		client := redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: "",
			DB:       redisDb,
		})

		redisInstance = &redisUtil{
			rds: client,
		}
	})

	return redisInstance.rds
}

// RedisSet sets/enter data into redis DB
func RedisSet(key string, value interface{}, expiration time.Duration) error {
	redisClient := newRedisClient()
	expDuration := expiration * time.Second

	err := redisClient.Set(key, value, expDuration).Err()

	if err != nil {
		return err
	}

	return nil
}

// RedisGet gets data with key from redisDB
func RedisGet(key string) (interface{}, error) {
	redisClient := newRedisClient()

	data, err := redisClient.Get(key).Result()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return data, nil
}

// RedisHGet gets data with Redis Command HGET based on given key and field
func RedisHGet(key string, field string) (interface{}, error) {
	redisClient := newRedisClient()

	data, err := redisClient.HGet(key, field).Result()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return data, nil
}

// RedisHSet sets data with Redis Command HSET based on given key, field, and value
func RedisHSet(key string, field string, value interface{}) error {
	redisClient := newRedisClient()
	err := redisClient.HSet(key, field, value).Err()

	if err == redis.Nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

// RedisIncr increment data by one with Redis Command INCR based on give key
func RedisIncr(key string) error {
	redisClient := newRedisClient()
	err := redisClient.Incr(key).Err()

	if err != nil {
		return err
	}

	return err
}

// RedisIncrby increment data with Redis Command INCRBY based on given key and value
func RedisIncrby(key string, value int64) error {
	redisClient := newRedisClient()
	err := redisClient.IncrBy(key, value).Err()

	if err == redis.Nil {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

// RedisDel is function to send a data array from redis database
// Return array data and error redis
func RedisDel(key string) error {
	redisClient := newRedisClient()
	err := redisClient.Del(key).Err()

	if err != nil {
		return err
	}

	return nil
}
