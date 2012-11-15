package main

import (
	"myutil"
	"fmt"
	"net/http"
	"strings"
	"log"
	"html/template"
	"io"
	"time"
	"os"
	"bytes"
)

func welcome(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	myutil.PrintRequestInfo("welcome", r)
    for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	userInfoMap := myutil.GetStagedUserInfo(w, r)
	loginName := userInfoMap[myutil.LoginNameKey]
	attrMap := myutil.GenerateBasicAttrMap(r, (len(loginName) > 0))
	attrMap[myutil.LoginNameKey] = loginName
	t, err := template.ParseFiles("web/page/welcome.gtpl")
	if err != nil {
		fmt.Println("TemplateParseErr:", err)
	}
	err = t.Execute(w, attrMap)
    if err != nil {
        fmt.Println("PageWriteErr:", err)
    }
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	myutil.PrintRequestInfo("login", r)
	userInfoMap := myutil.GetStagedUserInfo(w, r)
	loginName := userInfoMap[myutil.LoginNameKey]
	if r.Method == "GET" {
		tokenKey := myutil.GenerateTokenKey(loginName, r)
		fmt.Println("TokenKey:", tokenKey)
		token := myutil.GenerateToken()
		fmt.Println("Token:", token)
		myutil.SetToken(tokenKey, token)
		attrMap := myutil.GenerateBasicAttrMap(r, false)
		attrMap["token"] = token
		t, err := template.ParseFiles("web/page/login.gtpl")
		if err != nil {
			fmt.Println("TemplateParseErr:", err)
		}
		err = t.Execute(w, attrMap)
        if err != nil {
            fmt.Println("PageWriteErr:", err)
        }
	} else {
		r.ParseForm()
		token := r.Form.Get("token")
		fmt.Println("Token:", token)
		validToken := false
		if token != "" {
			tokenKey := myutil.GenerateTokenKey(loginName, r)
			fmt.Println("TokenKey:", tokenKey)
			storedToken := myutil.GetToken(tokenKey)
			fmt.Println("StoredToken:", storedToken)
			if len(token) > 0 && len(storedToken)> 0 && token == storedToken {
				validToken = true
			}
		}
		loginName = template.HTMLEscapeString(r.Form.Get(myutil.LoginNameKey))
		fmt.Println("login - loginName:", loginName)
		password := template.HTMLEscapeString(r.Form.Get(myutil.PasswordKey))
		fmt.Println("login - password:", password)
		rememberMe := r.Form.Get("remember-me")
		fmt.Println("login - remember-me:", rememberMe)
		validLogin := myutil.VerifyUser(loginName, password)
		rememberMeTag := r.Form.Get("remember-me")
		if validLogin {
			if validToken {
				userInfoMap[myutil.LoginNameKey] = loginName
				onlySession := len(rememberMeTag) == 0 || rememberMeTag != "y"
				myutil.SetUserInfoToStage(userInfoMap, w, r, onlySession)
			}
		} else {
		  //
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	myutil.PrintRequestInfo("logout", r)
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	userInfoMap := myutil.GetStagedUserInfo(w, r)
	loginName := userInfoMap[myutil.LoginNameKey]
	if len(loginName) > 0 {
		myutil.RemoveUserInfoFromStage(userInfoMap, w, r)
		fmt.Printf("Logout: The user '%s' has  logout.\n", loginName)
	} else {
		fmt.Printf("Logout: The user '%s' has yet login.\n", loginName)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	myutil.PrintRequestInfo("register", r)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("web/page/register.gtpl")
		err := t.Execute(w, nil)
		if err != nil {
			fmt.Println("PageWriteErr:", err)
		}
	} else {
		invalidFields := myutil.VerifyRegisterForm(r)
		if len(invalidFields) > 0 {
			fmt.Println("There are a/some invalid field(s):", invalidFields)
			t, err := template.ParseFiles("web/page/register-after.gtpl")
			if err != nil {
				fmt.Println("TemplateParseErr:", err)
			}
			err = t.Execute(w, nil)
            if err != nil {
                fmt.Println("PageWriteErr:", err)
            }
		}
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		token := r.Form.Get("token")
		t, _ := template.ParseFiles("web/page/upload.gtpl")
		err := t.Execute(w, token)
        if err != nil {
            fmt.Println("PageWriteErr:", err)
        }
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		var buffer bytes.Buffer
		buffer.WriteString(os.TempDir())
		buffer.WriteString("/")
		buffer.WriteString(handler.Filename)
		tempFilePath := buffer.String()
		f, err := os.OpenFile(tempFilePath, os.O_WRONLY | os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		fmt.Printf("Receive a file & save to %s ...\n", tempFilePath)
		io.Copy(f, file)
		go myutil.DeleteTempFile(time.Duration(time.Minute * 5), tempFilePath)
	}
}

func main() {
	fileServer := http.FileServer(http.Dir("web"))
    http.Handle("/css/", fileServer)
    http.Handle("/js/", fileServer)
    http.Handle("/img/", fileServer)
    http.HandleFunc("/", welcome)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/upload", upload)
	config, err := myutil.ReadConfig(true)
	if err != nil {
		log.Fatal("ConfigLoadError: ", err)
	} else {
		addr := ":" + fmt.Sprintf("%v", config.ServerPort)
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Fatal("ListenAndServeError: ", err)
		}
	}
}

