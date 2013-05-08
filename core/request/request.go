package request

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hypermind/core/base"
	"hypermind/core/rights"
	"hypermind/core/session"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

type RequestInfo struct {
	Form   url.Values
	Method string
	Path   string
	Scheme string
}

var pageParameterMap map[string]string = map[string]string{
	base.HOME_PAGE_KEY:           base.HOME_PAGE,
	base.ABOUT_ME_PAGE_KEY:       base.ABOUT_ME_PAGE,
	base.ABOUT_WEBSITE_PAGE_KEY:  base.ABOUT_WEBSITE_PAGE,
	base.MEETING_KANBAN_PAGE_KEY: base.MEETING_KANBAN_PAGE,
	base.PROJECT_HASH_RING_KEY:   base.PROJECT_HASH_RING_PAGE,
	base.ADMIN_AUTH_CODE_KEY:     base.ADMIN_AUTH_CODE_PAGE,
	base.ADMIN_USER_LIST_KEY:     base.ADMIN_USER_LIST_PAGE,
}

func GeneratePagePath(reqPage string) string {
	page := "home"
	if len(reqPage) > 0 {
		page = reqPage
	}
	return "web/page/" + page + ".gtpl"
}

func GetRequestInfo(r *http.Request) string {
	requestInfo := RequestInfo{Form: r.Form, Method: r.Method, Path: r.URL.Path, Scheme: r.URL.Scheme}
	b, err := json.Marshal(requestInfo)
	if err != nil {
		base.Logger().Errorln("JsonMarshalError:", err)
	}
	return string(b)
}

func GenerateBasicAttrMap(w http.ResponseWriter, r *http.Request) map[string]string {
	attrMap := make(map[string]string)
	host, port := splitHostPort(r.Host)
	attrMap["serverAddr"] = host
	attrMap["serverPort"] = port
	for pageKey, page := range pageParameterMap {
		attrMap[pageKey] = page
	}
	hmSession, err := GetSession(w, r)
	if err != nil {
		base.Logger().Errorln("GetSessionError: %s\n", err)
	} else {
		var pageRights map[string]string
		var loginName string
		var groupName string
		if hmSession != nil {
			pageRights = getPageRights(hmSession)
			loginName = getLoginName(hmSession)
			groupName = getGroupName(hmSession)
		} else {
			pageRights = rights.GetGuestPageRights()
			loginName = ""
			groupName = ""
		}
		for p, pr := range pageRights {
			attrMap[p] = pr
		}
		attrMap[LOGIN_NAME_KEY] = loginName
		attrMap[GROUP_NAME_KEY] = groupName
	}
	return attrMap
}

func GetSession(w http.ResponseWriter, r *http.Request) (*session.MySession, error) {
	hmSession, err := session.GetMatchedSession(w, r)
	if err != nil {
		return nil, err
	}
	return hmSession, nil
}

func getLoginName(hmSession *session.MySession) string {
	if hmSession == nil {
		return ""
	}
	grantors, err := hmSession.Get(session.SESSION_GRANTORS_KEY)
	if err != nil {
		base.Logger().Errorln("SessionGetError (field=%s): %s\n", session.SESSION_GRANTORS_KEY, err)
		return ""
	}
	return grantors
}

func getGroupName(hmSession *session.MySession) string {
	if hmSession == nil {
		return ""
	}
	groupName, err := hmSession.Get(session.SESSION_GROUP_KEY)
	if err != nil {
		base.Logger().Errorln("SessionGetError (field=%s): %s\n", session.SESSION_GROUP_KEY, err)
		return ""
	}
	return groupName
}

func getPageRights(hmSession *session.MySession) map[string]string {
	var pageRights map[string]string
	if hmSession == nil {
		return pageRights
	}
	groupName, err := hmSession.Get(session.SESSION_GROUP_KEY)
	if err != nil {
		base.Logger().Errorln("SessionGetError (field=%s): %s\n", session.SESSION_GROUP_KEY, err)
		return pageRights
	}
	userGroup, err := rights.GetUserGroup(groupName)
	if err != nil {
		base.Logger().Errorln("GetUserGroupError (groupName=%s): %s\n", groupName, err)
		return pageRights
	}
	if userGroup != nil {
		pageRights = userGroup.Rights.PageRights
	}
	return pageRights
}

func GetSessionMap(w http.ResponseWriter, r *http.Request) map[string]string {
	var sessionMap map[string]string
	hmSession, err := GetSession(w, r)
	if err != nil {
		base.Logger().Errorln("GetSessionError: %s\n", err)
		return sessionMap
	}
	if hmSession != nil {
		sessionMap, err = hmSession.GetAll()
		if err != nil {
			base.Logger().Errorln("SessionGetAllError: %s\n", err)
			return sessionMap
		}
	}
	return sessionMap
}

func splitHostPort(requestHost string) (host string, port string) {
	if splitIndex := strings.Index(requestHost, ":"); splitIndex > 0 {
		host = requestHost[0:splitIndex]
		port = requestHost[splitIndex+1 : len(requestHost)]
	} else {
		host = requestHost
		config := base.GetHmConfig()
		err := config.ReadConfig(false)
		if err != nil {
			base.Logger().Errorln("ConfigLoadError: ", err)
			port = "80"
		} else {
			port = fmt.Sprintf("%v", config.Dict["server_port"])
		}
	}
	return
}

func VerifyRegisterForm(r *http.Request) (fieldMap map[string]string, invalidFields []string) {
	fieldMap = make(map[string]string)
	invalidFields = make([]string, 1)
	loginName := r.FormValue(LOGIN_NAME_KEY)
	fieldMap[LOGIN_NAME_KEY] = loginName
	if m, _ := regexp.MatchString("^[a-zA-Z-_\\.]{1, 10}$", loginName); !m {
		invalidFields = append(invalidFields, LOGIN_NAME_KEY)
	}
	password := r.FormValue(PASSWORD_KEY)
	fieldMap[PASSWORD_KEY] = password
	if m, _ := regexp.MatchString("^[a-zA-Z-_\\.]{1, 20}$", password); !m {
		invalidFields = append(invalidFields, PASSWORD_KEY)
	}
	cnName := r.FormValue(CN_NAME_KEY)
	fieldMap[CN_NAME_KEY] = cnName
	if len(cnName) > 0 {
		if m, _ := regexp.MatchString("^[\\x{4e00}-\\x{9fa5}]+$", cnName); !m {
			invalidFields = append(invalidFields, CN_NAME_KEY)
		}
	}
	email := r.FormValue(EMAIL_KEY)
	fieldMap[EMAIL_KEY] = email
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, email); !m {
		invalidFields = append(invalidFields, EMAIL_KEY)
	}
	mobilePhone := r.FormValue(MOBILE_PHONE_KEY)
	fieldMap[MOBILE_PHONE_KEY] = mobilePhone
	if m, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, mobilePhone); !m {
		invalidFields = append(invalidFields, MOBILE_PHONE_KEY)
	}
	return
}

func DeleteTempFile(delay time.Duration, filePath string) (err error) {
	time.Sleep(delay)
	err = os.Remove(filePath)
	if err != nil {
		base.Logger().Errorf("Occur error when delete file '%s': %s\n", filePath, err)
	} else {
		base.Logger().Infof("The file '%s' is deleted.\n", filePath, err)
	}
	return
}

func EncodePassport(originalPasspord string) (result string) {
	h := md5.New()
	io.WriteString(h, originalPasspord)
	result = fmt.Sprintf("%x", h.Sum(nil))
	return
}

func UrlEncoding(s string) string {
	var buf bytes.Buffer
	var encoder = base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write([]byte(s))
	encoder.Close()
	return buf.String()
}

func UrlDecoding(s string) string {
	var buf = bytes.NewBufferString(s)
	decoder := base64.NewDecoder(base64.StdEncoding, buf)
	var res bytes.Buffer
	res.ReadFrom(decoder)
	return res.String()
}

func AppendParameter(urlPath string, parameters map[string]string) string {
	if len(urlPath) == 0 || len(parameters) == 0 {
		return urlPath
	}
	var newParameterBuffer bytes.Buffer
	for k, v := range parameters {
		if newParameterBuffer.Len() > 0 {
			newParameterBuffer.WriteString("&")
		}
		newParameterBuffer.WriteString(k)
		newParameterBuffer.WriteString("=")
		newParameterBuffer.WriteString(url.QueryEscape(v))
	}
	prefix := "?"
	qmIndex := strings.LastIndex(urlPath, "?")
	if qmIndex > 0 {
		pcIndex := strings.LastIndex(urlPath[qmIndex:], "=")
		if pcIndex > 0 {
			prefix = "&"
		}
	}
	var newUrlPathBuffer bytes.Buffer
	newUrlPathBuffer.WriteString(urlPath)
	newUrlPathBuffer.WriteString(prefix)
	newUrlPathBuffer.WriteString(newParameterBuffer.String())
	return newUrlPathBuffer.String()
}
