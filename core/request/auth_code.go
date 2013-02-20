package request

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"go_lib"
	"hypermind/core/dao"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

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
				conn := dao.RedisPool.Get()
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
	conn := dao.RedisPool.Get()
	defer conn.Close()
	values, err := redis.Values(conn.Do("LRANGE", dao.AUTH_CODE_KEY, 0, 0))
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
	conn := dao.RedisPool.Get()
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
	n, err := redis.Int(conn.Do("LPUSH", dao.AUTH_CODE_KEY, code))
	if err != nil {
		return err
	}
	if n < 0 {
		errorMsg := fmt.Sprintf("Redis operation failed! (cmd='LPUSH %v %v', n=%d)", dao.AUTH_CODE_KEY, code, n)
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
