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
	attrMap, err := myutil.GenerateBasicAttrMap(r)
	if err != nil {
		fmt.Println("BasicAttrMapGenErr:", err)
	}
	loginName := attrMap["loginName"]
	if len(loginName) > 0 {
		attrMap["welcomePrefix"] = loginName + ", "
	} else {
		attrMap["welcomePrefix"] = ""
	}
	t, err := template.ParseFiles(myutil.GetAbsolutePathOfPage("welcome.gtpl"))
	if err != nil {
		fmt.Println("TemplateParseErr:", err)
	}
	t.Execute(w, attrMap)
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	myutil.PrintRequestInfo("login", r)
	session := myutil.GetSession(w, r)
	var loginName string
	if v := session.Get(myutil.LoginNameKey);v != nil {
		loginName = v.(string)
	}
	if r.Method == "GET" {
		tokenKey := myutil.GenerateTokenKey(loginName, r)
		fmt.Println("TokenKey:", tokenKey)
		token := myutil.GenerateToken()
		fmt.Println("Token:", token)
		myutil.SetToken(tokenKey, token)
		attrMap, err := myutil.GenerateBasicAttrMap(r)
		if err != nil {
			fmt.Println("BasicAttrMapGenErr:", err)
		}
		attrMap["token"] = token
		t, err := template.ParseFiles(myutil.GetAbsolutePathOfPage("login.gtpl"))
		if err != nil {
			fmt.Println("TemplateParseErr:", err)
		}
		t.Execute(w, attrMap)
	} else {
		r.ParseForm()
		loginName = template.HTMLEscapeString(r.Form.Get(myutil.LoginNameKey))
		fmt.Println("login - login_name:", loginName)
		password := template.HTMLEscapeString(r.Form.Get(myutil.PasswordKey))
		fmt.Println("login - password:", password)
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
		if validToken {
			myutil.SetCookie(w, myutil.LoginNameKey, loginName, 5)
		}
		attrMap, err := myutil.GenerateBasicAttrMap(r)
		if err != nil {
			fmt.Println("BasicAttrMapGenErr:", err)
		}
		attrMap["loginName"] = loginName
		attrMap["validToken"] = fmt.Sprintf("%v", validToken)
		t, err := template.ParseFiles(myutil.GetAbsolutePathOfPage("login-after.gtpl"))
		if err != nil {
			fmt.Println("TemplateParseErr:", err)
		}
		t.Execute(w, attrMap)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	myutil.PrintRequestInfo("register", r)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("register.gtpl")
		t.Execute(w, nil)
	} else {
		invalidFields := myutil.VerifyRegisterForm(r)
		if len(invalidFields) > 0 {
			fmt.Println("There are a/some invalid field(s):", invalidFields)
			t, err := template.ParseFiles("register-after.gtpl")
			if err != nil {
				fmt.Println("TemplateParseErr:", err)
			}
			t.Execute(w, nil)
		}
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		token := r.Form.Get("token")
		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
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
	http.HandleFunc("/", welcome)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
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

