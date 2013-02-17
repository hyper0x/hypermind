package dao

import (
	"errors"
	"hypermind/core/base"
)

type User struct {
	LoginName   string
	Password    string
	Email       string
	MobilePhone string
	Rights      string
	Remark      string
}

func AddUserToDb(user *User) error {
	loginName := user.LoginName
	if user == nil || user.LoginName == "" {
		return errors.New("The parameter named user is NOT Ready! (loginName=\"\")")
	}
	userKey := getUserKey(loginName)
	userInfoMap := make(map[string]string, 6)
	userInfoMap[LOGIN_NAME_FIELD] = loginName
	userInfoMap[PASSWORD_FIELD] = encryptPassword(user.Password)
	userInfoMap[EMAIL_FIELD] = user.Email
	userInfoMap[MOBILE_PHONE_FIELD] = user.MobilePhone
	userInfoMap[RIGHTS_FIELD] = user.Rights
	userInfoMap[REMARK_FIELD] = user.Remark
	conn := redisPool.Get()
	defer conn.Close()
	err := SetHashBatch(userKey, userInfoMap)
	if err != nil {
		return err
	}
	return nil
}

func GetUserFromDb(loginName string) (*User, error) {
	if len(loginName) == 0 {
		return nil, errors.New("The parameter named loginName is EMPTY!")
	}
	userKey := getUserKey(loginName)
	userInfoMap, err := GetHashAll(userKey)
	if err != nil {
		return nil, err
	}
	user := new(User)
	user.LoginName = userInfoMap[LOGIN_NAME_FIELD]
	user.Password = userInfoMap[PASSWORD_FIELD]
	user.Email = userInfoMap[EMAIL_FIELD]
	user.MobilePhone = userInfoMap[MOBILE_PHONE_FIELD]
	user.Rights = userInfoMap[RIGHTS_FIELD]
	user.Remark = userInfoMap[REMARK_FIELD]
	return user, nil
}

func DeleteUserFromDb(loginName string) error {
	if len(loginName) == 0 {
		return errors.New("The parameter named loginName is EMPTY!")
	}
	userKey := getUserKey(loginName)
	err := DelKey(userKey)
	if err != nil {
		return err
	}
	return nil
}

func VerifyUser(loginName string, password string) (bool, error) {
	user, err := GetUserFromDb(loginName)
	if err != nil {
		return false, err
	}
	var pass bool = false
	if user != nil && user.Password == encryptPassword(password) {
		pass = true
	}
	return pass, nil
}

func getUserKey(loginName string) string {
	return USER_KEY_PREFIX + loginName
}

func encryptPassword(password string) (result string) {
	result = base.EncryptWithMd5(password)
	return
}
