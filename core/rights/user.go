package rights

import (
	"errors"
	"hypermind/core/base"
	"hypermind/core/dao"
)

type User struct {
	LoginName   string
	Password    string
	Email       string
	MobilePhone string
	Group       string
	Remark      string
}

func AddUser(user *User) error {
	if user == nil || user.LoginName == "" {
		return errors.New("The parameter named user is NOT Ready!")
	}
	loginName := user.LoginName
	userKey := getUserKey(loginName)
	userInfoMap := make(map[string]string)
	userInfoMap[LOGIN_NAME_FIELD] = loginName
	userInfoMap[PASSWORD_FIELD] = encryptPassword(user.Password)
	userInfoMap[EMAIL_FIELD] = user.Email
	userInfoMap[MOBILE_PHONE_FIELD] = user.MobilePhone
	userInfoMap[GROUP_FIELD] = user.Group
	userInfoMap[REMARK_FIELD] = user.Remark
	conn := dao.RedisPool.Get()
	defer conn.Close()
	_, err := dao.SetHashBatch(userKey, userInfoMap)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(loginName string) (*User, error) {
	if len(loginName) == 0 {
		return nil, errors.New("The parameter named loginName is EMPTY!")
	}
	userKey := getUserKey(loginName)
	userInfoMap, err := dao.GetHashAll(userKey)
	if err != nil {
		return nil, err
	}
	if len(userInfoMap) == 0 {
		return nil, nil
	}
	user := new(User)
	user.LoginName = userInfoMap[LOGIN_NAME_FIELD]
	user.Password = userInfoMap[PASSWORD_FIELD]
	user.Email = userInfoMap[EMAIL_FIELD]
	user.MobilePhone = userInfoMap[MOBILE_PHONE_FIELD]
	user.Group = userInfoMap[GROUP_FIELD]
	user.Remark = userInfoMap[REMARK_FIELD]
	return user, nil
}

func DeleteUser(loginName string) error {
	if len(loginName) == 0 {
		return errors.New("The parameter named loginName is EMPTY!")
	}
	userKey := getUserKey(loginName)
	_, err := dao.DelKey(userKey)
	if err != nil {
		return err
	}
	return nil
}

func VerifyUser(loginName string, password string) (bool, error) {
	user, err := GetUser(loginName)
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
	return dao.USER_KEY_PREFIX + loginName
}

func encryptPassword(password string) (result string) {
	result = base.EncryptWithMd5(password)
	return
}
