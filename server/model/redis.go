package model

import (
	"context"
	"fmt"
	"log"
	"syspulse/common"
	"time"

	redis "github.com/go-redis/redis/v8"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", common.SysArgs.Cache.Host, common.SysArgs.Cache.Port),
		DB:   common.SysArgs.Cache.DBIndex,
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
		return
	}
	log.Println("Connected to Redis:", pong)
}

func CacheSet(key string, value string, expiration time.Duration) {
	err := client.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		log.Println("Failed to set key:", err)
		return
	}
}

func CacheGet(key string) string {
	val, err := client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		log.Printf("There is no value with key '%s'", key)
		return ""
	} else if err != nil {
		log.Panic("Failed to get key:", err)
	}
	return val
}

func CacheExpire(key string, expiration time.Duration) {
	err := client.Expire(context.Background(), key, expiration).Err()
	if err != nil {
		log.Println("Failed to set key:", err)
	}
}
