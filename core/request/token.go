package request

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Token struct {
	Key   string
	Value string
}

var userTokenMap map[string]string

func init() {
	userTokenMap = make(map[string]string)
}

func GenerateToken(r *http.Request, extraInfo ...string) Token {
	var buffer bytes.Buffer
	buffer.WriteString(r.RequestURI)
	buffer.WriteString("|")
	buffer.WriteString(strings.Split(r.RemoteAddr, ":")[0])
	buffer.WriteString("|")
	buffer.WriteString(strconv.FormatInt(time.Now().UnixNano(), 10))
	for _, v := range extraInfo {
		buffer.WriteString("|")
		buffer.WriteString(v)
	}
	info := buffer.String()
	h := sha1.New()
	io.WriteString(h, info)
	tokenKey := fmt.Sprintf("%x", h.Sum(nil))
	userTokenMap[tokenKey] = info
	token := Token{Key: tokenKey, Value: info}
	return token
}

func SaveToken(token Token) bool {
	if len(token.Key) == 0 {
		return false
	}
	userTokenMap[token.Key] = token.Value
	return true
}

func CheckToken(tokenKey string) bool {
	if len(tokenKey) == 0 {
		return false
	}
	_, ok := userTokenMap[tokenKey]
	return ok
}

func RemoveToken(tokenKey string) bool {
	if len(tokenKey) == 0 {
		return false
	}
	delete(userTokenMap, tokenKey)
	return true
}
