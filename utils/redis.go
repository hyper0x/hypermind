package utils

import (
	"github.com/garyburd/redigo/redis"
	"errors"
	"time"
	"go_lib"
)

var redisServerIp string
var redisServerPort string
var redisServerPassword string
var redisPool *redis.Pool

func init() {
	err := myConfig.ReadConfig(false)
	if err != nil {
		go_lib.LogErrorln("ConfigLoadError: ", err)
	}
	redisServerIp = myConfig.Dict["redis_server_ip"]
	if len(redisServerIp) == 0 {
		redisServerIp = DEFAULT_REDIS_SERVER_IP
	}
	redisServerPort = myConfig.Dict["redis_server_port"]
	if len(redisServerPort) == 0 {
		redisServerPort = DEFAULT_REDIS_SERVER_PORT
	}
	redisServerPassword = myConfig.Dict["redis_server_password"]
	if len(redisServerPassword) == 0 {
		redisServerPassword = DEFAULT_REDIS_SERVER_PASSWORD
	}
	redisServerAddr := "127.0.0.1" + ":" + redisServerPort
	redisPool = &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func () (redis.Conn, error) {
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

func SetUserToDb(user User) error {
	if user.LoginName == "" {
		return errors.New("The parameter named user is NOT Ready! (loginName=\"\")")
	}
	field := user.LoginName
	value, err := MarshalUser(user)
	if err != nil {
		return err
	}
	conn := redisPool.Get()
	defer conn.Close()
	n, err := redis.Int(conn.Do("HSET", REDIS_USER_KEY, field, value))
	if err != nil {
		return err
	}
	if n !=0 && n != 1 {
		return errors.New("The INVALID result: " + string(n))
	}
	return nil
}

func GetUserFromDb(loginName string) (*User, error) {
	if loginName == "" {
		return nil, errors.New("The parameter named loginName is EMPTY!")
	}
	if !exists(REDIS_USER_KEY) {
		return nil, nil
	}
	if !hashFieldExists(REDIS_USER_KEY, loginName) {
		return nil, nil
	}
	conn := redisPool.Get()
	defer conn.Close()
	literals, err := redis.String(conn.Do("HGET", REDIS_USER_KEY, loginName))
	if err != nil {
		return nil, err
	}
	if len(literals) > 0 {
		user, err := UnmarshalUser(literals)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, nil
}

func GetAllUsersFromDb() (map[string]*User, error) {
	var tempContainer map[string]string
	var result map[string]*User
	if !exists(REDIS_USER_KEY) {
	   return result, nil
	}
	conn := redisPool.Get()
	defer conn.Close()
	values, err := redis.Values(conn.Do("HGETALL", REDIS_USER_KEY))
	if err != nil {
		return nil, err
	}
	err = redis.ScanStruct(values, tempContainer)
	if err != nil {
		return nil, err
	}
	var tempUser *User
	for k, v := range tempContainer {
		*tempUser, err = UnmarshalUser(string(v))
		if err != nil {
			go_lib.LogErrorf("UnmarshalUserError (json=%s): %s\n", v, err)
		} else {
			result[k] = tempUser
		}
	}
	return result, nil
}

func DeleteUserFromDb(loginName string) (error) {
	if loginName == "" {
		return errors.New("The parameter named loginName is NULL!")
	}
	if !exists(REDIS_USER_KEY) {
		return nil
	}
	conn := redisPool.Get()
	defer conn.Close()
	n, err := redis.Int(conn.Do("HDEL", REDIS_USER_KEY, loginName))
	if err != nil {
		return err
	}
	if n !=0 && n != 1 {
		return errors.New("The INVALID result: " + string(n))
	}
	return nil
}

func exists(key string) bool {
	conn := redisPool.Get()
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

func hashFieldExists(key string, field string) bool {
	conn := redisPool.Get()
	fieldExists, err := redis.Bool(conn.Do("HEXISTS", REDIS_USER_KEY, field))
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

