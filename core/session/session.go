package session

import (
	"bytes"
	"errors"
	"fmt"
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
	// survivalSeconds: 
	//  '<0'- Don't set/Delete cookie; 
	//  '0' - Set temporary cookie; 
	//  '>0'- Set long term cookie according to this value & Set session expires
	survivalSeconds int,
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
	base.Logger().Infof("Initialize session (key=%s, grantors=%s)...\n", self.key, grantors)
	exists, err := dao.Exists(self.key)
	if err != nil {
		return err
	}
	if exists {
		_, err = dao.DelKey(self.key)
		if err != nil {
			return err
		}
	}
	_, err = dao.SetHash(self.key, SESSION_GRANTORS_KEY, grantors)
	if err != nil {
		return err
	}
	_, err = dao.SetHash(self.key, SESSION_SURVIVAL_SECONDS_KEY, strconv.FormatInt(int64(survivalSeconds), 10))
	if err != nil {
		return err
	}
	cookieMaxAge := survivalSeconds
	if survivalSeconds <= 0 {
		cookieMaxAge = -1
	}
	base.Logger().Infof("Set session cookie (value=%s, grantors=%s, maxAge=%d)...\n", self.sessionId, grantors, cookieMaxAge)
	result := hmSessionCookie.SetOne(self.w, SESSION_COOKIE_KEY, self.sessionId, cookieMaxAge)
	if result {
		base.Logger().Infof("Session cookie setting (value=%s, grantors=%s, maxAge=%d) is successful.\n", self.sessionId, grantors, cookieMaxAge)
	} else {
		base.Logger().Warnf("Session cookie setting (value=%s, grantors=%s, maxAge=%d) is failing!\n", self.sessionId, grantors, cookieMaxAge)
	}
	if survivalSeconds > 0 {
		done, err := dao.SetExpires(self.key, uint64(survivalSeconds))
		if err != nil || !done {
			warningMsg := fmt.Sprintf("Setting session expires failed! (key=%s, survivalSeconds=%d, error=%s)\n", self.key, survivalSeconds, err)
			base.Logger().Warnln(warningMsg)
		}
	}
	_, err = dao.SetHash(SESSION_MAP_KEY, grantors, self.key)
	if err != nil {
		return err
	}
	return nil
}

func (self *MySession) Destroy() (bool, error) {
	if len(self.key) == 0 || len(self.sessionId) == 0 {
		errorMsg := fmt.Sprintln("Uninitialized yet!")
		return false, errors.New(errorMsg)
	}
	base.Logger().Infof("Destroy session (key=%s)... \n", self.key)
	grantors, err := dao.GetHash(self.key, SESSION_GRANTORS_KEY)
	if err != nil {
		return false, err
	}
	_, err = dao.DelKey(self.key)
	if err != nil {
		return false, err
	}
	if len(grantors) > 0 {
		_, err = dao.DelHashField(SESSION_MAP_KEY, grantors)
		if err != nil {
			return false, err
		}
	}
	base.Logger().Infof("Delete session cookie (value=%s)...\n", self.sessionId)
	hmSessionCookie.Delete(SESSION_COOKIE_KEY, self.w)
	base.Logger().Infof("The session (key=%s) is destroyed. \n", self.key)
	return true, nil
}

func (self *MySession) Set(name string, value string) error {
	if len(name) == 0 {
		errorMsg := fmt.Sprintln("The parameter named name is EMPTY!")
		return errors.New(errorMsg)
	}
	_, err := dao.SetHash(self.key, name, value)
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
	_, err := dao.DelHashField(self.key, name)
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
		base.Logger().Warnln(warningMsg)
		return nil, nil
	}
	sessionkey := generateSessionKey(sessionId)
	exists, err := dao.Exists(sessionkey)
	if err != nil {
		return nil, err
	}
	if !exists {
		warningMsg := fmt.Sprintf("Not found matched session! No session in storage! (sessionId=%s, sessionKey=%s)", sessionId, sessionkey)
		base.Logger().Warnln(warningMsg)
		return nil, nil
	}
	grantors, err := dao.GetHash(sessionkey, SESSION_GRANTORS_KEY)
	if err != nil {
		return nil, err
	}
	if len(grantors) == 0 {
		warningMsg := fmt.Sprintf("Not found grantor from session (sessionKey=%s, attribute=%s)!\n", sessionkey, SESSION_GRANTORS_KEY)
		base.Logger().Warnln(warningMsg)
		return nil, nil
	}
	groupName, err := dao.GetHash(sessionkey, SESSION_GROUP_KEY)
	if err != nil {
		return nil, err
	}
	if len(groupName) == 0 {
		warningMsg := fmt.Sprintf("Not found group name from session (sessionKey=%s, attribute=%s)!\n", sessionkey, SESSION_GROUP_KEY)
		base.Logger().Warnln(warningMsg)
		return nil, err
	}
	servivalSecondsLiterals, err := dao.GetHash(sessionkey, SESSION_SURVIVAL_SECONDS_KEY)
	if err != nil {
		return nil, err
	}
	if len(servivalSecondsLiterals) == 0 {
		warningMsg := fmt.Sprintf("Not found session servival seconds. Use default value '0'. (sessionKey=%s, attribute=%s)!\n", sessionkey, SESSION_SURVIVAL_SECONDS_KEY)
		base.Logger().Warnln(warningMsg)
		return nil, err
	}
	hmSession := &MySession{key: sessionkey, sessionId: sessionId, w: w, r: r}
	return hmSession, nil
}

func newSession(grantors string, servivalSeconds int, w http.ResponseWriter, r *http.Request) (*MySession, error) {
	hmSession := &MySession{}
	err := hmSession.Initialize(grantors, servivalSeconds, w, r)
	if err != nil {
		return nil, err
	}
	user, err := rights.GetUser(grantors)
	if err != nil {
		return nil, err
	}
	base.Logger().Infof("Set session core attribute '%s' to value '%s' for grantors '%s'. (key=%s)\n", SESSION_GROUP_KEY, user.Group, grantors, hmSession.Key())
	err = hmSession.Set(SESSION_GROUP_KEY, user.Group)
	if err != nil {
		return nil, err
	}
	return hmSession, nil
}

func NewSession(grantors string, longTerm bool, w http.ResponseWriter, r *http.Request) (*MySession, error) {
	servivalSeconds := -1
	if longTerm {
		servivalSeconds = SESSION_SURVIVAL_SECONDS
	}
	return newSession(grantors, servivalSeconds, w, r)
}
