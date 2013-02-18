package session

import (
	"bytes"
	"errors"
	"fmt"
	"go_lib"
	"hypermind/core/base"
	"hypermind/core/dao"
	"hypermind/core/rights"
	"net/http"
	"strconv"
)

var hmSessionCookie *MyCookie

func init() {
	hmSessionCookie = &MyCookie{key: COOKIE_KEY_PREFIX}
}

type MySession struct {
	key       string
	sessionId string
	w         http.ResponseWriter
	r         *http.Request
}

func (self *MySession) Initialize(
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
	go_lib.LogInfof("Initialize session (key=%s, grantors=%s)...\n", self.key, grantors)
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

func (self *MySession) Destroy() (bool, error) {
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

func (self *MySession) Set(name string, value string) error {
	if len(name) == 0 {
		errorMsg := fmt.Sprintln("The parameter named name is EMPTY!")
		return errors.New(errorMsg)
	}
	err := dao.SetHash(self.key, name, value)
	return err
}

func (self *MySession) SetBatch(contentMap map[string]string) error {
	if len(contentMap) == 0 {
		return nil
	}
	for k, v := range contentMap {
		err := self.Set(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *MySession) Get(name string) (string, error) {
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

func (self *MySession) GetAll() (map[string]string, error) {
	sessionMap, err := dao.GetHashAll(self.sessionId)
	if err != nil {
		return nil, err
	}
	return sessionMap, nil
}

func (self *MySession) Delete(name string) error {
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

func (self *MySession) Key() string {
	return self.key
}

func (self *MySession) SessionID() string {
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

func GetMatchedSession(w http.ResponseWriter, r *http.Request) (*MySession, error) {
	sessionId := hmSessionCookie.GetOne(SESSION_COOKIE_KEY, r)
	if len(sessionId) == 0 {
		warningMsg := fmt.Sprintf("Not found matched session! No session cookie!")
		go_lib.LogWarnln(warningMsg)
		return nil, nil
	}
	sessionkey := generateSessionKey(sessionId)
	if !dao.Exists(sessionkey) {
		warningMsg := fmt.Sprintf("Not found matched session! No session in storage! (sessionId=%s, sessionKey=%s)", sessionId, sessionkey)
		go_lib.LogWarnln(warningMsg)
		return nil, nil
	}
	grantors, err := dao.GetHash(sessionkey, SESSION_GRANTORS_KEY)
	if err != nil {
		return nil, err
	}
	if len(grantors) == 0 {
		warningMsg := fmt.Sprintf("Not found grantor from session (sessionKey=%s, field=%s)!\n", sessionkey, SESSION_GRANTORS_KEY)
		go_lib.LogWarnln(warningMsg)
		return nil, nil
	}
	signLiterals, err := dao.GetHash(sessionkey, SESSION_COOKIE_SIGN_KEY)
	if err != nil {
		return nil, err
	}
	var cookieSign int64
	if len(signLiterals) == 0 {
		warningMsg := fmt.Sprintf("Not found session cookie sign from session. Use default value '0'. (sessionKey=%s, field=%s)!\n", sessionkey, SESSION_COOKIE_SIGN_KEY)
		go_lib.LogWarnln(warningMsg)
		cookieSign = 0
	} else {
		cookieSign, err = strconv.ParseInt(signLiterals, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	hmSession := &MySession{}
	err = hmSession.Initialize(grantors, int(cookieSign), w, r)
	if err != nil {
		return nil, err
	}
	return hmSession, nil
}

func NewSession(grantors string, longTerm bool, w http.ResponseWriter, r *http.Request) (*MySession, error) {
	hmSession := &MySession{}
	cookieSign := 0
	if longTerm {
		cookieSign = 1
	}
	err := hmSession.Initialize(grantors, cookieSign, w, r)
	if err != nil {
		return nil, err
	}

	user, err := rights.GetUser(grantors)
	if err != nil {
		return nil, err
	}
	contentMap := make(map[string]string)
	contentMap[SESSION_GRANTORS_KEY] = grantors
	contentMap[SESSION_COOKIE_SIGN_KEY] = strconv.FormatInt(int64(cookieSign), 10)
	contentMap[SESSION_GROUP_KEY] = user.Group
	err = hmSession.SetBatch(contentMap)
	if err != nil {
		return nil, err
	}
	return hmSession, nil
}
