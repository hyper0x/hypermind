package myutil

import (
	"net/http"
	"time"
)

func SetCookie(w http.ResponseWriter,
               name string,
               value string,
               expirationMinute int) bool {
	expiration := time.Now()
	expirationNano := time.Minute * time.Duration(expirationMinute)
	expiration = expiration.Add(expirationNano)
	cookie := http.Cookie{Name: name, Value: value, Expires: expiration}
    http.SetCookie(w, &cookie)
    return true
}

func SetCookies(w http.ResponseWriter,
                data map[string]string,
                expirationMinute int) bool {
	expiration := time.Now()
	expirationNano := time.Minute * time.Duration(expirationMinute)
	expiration = expiration.Add(expirationNano)
	for k, v := range data {
		cookie := http.Cookie{Name: k, Value: v, Expires: expiration}
	    http.SetCookie(w, &cookie)
    }
    return true
}

func GetCookies(r *http.Request) (data map[string]string) {
	for _, cookie := range r.Cookies() {
		data[cookie.Name] = cookie.Value
	}
	return
}

func GetCookie(r *http.Request, cookieName string) (value string) {
	if len(cookieName) == 0 {
		return
	}
	for _, cookie := range r.Cookies() {
		if cookie.Name == cookieName {
			value = cookie.Value
		}
	}
	return
}

