package redisUtil

import (
	"XcxcPan/common/models"
	"context"
	"encoding/json"
	"strconv"
	"time"
)

func Set(key string, value interface{}, duration time.Duration) {
	valueJson, _ := json.Marshal(value)
	models.RDb.Set(context.Background(), key, valueJson, duration)
}
func SetWithLogicalExpire(key string, value interface{}, duration time.Duration) {
	cacheValue := struct {
		Value      interface{}
		ExpireTime int64
	}{
		Value:      value,
		ExpireTime: time.Now().Add(duration).Unix(),
	}
	cacheValueJson, _ := json.Marshal(cacheValue)
	models.RDb.Set(context.Background(), key, cacheValueJson, 0)
}

func Get(key string, value any) {
	valueJson, _ := models.RDb.Get(context.Background(), key).Result()
	json.Unmarshal([]byte(valueJson), &value)
}

func SetAdd(key string, value any) {
	models.RDb.SAdd(context.Background(), key, value)
}
func GetSet(key string) []string {
	return models.RDb.SMembers(context.Background(), key).Val()
}

func SetHash(key string, field int, value int) {
	fieldStr := strconv.Itoa(field)
	valueStr := strconv.Itoa(value)
	models.RDb.HSet(context.Background(), key, fieldStr, valueStr)
}

func GetHashInt(key string) map[int]int {
	hash, _ := models.RDb.HGetAll(context.Background(), key).Result()
	// 转换为 map[int]int
	intMap := make(map[int]int)
	for k, v := range hash {
		// 将 key 转换为 int
		intKey, err := strconv.Atoi(k)
		if err != nil {
			continue
		}

		// 将 value 转换为 int
		intValue, err := strconv.Atoi(v)
		if err != nil {
			continue
		}

		// 存入 map[int]int
		intMap[intKey] = intValue
	}
	return intMap
}
