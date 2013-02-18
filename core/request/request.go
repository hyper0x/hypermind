package request

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go_lib"
	"hypermind/core/base"
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
	HOME_PAGE_KEY:           HOME_PAGE,
	ABOUT_ME_PAGE_KEY:       ABOUT_ME_PAGE,
	ABOUT_WEBSITE_PAGE_KEY:  ABOUT_WEBSITE_PAGE,
	MEETING_KANBAN_PAGE_KEY: MEETING_KANBAN_PAGE,
	PROJECT_HASH_RING_KEY:   PROJECT_HASH_RING_PAGE,
	ADMIN_CV_KEY:            ADMIN_CV_PAGE,
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
		go_lib.LogErrorln("JsonMarshalError:", err)
	}
	return string(b)
}

func GenerateBasicAttrMap(r *http.Request, validLogin bool) map[string]string {
	attrMap := make(map[string]string)
	host, port := splitHostPort(r.Host)
	attrMap["serverAddr"] = host
	attrMap["serverPort"] = port
	if validLogin {
		attrMap["validLogin"] = "true"
	}
	for pageKey, page := range pageParameterMap {
		attrMap[pageKey] = page
	}
	return attrMap
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
			go_lib.LogErrorln("ConfigLoadError: ", err)
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
		go_lib.LogErrorf("Occur error when delete file '%s': %s\n", filePath, err)
	} else {
		go_lib.LogInfof("The file '%s' is deleted.\n", filePath, err)
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
