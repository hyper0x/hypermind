package request

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"hypermind/core/base"
	"hypermind/core/dao"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type NewAuthCodeTrigger func(string)

var newAuthCodeTriggerMap map[string]NewAuthCodeTrigger

func init() {
	newAuthCodeTriggerMap = make(map[string]NewAuthCodeTrigger)
	firstNewAuthCodeTrigger := func(newAuthCode string) {
		base.Logger().Infof("There has a new auth code '%s'.", newAuthCode)
	}
	AddNewAuthCodeTrigger("monitoring", firstNewAuthCodeTrigger)
}

func AddNewAuthCodeTrigger(id string, trigger NewAuthCodeTrigger) {
	newAuthCodeTriggerMap[id] = trigger
}

func HasNewAuthCodeTrigger(id string) bool {
	_, ok := newAuthCodeTriggerMap[id]
	return ok
}

func DelNewAuthCodeTrigger(id string) {
	delete(newAuthCodeTriggerMap, id)
}

func activateNewAuthCodeTriggers(newAuthCode string) {
	for _, f := range newAuthCodeTriggerMap {
		f(newAuthCode)
	}
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
					break
				}
			}
			if len(newAuthCode) > 0 {
				conn := dao.RedisPool.Get()
				defer conn.Close()
				err = pushAuthCode(newAuthCode, conn)
				if err != nil {
					base.Logger().Errorf("New auth code pushing error: %s\n", err)
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
		base.Logger().Infof("Current Code: '%s'\n", currentAuthCode)
	} else {
		initialAuthCode := generateInitialAuthCode()
		base.Logger().Infof("Initial Auth Code: '%s'\n", initialAuthCode)
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
	if n <= 0 {
		errorMsg := fmt.Sprintf("Redis operation failed! (cmd='LPUSH %v %v', n=%d)", dao.AUTH_CODE_KEY, code, n)
		return errors.New(errorMsg)
	} else {
		activateNewAuthCodeTriggers(code)
	}
	return nil
}

func generateInitialAuthCode() string {
	var buffer bytes.Buffer
	ns := fmt.Sprintf("%v", time.Now().Nanosecond())
	buffer.WriteString(ns)
	if buffer.Len() < AUTH_CODE_LENGTH {
		for {
			infilling := fmt.Sprintf("%v", rand.Intn(99))
			buffer.WriteString(infilling)
			if buffer.Len() >= AUTH_CODE_LENGTH {
				break
			}
		}
	}
	code := buffer.String()
	if len(code) > AUTH_CODE_LENGTH {
		code = code[:AUTH_CODE_LENGTH]
	}
	return code
}

func generateAuthCode() string {
	var limit int64 = 255
	var buffer bytes.Buffer
	var temp string
	for {
		temp = strconv.FormatInt(rand.Int63n(limit), 16)
		buffer.WriteString(temp)
		if buffer.Len() >= AUTH_CODE_LENGTH {
			break
		}
	}
	code := buffer.String()
	if len(code) > AUTH_CODE_LENGTH {
		code = code[:AUTH_CODE_LENGTH]
	}
	return code
}
