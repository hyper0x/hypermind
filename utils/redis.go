package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"go_lib"
	"math/rand"
	"strconv"
	"strings"
	"time"
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
	n, err := redis.Int(conn.Do("HSET", USER_KEY, field, value))
	if err != nil {
		return err
	}
	if n != 0 && n != 1 {
		errorMsg := fmt.Sprintf("Redis operation failed! (cmd='HSET %v %v %v', n=%d)", USER_KEY, field, value, n)
		return errors.New(errorMsg)
	}
	return nil
}

func GetUserFromDb(loginName string) (*User, error) {
	if loginName == "" {
		return nil, errors.New("The parameter named loginName is EMPTY!")
	}
	if !exists(USER_KEY) {
		return nil, nil
	}
	if !hashFieldExists(USER_KEY, loginName) {
		return nil, nil
	}
	conn := redisPool.Get()
	defer conn.Close()
	literals, err := redis.String(conn.Do("HGET", USER_KEY, loginName))
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
	if !exists(USER_KEY) {
		return result, nil
	}
	conn := redisPool.Get()
	defer conn.Close()
	values, err := redis.Values(conn.Do("HGETALL", USER_KEY))
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

func DeleteUserFromDb(loginName string) error {
	if loginName == "" {
		return errors.New("The parameter named loginName is NULL!")
	}
	if !exists(USER_KEY) {
		return nil
	}
	conn := redisPool.Get()
	defer conn.Close()
	n, err := redis.Int(conn.Do("HDEL", USER_KEY, loginName))
	if err != nil {
		return err
	}
	if n != 0 && n != 1 {
		errorMsg := fmt.Sprintf("Redis operation failed! (cmd='HDEL %v %v', n=%d)", USER_KEY, loginName, n)
		return errors.New(errorMsg)
	}
	return nil
}

func VerifyAuthCode(authCode string) (bool, error) {
	if len(authCode) == 0 {
		return false, nil
	}
	currentAuthCode, err := GetCurrentAuthCode()
	var pass bool
	if err == nil {
		pass = (currentAuthCode == strings.TrimSpace(authCode))
	}
	if pass {
		go func() {
			var newAuthCode string
			for {
				newAuthCode = generateAuthCode()
				if newAuthCode != currentAuthCode {
					go_lib.LogInfof("New Auth Code: %s\n", newAuthCode)
					break
				}
			}
			if len(newAuthCode) > 0 {
				conn := redisPool.Get()
				defer conn.Close()
				err = pushAuthCode(newAuthCode, conn)
				if err != nil {
					go_lib.LogErrorf("New auth code pushing error: %s\n", err)
				}
			}
		}()
	}
	return pass, err
}

func GetCurrentAuthCode() (string, error) {
	conn := redisPool.Get()
	defer conn.Close()
	values, err := redis.Values(conn.Do("LRANGE", AUTH_CODE_KEY, 0, 0))
	if err != nil {
		return "", err
	}
	var currentAuthCode string
	if len(values) > 0 {
		var buffer bytes.Buffer
	        value := values[0]
	        valueBytes := value.([]byte)
	        for _, v := range valueBytes {
			buffer.WriteByte(v)
	        }
	        currentAuthCode = buffer.String()
	        go_lib.LogInfof("Current Code: %v\n", currentAuthCode)
	} else {
		initialAuthCode := generateInitialAuthCode()
		go_lib.LogInfof("Initial Auth Code: %s\n", initialAuthCode)
		err = pushAuthCode(initialAuthCode, conn)
		if err != nil {
			return "", err
		}
		currentAuthCode = initialAuthCode
	}
	return currentAuthCode, nil
}

func GetAndNewAuthCode() (string, error) {
	conn := redisPool.Get()
	defer conn.Close()
	newAuthCode := generateAuthCode()
	go_lib.LogInfof("New Auth Code: %s\n", newAuthCode)
	err := pushAuthCode(newAuthCode, conn)
	if err != nil {
		return "", err
	}
	return newAuthCode, nil
}

func pushAuthCode(code string, conn redis.Conn) error {
	n, err := redis.Int(conn.Do("LPUSH", AUTH_CODE_KEY, code))
	if err != nil {
		return err
	}
	if n < 0 {
		errorMsg := fmt.Sprintf("Redis operation failed! (cmd='LPUSH %v %v', n=%d)", AUTH_CODE_KEY, code, n)
		return errors.New(errorMsg)
	}
	return nil
}

func generateInitialAuthCode() string {
	var buffer bytes.Buffer
	now := time.Now()
	hour := fmt.Sprintf("%v", now.Hour())
	buffer.WriteString(hour)
	minute := fmt.Sprintf("%v", now.Minute())
	buffer.WriteString(minute)
	if buffer.Len() < 6 {
		for {
			infilling := fmt.Sprintf("%v", rand.Intn(99))
		        buffer.WriteString(infilling)
			if buffer.Len() >= 6 {
				break
			}
	        }
	}
	code := buffer.String()
	if len(code) > 6 {
		code = code[:6]
	}
	return code
}

func generateAuthCode() string {
	var limit int64 = 65535
	var buffer bytes.Buffer
	var temp string
	for {
		temp = strconv.FormatInt(rand.Int63n(limit), 16)
		buffer.WriteString(temp)
		if buffer.Len() >= 6 {
			break
		}
	}
	code := buffer.String()
	if len(code) > 6 {
		code = code[:6]
	}
	return code
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
	fieldExists, err := redis.Bool(conn.Do("HEXISTS", USER_KEY, field))
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
