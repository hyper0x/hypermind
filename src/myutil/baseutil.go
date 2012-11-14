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
)

const (
	LoginNameKey string = "login_name"
	PasswordKey string = "password"
	CnNameKey string = "cn_name"
	EmailKey string = "email"
	MobilePhoneKey string = "mobile_phone"
)

func PrintRequestInfo(prefix string, r *http.Request) {
	fmt.Println(prefix, "- form:", r.Form)
	fmt.Println(prefix, "- method:", r.Method)
	fmt.Println(prefix, "- path:", r.URL.Path)
	fmt.Println(prefix, "- scheme:", r.URL.Scheme)
}

func GenerateBasicAttrMap(r *http.Request) (attrMap map[string]string, err error) {
	attrMap = make(map[string]string)
	host, port := splitHostPort(r.Host)
	attrMap["serverAddr"] = host
	attrMap["serverPort"] = port
	loginName := GetCookie(r, LoginNameKey)
	attrMap["loginName"] = loginName
	return
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

func VerifyRegisterForm(r *http.Request) (invalidFields []string) {
	invalidFields = make([]string, 1)
	loginName := r.FormValue(LoginNameKey)
	if m, _ := regexp.MatchString("^[a-zA-Z-_\\.]{1, 10}$", loginName); !m {
		invalidFields = append(invalidFields, LoginNameKey)
	}
	password := r.FormValue(PasswordKey)
	if m, _ := regexp.MatchString("^[a-zA-Z-_\\.]{1, 20}$", password); !m {
		invalidFields = append(invalidFields, PasswordKey)
	}
	cnName := r.FormValue(CnNameKey)
	if m, _ := regexp.MatchString("^[\\x{4e00}-\\x{9fa5}]+$", cnName); !m {
		invalidFields = append(invalidFields, CnNameKey)
	}
	email := r.FormValue(EmailKey)
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, email); !m {
		invalidFields = append(invalidFields, EmailKey)
	}
	mobilePhone := r.FormValue(MobilePhoneKey)
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

