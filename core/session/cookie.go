package session

import (
	"net/http"
	"strings"
)

type HmCookie struct {
	key string
}

func (self *HmCookie) SetOne(
	w http.ResponseWriter,
	name string,
	value string,
	maxAge int) bool {
	if maxAge <= 0 {
		maxAge = -1
	}
	fullName := self.key + name
	cookie := http.Cookie{Name: fullName, Value: value, MaxAge: maxAge}
	http.SetCookie(w, &cookie)
	return true
}

func (self *HmCookie) Set(
	w http.ResponseWriter,
	data map[string]string,
	maxAge int) bool {
	if maxAge <= 0 {
		maxAge = -1
	}
	for k, v := range data {
		fullName := self.key + k
		cookie := http.Cookie{Name: fullName, Value: v, Path: "/", HttpOnly: true, MaxAge: maxAge}
		http.SetCookie(w, &cookie)
	}
	return true
}

func (self *HmCookie) GetOne(name string, r *http.Request) (value string) {
	if len(name) == 0 {
		return
	}
	keyLength := len(self.key)
	for _, cookie := range r.Cookies() {
		if !strings.HasPrefix(cookie.Name, self.key) {
			continue
		}
		cookieName := cookie.Name[keyLength:]
		if cookieName == name {
			value = cookie.Value
		}
	}
	return
}

func (self *HmCookie) Get(r *http.Request) map[string]string {
	keyLength := len(self.key)
	cookieMap := make(map[string]string)
	for _, cookie := range r.Cookies() {
		if !strings.HasPrefix(cookie.Name, self.key) {
			continue
		}
		name := cookie.Name[keyLength:]
		cookieMap[name] = cookie.Value
	}
	return cookieMap
}

func (self *HmCookie) Delete(cookieName string, w http.ResponseWriter) bool {
	if len(cookieName) == 0 {
		return false
	}
	fullName := self.key + cookieName
	cookie := http.Cookie{Name: fullName, Path: "/", HttpOnly: true, MaxAge: 0}
	http.SetCookie(w, &cookie)
	return true
}
