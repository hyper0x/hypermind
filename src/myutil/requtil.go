package myutil

import (
	"net/http"
	"time"
	"regexp"
	"fmt"
	"io"
	"crypto/md5"
	"os"
	"strings"
	"encoding/base64"
	"bytes"
)

var pageParameterMap map[string]string = map[string]string{
	HomePageKey: HomePage,
	AboutMePageKey: AboutMePage,
	AboutWebsitePageKey:  AboutWebsitePage,
	MeetingKanbanPageKey: MeetingKanbanPage,
}

func GeneratePagePath(reqPage string) string {
	page := "home"
	if len(reqPage) > 0 {
		page = reqPage
	}
	return "web/page/" + page + ".gtpl"
}

func PrintRequestInfo(prefix string, r *http.Request) {
	fmt.Println(prefix, "- form:", r.Form)
	fmt.Println(prefix, "- method:", r.Method)
	fmt.Println(prefix, "- path:", r.URL.Path)
	fmt.Println(prefix, "- scheme:", r.URL.Scheme)
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
		port = requestHost[splitIndex + 1:len(requestHost)]
	} else {
		host = requestHost
		config, err := ReadConfig(false)
		if err != nil {
			fmt.Println("ConfigLoadError: ", err)
			port = "80"
		} else {
			port = fmt.Sprintf("%v", config.ServerPort)
		}
	}
	return
}

func VerifyRegisterForm(r *http.Request) (fieldMap map[string]string, invalidFields []string) {
	fieldMap = make(map[string]string)
	invalidFields = make([]string, 1)
	loginName := r.FormValue(LoginNameKey)
	fieldMap[LoginNameKey] = loginName
	if m, _ := regexp.MatchString("^[a-zA-Z-_\\.]{1, 10}$", loginName); !m {
		invalidFields = append(invalidFields, LoginNameKey)
	}
	password := r.FormValue(PasswordKey)
	fieldMap[PasswordKey] = password
	if m, _ := regexp.MatchString("^[a-zA-Z-_\\.]{1, 20}$", password); !m {
		invalidFields = append(invalidFields, PasswordKey)
	}
	cnName := r.FormValue(CnNameKey)
	fieldMap[CnNameKey] = cnName
	if m, _ := regexp.MatchString("^[\\x{4e00}-\\x{9fa5}]+$", cnName); !m {
		invalidFields = append(invalidFields, CnNameKey)
	}
	email := r.FormValue(EmailKey)
	fieldMap[EmailKey] = email
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, email); !m {
		invalidFields = append(invalidFields, EmailKey)
	}
	mobilePhone := r.FormValue(MobilePhoneKey)
	fieldMap[MobilePhoneKey] = mobilePhone
	if m, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, mobilePhone); !m {
		invalidFields = append(invalidFields, MobilePhoneKey)
	}
	return
}

func DeleteTempFile(delay time.Duration, filePath string) (err error) {
	time.Sleep(delay)
	err = os.Remove(filePath)
	if err != nil {
		fmt.Printf("Occur error when delete file '%s': %s\n", filePath, err)
	} else {
		fmt.Printf("The file '%s' is deleted.\n", filePath, err)
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

