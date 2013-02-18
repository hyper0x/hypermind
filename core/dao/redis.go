package dao

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"go_lib"
	. "hypermind/core/base"
	"time"
)

var redisServerIp string
var redisServerPort string
var redisServerPassword string

var RedisPool *redis.Pool

func init() {
	config := GetHmConfig()
	err := config.ReadConfig(false)
	if err != nil {
		go_lib.LogErrorln("ConfigLoadError: ", err)
	}
	redisServerIp = config.Dict["redis_server_ip"]
	if len(redisServerIp) == 0 {
		redisServerIp = DEFAULT_REDIS_SERVER_IP
	}
	redisServerPort = config.Dict["redis_server_port"]
	if len(redisServerPort) == 0 {
		redisServerPort = DEFAULT_REDIS_SERVER_PORT
	}
	redisServerPassword = config.Dict["redis_server_password"]
	if len(redisServerPassword) == 0 {
		redisServerPassword = DEFAULT_REDIS_SERVER_PASSWORD
	}
	redisServerAddr := "127.0.0.1" + ":" + redisServerPort
	RedisPool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisServerAddr)
			if err != nil {
				return nil, err
			}
			if len(redisServerPassword) > 0 {
				if _, err := c.Do("AUTH", redisServerPassword); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
	}
}

func SetHash(key string, field string, value string) error {
	conn := RedisPool.Get()
	defer conn.Close()
	n, err := redis.Int(conn.Do("HSET", key, field, value))
	if err != nil {
		return err
	}
	if n != 0 && n != 1 {
		errorMsg := fmt.Sprintf("Redis operation failed! (cmd='HSET %v %v %v', n=%d)", key, field, value, n)
		return errors.New(errorMsg)
	}
	return nil
}

func SetHashBatch(key string, fieldMap map[string]string) error {
	conn := RedisPool.Get()
	defer conn.Close()
	for f, v := range fieldMap {
		err := SetHash(key, f, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetHash(key string, field string) (string, error) {
	if !Exists(key) {
		return "", nil
	}
	if !HashFieldExists(key, field) {
		return "", nil
	}
	conn := RedisPool.Get()
	defer conn.Close()
	value, err := redis.String(conn.Do("HGET", key, field))
	if err != nil {
		return "", err
	}
	return value, nil
}

func GetHashAll(key string) (map[string]string, error) {
	conn := RedisPool.Get()
	defer conn.Close()
	values, err := redis.Values(conn.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	length := len(values)
	for i := 0; i < length; i += 2 {
		key := fmt.Sprintf("%s", values[i])
		value := fmt.Sprintf("%s", values[i+1])
		result[key] = value
	}
	return result, nil
}

func DelKey(key string) error {
	if !Exists(key) {
		return nil
	}
	conn := RedisPool.Get()
	defer conn.Close()
	n, err := redis.Int(conn.Do("DEL", key))
	if err != nil {
		return err
	}
	if n != 0 && n != 1 {
		errorMsg := fmt.Sprintf("Redis operation failed! (cmd='DEL %v', n=%d)", key, n)
		return errors.New(errorMsg)
	}
	return nil
}

func DelHashField(key string, field string) error {
	if !Exists(key) {
		return nil
	}
	conn := RedisPool.Get()
	defer conn.Close()
	n, err := redis.Int(conn.Do("HDEL", key, field))
	if err != nil {
		return err
	}
	if n != 0 && n != 1 {
		errorMsg := fmt.Sprintf("Redis operation failed! (cmd='HDEL %v %v', n=%d)", key, field, n)
		return errors.New(errorMsg)
	}
	return nil
}

func Exists(key string) bool {
	conn := RedisPool.Get()
	defer conn.Close()
	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		go_lib.LogErrorf("JudgeKeyExistenceError (key=%s): %s\n ", key, err)
		return false
	}
	if !exists {
		go_lib.LogWarnf("The key '%s' is NONEXISTENCE.\n", key)
	}
	return exists
}

func HashFieldExists(key string, field string) bool {
	conn := RedisPool.Get()
	fieldExists, err := redis.Bool(conn.Do("HEXISTS", key, field))
	if err != nil {
		go_lib.LogErrorf("JudgeHashFieldExistenceError (key=%s, field=%s): %s\n ", key, field, err)
		return false
	}
	if !fieldExists {
		go_lib.LogWarnf("The field '%s' in hash key '%s' is NONEXISTENCE.\n", field, key)
		return false
	}
	return fieldExists
}
