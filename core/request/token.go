package request

import (
	"net/http"
	"net/url"
	"time"
	"strconv"
	"fmt"
	"io"
	"bytes"
	"crypto/md5"
)

var userTokenMap map[string]string

func init() {
	userTokenMap = make(map[string]string)
}

func GenerateTokenKey(loginName string, r *http.Request) (tokenKey string) {
	var buffer bytes.Buffer
	buffer.WriteString(loginName)
	buffer.WriteString("_")
	buffer.WriteString(r.RemoteAddr)
	buffer.WriteString("_")
	buffer.WriteString(url.QueryEscape(r.RequestURI))
	return buffer.String()
}

func GenerateToken() (token string) {
	currentTime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(currentTime, 10))
	token = fmt.Sprintf("%x", h.Sum(nil))
	return
}

func GetToken(loginName string) (token string) {
	if len(loginName) == 0 {
		return
	}
	token, _ = userTokenMap[loginName]
	return
}

func SetToken(loginName string, token string) (bool) {
	if len(loginName) == 0 || len(token) == 0 {
		return false
	}
	userTokenMap[loginName] = token
	return true
}

func DeleteToken(loginName string) (result bool) {
	if len(loginName) == 0 {
		return false
	}
	delete(userTokenMap, loginName)
	return true
}