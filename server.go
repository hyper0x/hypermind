package main

import (
	"bytes"
	"flag"
	"fmt"
	"go_lib"
	"html/template"
	"hypermind/core/base"
	"hypermind/core/dao"
	"hypermind/core/request"
	"hypermind/core/rights"
	"hypermind/core/session"
	"io"
	"net/http"
	"os"
	"time"
)

var serverPort int = *flag.Int("port", 9091, "the server (http listen) port")

func welcome(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	go_lib.LogInfoln(request.GetRequestInfo(r))
	attrMap := request.GenerateBasicAttrMap(w, r)
	currentPage := r.FormValue("page")
	if len(currentPage) == 0 {
		currentPage = request.HOME_PAGE
	}
	pageRightsTag := attrMap[currentPage]
	if pageRightsTag != "true" {
		currentPage = request.HOME_PAGE
	}
	t := template.New("welcome page")
	t.Funcs(template.FuncMap{
		"equal":   request.SimpleEqual,
		"match":   request.MatchString,
		"allTrue": request.AllTrue,
	})
	t, err := t.ParseFiles(request.GeneratePagePath(currentPage),
		request.GeneratePagePath("header"),
		request.GeneratePagePath("footer"),
		request.GeneratePagePath("navbar"))
	if err != nil {
		go_lib.LogErrorln("ParseFilesErr:", err)
	}
	attrMap["currentPage"] = currentPage
	err = t.ExecuteTemplate(w, "page", attrMap)
	if err != nil {
		go_lib.LogErrorln("ExecuteTemplateErr:", err)
	}
}

func getCv(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	attrMap := request.GenerateBasicAttrMap(w, r)
	loginName := attrMap[request.LOGIN_NAME_KEY]
	go_lib.LogInfoln(request.GetRequestInfo(r))
	auth_code := r.FormValue(request.AUTH_CODE)
	go_lib.LogInfof("Getting CV by user '%s' with input '%s'...\n", loginName, auth_code)
	pass, err := dao.VerifyAuthCode(auth_code)
	if err != nil {
		go_lib.LogErrorf("Occur error when verify auth code: %s\n", err)
		// w.WriteHeader(500)
		fmt.Fprintln(w, "Error: Somethin wrong when verify auth code!")
		return
	}
	if !pass {
		go_lib.LogWarnf("Unauthorized CV getting by user '%s' with input '%s'.\n", loginName, auth_code)
		// w.WriteHeader(401)
		fmt.Fprintln(w, "FAIL: Wrong authorization code.")
		return
	}
	cvContent, err := base.GetCvContent()
	if err != nil {
		go_lib.LogErrorf("Occur error when get cv content: %s.\n", err)
		// w.WriteHeader(500)
		fmt.Fprintln(w, "Error: Somethin wrong when get CV content!")
		return
	}
	fmt.Fprintln(w, cvContent)
	go_lib.LogInfof("The CV had taken by user '%s' with input '%s'.\n", loginName, auth_code)
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	go_lib.LogInfoln(request.GetRequestInfo(r))
	attrMap := request.GenerateBasicAttrMap(w, r)
	loginName := attrMap[request.LOGIN_NAME_KEY]
	if r.Method == "GET" {
		tokenKey := request.GenerateTokenKey(loginName, r)
		go_lib.LogInfoln("TokenKey:", tokenKey)
		token := request.GenerateToken()
		go_lib.LogInfo("Token:", token)
		request.SetToken(tokenKey, token)
		attrMap := request.GenerateBasicAttrMap(w, r)
		attrMap["token"] = token
		t, err := template.ParseFiles(request.GeneratePagePath("login"))
		if err != nil {
			go_lib.LogErrorln("TemplateParseErr:", err)
		}
		err = t.Execute(w, attrMap)
		if err != nil {
			go_lib.LogErrorln("PageWriteErr:", err)
		}
	} else {
		r.ParseForm()
		token := r.Form.Get("token")
		go_lib.LogInfoln("Token:", token)
		validToken := false
		if token != "" {
			tokenKey := request.GenerateTokenKey(loginName, r)
			go_lib.LogInfoln("TokenKey:", tokenKey)
			storedToken := request.GetToken(tokenKey)
			go_lib.LogInfoln("StoredToken:", storedToken)
			if len(token) > 0 && len(storedToken) > 0 && token == storedToken {
				validToken = true
			}
		}
		loginName = template.HTMLEscapeString(r.Form.Get(request.LOGIN_NAME_KEY))
		go_lib.LogInfoln("login - loginName:", loginName)
		password := template.HTMLEscapeString(r.Form.Get(request.PASSWORD_KEY))
		go_lib.LogInfoln("login - password:", password)
		rememberMe := r.Form.Get("remember-me")
		go_lib.LogInfoln("login - remember-me:", rememberMe)
		validLogin, err := rights.VerifyUser(loginName, password)
		go_lib.LogInfoln("Verify user:", validLogin)
		if err != nil {
			go_lib.LogErrorf("VerifyUserError (loginName=%s): %s\n", loginName, err)
		} else {
			if validLogin {
				if validToken {
					longTerm := len(rememberMe) == 0 || rememberMe != "y"
					_, err = session.NewSession(loginName, longTerm, w, r)
					if err != nil {
						go_lib.LogErrorf("SetSessionError (loginName=%s): %s\n", loginName, err)
					}
				}
			}
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	go_lib.LogInfoln(request.GetRequestInfo(r))
	hmSession, err := session.GetMatchedSession(w, r)
	if err != nil {
		go_lib.LogErrorln("GetSessionErr:", err)
	}
	if hmSession != nil {
		loginName, err := hmSession.Get(session.SESSION_GRANTORS_KEY)
		if err != nil {
			go_lib.LogErrorln("GetLoginNameErr:", err)
		}
		done, err := hmSession.Destroy()
		if err != nil {
			go_lib.LogErrorln("DestroySessionErr:", err)
		} else {
			go_lib.LogInfoln("Logout: User '%s' logout. (result=%v)\n", loginName, done)
		}
	} else {
		go_lib.LogInfoln("Logout: Current visitor has yet login.\n")
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	go_lib.LogInfoln(request.GetRequestInfo(r))
	if r.Method == "GET" {
		attrMap := request.GenerateBasicAttrMap(w, r)
		encodedHint := r.FormValue("hint")
		if len(encodedHint) > 0 {
			hint := request.UrlDecoding(encodedHint)
			attrMap["hint"] = hint
		}
		t, _ := template.ParseFiles(request.GeneratePagePath("register"))
		err := t.Execute(w, attrMap)
		if err != nil {
			go_lib.LogErrorln("PageWriteErr:", err)
		}
	} else {
		fieldMap, invalidFields := request.VerifyRegisterForm(r)
		go_lib.LogInfoln("The field map:", fieldMap)
		if len(invalidFields) > 0 {
			hint := fmt.Sprintln("There are some invalid fields of '':", invalidFields, ".")
			go_lib.LogInfoln(hint)
			encodedHint := request.UrlEncoding(hint)
			redirectUrl := "/register?hint=" + encodedHint
			http.Redirect(w, r, redirectUrl, http.StatusFound)
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	go_lib.LogInfoln(request.GetRequestInfo(r))
	if r.Method == "GET" {
		token := r.Form.Get("token")
		t, _ := template.ParseFiles(request.GeneratePagePath("upload"))
		err := t.Execute(w, token)
		if err != nil {
			go_lib.LogErrorln("PageWriteErr:", err)
		}
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			go_lib.LogErrorln("UploadFileParsError:", err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		var buffer bytes.Buffer
		buffer.WriteString(os.TempDir())
		buffer.WriteString("/")
		buffer.WriteString(handler.Filename)
		tempFilePath := buffer.String()
		f, err := os.OpenFile(tempFilePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			go_lib.LogErrorln(err)
			return
		}
		defer f.Close()
		go_lib.LogInfoln("Receive a file & save to %s ...\n", tempFilePath)
		io.Copy(f, file)
		go request.DeleteTempFile(time.Duration(time.Minute*5), tempFilePath)
	}
}

func main() {
	flag.Parse()
	fileServer := http.FileServer(http.Dir("web"))
	http.Handle("/css/", fileServer)
	http.Handle("/js/", fileServer)
	http.Handle("/img/", fileServer)
	http.HandleFunc("/", welcome)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/get-cv", getCv)
	go_lib.LogInfof("Starting hypermind http server (port=%d)...\n", serverPort)
	err := http.ListenAndServe(":"+fmt.Sprintf("%d", serverPort), nil)
	if err != nil {
		go_lib.LogFatalln("ListenAndServeError: ", err)
	} else {
		go_lib.LogInfoln("Hypermind http server is started.")
	}
}
