package utils

import (
	"net/http"
)

func GetStagedUserInfo(w http.ResponseWriter, r *http.Request) map[string]string {
	userInfoMap := make(map[string]string)
	loginName := GetOneCookie(r, LoginNameKey)
	if len(loginName) == 0 {
		LogErrorln("The login name is NOT in the cookie of client!")
		session := GetSession(w, r)
		if v := session.Get(LoginNameKey);v != nil {
			loginName = v.(string)
		}
	}
	userInfoMap[LoginNameKey] = loginName
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

