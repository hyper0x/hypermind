package session

import (
	"bytes"
	"errors"
	"fmt"
	"go_lib"
	"hypermind/core/base"
	"hypermind/core/dao"
	"net/http"
	"strconv"
)

var hmSessionCookie *HmCookie

func init() {
	hmSessionCookie = &HmCookie{key: COOKIE_KEY_PREFIX}
}

type HmSession struct {
	key       string
	sessionId string
	w         http.ResponseWriter
	r         *http.Request
}

func (self *HmSession) Initialize(
	grantors string,
	cookieSign int, // cookieSign: '<0'-Don't set cookie; '0'-Set temporary cookie, '>0'- Set long term cookie
	w http.ResponseWriter,
	r *http.Request) error {
	if len(grantors) == 0 {
		errorMsg := fmt.Sprintln("The session grantors is EMPTY!")
		return errors.New(errorMsg)
	}
	if w == nil {
		errorMsg := fmt.Sprintln("The pointer of http response writer is NIL!")
		return errors.New(errorMsg)
	}
	if r == nil {
		errorMsg := fmt.Sprintln("The pointer of http request is NIL!")
		return errors.New(errorMsg)
	}
	self.w = w
	self.r = r
	self.sessionId = generateSessionId(grantors, r)
	self.key = generateSessionKey(self.sessionId)
	go_lib.LogInfof("Initialize session (key=%s)...\n", self.key)
	err := dao.SetHash(self.key, SESSION_GRANTORS_KEY, grantors)
	if err != nil {
		return err
	}
	err = dao.SetHash(self.key, SESSION_COOKIE_SIGN_KEY, strconv.FormatInt(int64(cookieSign), 10))
	if err != nil {
		return err
	}
	cookieMaxAge := -1
	if cookieSign > 0 {
		cookieMaxAge = COOKIE_MAX_AGE
	} else if cookieSign == 0 {
		cookieMaxAge = -1
	} else {
		cookieMaxAge = 0
	}
	go_lib.LogInfof("Set session cookie (value=%s, grantors=%s, maxAge=%d)...\n", self.sessionId, grantors, cookieMaxAge)
	result := hmSessionCookie.SetOne(self.w, SESSION_COOKIE_KEY, self.sessionId, cookieMaxAge)
	if result {
		go_lib.LogInfof("Session cookie Setting (value=%s, grantors=%s, maxAge=%d) is successful.\n", self.sessionId, grantors, cookieMaxAge)
	} else {
		go_lib.LogWarnf("Session cookie Setting (value=%s, grantors=%s, maxAge=%d) is failing!\n", self.sessionId, grantors, cookieMaxAge)
	}
	return nil
}

func (self *HmSession) Destroy() (bool, error) {
	if len(self.key) == 0 || len(self.sessionId) == 0 {
		errorMsg := fmt.Sprintln("Uninitialized yet!")
		return false, errors.New(errorMsg)
	}
	go_lib.LogInfof("Destroy session (key=%s)...\n", self.key)
	err := dao.DelKey(self.key)
	if err != nil {
		return false, err
	}
	go_lib.LogInfof("Delete session cookie (value=%s)...\n", self.sessionId)
	hmSessionCookie.Delete(SESSION_COOKIE_KEY, self.w)
	return true, nil
}

func (self *HmSession) Set(name string, value string) error {
	if len(name) == 0 {
		errorMsg := fmt.Sprintln("The parameter named name is EMPTY!")
		return errors.New(errorMsg)
	}
	err := dao.SetHash(self.key, name, value)
	return err
}

func (self *HmSession) Get(name string) (string, error) {
	if len(name) == 0 {
		errorMsg := fmt.Sprintln("The parameter named name is EMPTY!")
		return "", errors.New(errorMsg)
	}
	value, err := dao.GetHash(self.key, name)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (self *HmSession) GetAll() (map[string]string, error) {
	sessionMap, err := dao.GetHashAll(self.sessionId)
	if err != nil {
		return nil, err
	}
	return sessionMap, nil
}

func (self *HmSession) Delete(name string) error {
	if len(name) == 0 {
		errorMsg := fmt.Sprintln("The parameter named name is EMPTY!")
		return errors.New(errorMsg)
	}
	err := dao.DelHashField(self.key, name)
	if err != nil {
		return err
	}
	return nil
}

func (self *HmSession) Key() string {
	return self.key
}

func (self *HmSession) SessionID() string {
	return self.sessionId
}

func generateSessionId(key string, r *http.Request) (tokenKey string) {
	var buffer bytes.Buffer
	buffer.WriteString(key)
	buffer.WriteString("_")
	buffer.WriteString(r.RemoteAddr)
	buffer.WriteString("_[")
	buffer.WriteString(r.UserAgent())
	buffer.WriteString("]")
	return base.EncryptWithSha1(buffer.String())
}

func generateSessionKey(name string) string {
	return SESSION_KEY_PREFIX + name
}

func GetMatchedSession(w http.ResponseWriter, r *http.Request) (*HmSession, error) {
	sessionId := hmSessionCookie.GetOne(SESSION_COOKIE_KEY, r)
	if len(sessionId) == 0 {
		return nil, errors.New("Not found matched session! (no session cookie)")
	}
	sessionkey := generateSessionKey(sessionId)
	if !dao.Exists(sessionkey) {
		return nil, errors.New("Not found matched session! (no session in storage)")
	}
	grantors, err := dao.GetHash(sessionkey, SESSION_GRANTORS_KEY)
	if err != nil {
		return nil, err
	}
	if len(grantors) == 0 {
		errorMsg := fmt.Sprintf("No found grantor from session (key=%s, field=%s)!\n", sessionkey, SESSION_GRANTORS_KEY)
		return nil, errors.New(errorMsg)
	}
	signLiterals, err := dao.GetHash(sessionkey, SESSION_COOKIE_SIGN_KEY)
	if err != nil {
		return nil, err
	}
	sign, err := strconv.ParseInt(signLiterals, 10, 64)
	if err != nil {
		return nil, err
	}
	hmSession := &HmSession{}
	err = hmSession.Initialize(grantors, int(sign), w, r)
	if err != nil {
		return nil, err
	}
	return hmSession, nil
}

func NewSession(grantors string, longTerm bool, w http.ResponseWriter, r *http.Request) (*HmSession, error) {
	hmSession := &HmSession{}
	cookieSign := 0
	if longTerm {
		cookieSign = 1
	}
	err := hmSession.Initialize(grantors, cookieSign, w, r)
	if err != nil {
		return nil, err
	}
	return hmSession, nil
}
