package model

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/syspulse/common"
	"go.uber.org/zap"

	redis "github.com/go-redis/redis/v8"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", common.SysArgs.Cache.Host, common.SysArgs.Cache.Port),
		DB:       common.SysArgs.Cache.DBIndex,
		Password: common.SysArgs.Cache.Passwd,
	})
	pong, err := client.Ping(client.Context()).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
		return
	}
	log.Println("Connected to Redis:", pong)
}

func CacheSet(key string, value string, expiration time.Duration) {
	err := client.Set(client.Context(), key, value, expiration).Err()
	if err != nil {
		log.Println("Failed to set key:", err)
		return
	}
}

func CacheGet(key string) string {
	val, err := client.Get(client.Context(), key).Result()
	if err == redis.Nil {
		zap.L().Warn("There is no value with key", zap.String("key", key))
		return ""
	} else if err != nil {
		zap.L().Panic("Failed to get key:", zap.Error(err))
	}
	return val
}

func CacheExpire(key string, expiration time.Duration) {
	err := client.Expire(client.Context(), key, expiration).Err()
	if err != nil {
		panic(err)
	}
}

func CacheAdd2Set(key string, member ...interface{}) {
	client.SAdd(client.Context(), key, member...)
}

func CacheAdd2HSetNX(key string, field string, value any) bool {
	result, err := client.HSetNX(client.Context(), key, field, value).Result()
	if err != nil {
		zap.L().Panic("error add field to hash in cache", zap.Error(err))
	}
	return result
}

func CacheHMSet(key string, entry any) {
	err := client.HMSet(client.Context(), key, entry).Err()
	if err != nil {
		panic(err)
	}
}

func CacheHGetAll(key string) map[string]string {
	result, err := client.HGetAll(client.Context(), key).Result()
	if err != nil {
		panic(err)
	}
	return result
}

func CacheExists(key string) bool {
	result, err := client.Exists(client.Context(), key).Result()
	if err != nil {
		zap.L().Panic("error check key exists.", zap.String("key", key))
	}
	return result == 1
}

func CacheHGet(key string, field string) string {
	result, err := client.HGet(client.Context(), key, field).Result()
	if err != nil {
		log.Default().Println(err)
		log.Default().Printf("key %s, field %s\n", key, field)
	}
	return result
}

func CacheHSet(key, field string, value any) int64 {
	result, err := client.HSet(client.Context(), key, field, value).Result()
	if err != nil {
		zap.L().Panic("error cache hset in cache", zap.Error(err))
	}
	return result
}

func CacheHDel(key string, field ...string) bool {
	result, err := client.HDel(client.Context(), key, field...).Result()
	if err != nil {
		zap.L().Panic("error cache hdel in cache", zap.Error(err))
	}
	return result == int64(len(field))
}

func CacheLPUSH(key string, value any) int64 {
	result, err := client.LPush(client.Context(), key, value).Result()
	if err != nil {
		log.Default().Println(err)
	}
	return result
}

// 尝试获取锁
func AcquireLock(lockKey string, lockValue string, expiration time.Duration) (bool, error) {
	// 使用SETNX命令实现锁，如果key不存在则设置成功并返回true，否则返回false
	ok, err := client.SetNX(client.Context(), lockKey, lockValue, expiration).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

// 释放锁
func ReleaseLock(lockKey string, lockValue string) error {
	// 使用Lua脚本删除锁，确保只有持有锁的客户端能够删除
	luaScript := `
        if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("DEL", KEYS[1])
        else
            return 0
        end
    `
	// 执行Lua脚本
	_, err := client.Eval(client.Context(), luaScript, []string{lockKey}, lockValue).Result()
	return err
}

func CacheLLen(key string) int64 {
	result, err := client.LLen(client.Context(), key).Result()
	if err != nil {
		log.Default().Println("error redis get lenght of list: ", err)
	}
	return result
}

func CacheRPop(key string) string {
	result, err := client.RPop(client.Context(), key).Result()
	if err != nil {
		log.Default().Println("error redis get lenght of list: ", err)
	}
	return result
}

func CacheLRange(key string, start int64, end int64) []string {
	result, err := client.LRange(client.Context(), key, start, end).Result()
	if err != nil {
		log.Default().Println("error redis get range of list:", err)
	}
	return result
}

func CacheDeleteByKey(keys ...string) int64 {
	effected, err := client.Del(client.Context(), keys...).Result()
	if err != nil {
		zap.L().Error("error delete by key: ", zap.Error(err))
	}
	return effected
}

func CacheGetKeysByPattern(pattern string) []string {
	result, err := client.Keys(client.Context(), pattern).Result()
	if err != nil {
		zap.L().Error("error query keys by pattern: ", zap.String("pattern", pattern), zap.Error(err))
	}
	return result
}

func SetIdentityAndIdMappingInCache(linux *Linux) {
	mins := 30 + common.GetRandomNumber(10)
	CacheSet(linux.LinuxId, strconv.FormatInt(linux.Id, 10), time.Duration(mins)*time.Minute)
}
