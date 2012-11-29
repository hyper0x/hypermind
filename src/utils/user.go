package utils

import (
	"net/http"
)

type User struct {
	LoginName string
	Password string
	Email string
	MobilePhone string `json:"omitempty"`
	Group int
	Privileges []Privilege
	Remark string
}

type Privilege struct {
	Name string
	Tag int
}

var defaultPrivileges []Privilege
var rootPrivileges []Privilege

func init() {
	defaultPrivileges = append(defaultPrivileges, Privilege{Name: "meeting-kanban", Tag: 1})
	rootPrivileges = append(rootPrivileges, Privilege{Name: "meeting-kanban", Tag: 1})
	rootPrivileges = append(rootPrivileges, Privilege{Name: "statistics-kanban", Tag: 1})
	rootPrivileges = append(rootPrivileges, Privilege{Name: "user-kanban", Tag: 1})
	user, err := GetUserFromDb(ROOT_USER_NAME)
	if err != nil {
		LogErrorf("RootInitializationError:%s\n", err)
	} else {
		if user == nil {
			LogInfo("Initialize root user...")
			root := User{
				LoginName: ROOT_USER_NAME,
				Password: "hypermind",
				Email: "freej.cn@gmail.com",
				Group: 0,
				Privileges: rootPrivileges,
				Remark: "root user"}
			err = SetUserToDb(root)
			if err != nil {
				LogErrorf("RootInitError:%s\n", err)
			} else {
				LogInfo("done\n")
			}
		}
	}
}

func GetStagedUserInfo(w http.ResponseWriter, r *http.Request) map[string]string {
	userInfoMap := make(map[string]string)
	loginName := GetOneCookie(r, LOGIN_NAME_KEY)
	if len(loginName) == 0 {
		LogWarnln("The login name is NOT in the cookie of client!")
		session := GetSession(w, r)
		if v := session.Get(LOGIN_NAME_KEY);v != nil {
			loginName = v.(string)
		}
	}
	userInfoMap[LOGIN_NAME_KEY] = loginName
	return userInfoMap
}

func SetUserInfoToStage(
        userInfoMap map[string]string,
		w http.ResponseWriter,
		r *http.Request,
		onlySession bool) bool {
	if len(userInfoMap) == 0 {
		return false
	}
	if !onlySession {
		SetCookies(w, userInfoMap, CookieLifecycleMinutes)
	}
	session := GetSession(w, r)
	for key, value := range userInfoMap {
		session.Set(key, value)
	}
	return true
}

func RemoveUserInfoFromStage(userInfoMap map[string]string, w http.ResponseWriter, r *http.Request) bool {
	if len(userInfoMap) == 0 {
		return false
	}
	SetCookies(w, userInfoMap, CookieLifecycleMinutes)
	session := GetSession(w, r)
	for key, _ := range userInfoMap {
		session.Delete(key)
		DeleteCookie(w, key)
	}
	return true
}

func VerifyUser(loginName string, password string) (bool, error) {
	user, err := GetUserFromDb(loginName)
	if err != nil {
		return false, err
	}
	var pass bool = false
	if user != nil && user.Password == password {
		pass = true
	}
	return pass, nil
}
