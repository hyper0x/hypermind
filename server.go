package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go_lib"
	"html/template"
	"hypermind/core/base"
	"hypermind/core/request"
	"hypermind/core/rights"
	"hypermind/core/session"
	"io"
	"net/http"
	"os"
	"time"
)

var serverPort int

func init() {
	flag.IntVar(&serverPort, "port", 9091, "the server (http listen) port")
}

func requestDispatcher(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	go_lib.LogInfoln(request.GetRequestInfo(r))
	attrMap := request.GenerateBasicAttrMap(w, r)
	currentPage := r.FormValue("page")
	if len(currentPage) == 0 {
		currentPage = base.HOME_PAGE
	}
	pageRightsTag := attrMap[currentPage]
	if pageRightsTag != "true" {
		currentPage = base.HOME_PAGE
	}
	t := template.New("welcome page")
	t.Funcs(template.FuncMap{
		"equal":   request.SimpleEqual,
		"match":   request.MatchString,
		"allTrue": request.AllTrue,
	})
	t, err := t.ParseFiles(request.GeneratePagePath(currentPage),
		request.GeneratePagePath("common"),
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
	pass, err := request.VerifyAuthCode(auth_code)
	if err != nil {
		go_lib.LogErrorf("Occur error when verify auth code: %s\n", err)
		// w.WriteHeader(500)
		fmt.Fprintln(w, "Error: Something wrong when verify auth code!")
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
		fmt.Fprintln(w, "Error: Something wrong when get CV content!")
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
		token := request.GenerateToken(r, loginName)
		go_lib.LogInfof("Token: %v\n", token)
		request.SaveToken(token)
		attrMap := request.GenerateBasicAttrMap(w, r)
		attrMap[request.TOKEN_KEY] = token.Key
		hint := r.FormValue(request.HINT_KEY)
		if len(hint) > 0 {
			attrMap[request.HINT_KEY] = hint
		}
		t, err := template.ParseFiles(request.GeneratePagePath("login"), request.GeneratePagePath("common"))
		if err != nil {
			go_lib.LogErrorln("TemplateParseErr:", err)
		}
		err = t.ExecuteTemplate(w, "page", attrMap)
		if err != nil {
			go_lib.LogErrorln("PageWriteErr:", err)
		}
	} else {
		r.ParseForm()
		tokenKey := r.Form.Get(request.TOKEN_KEY)
		go_lib.LogInfoln("Token Key:", tokenKey)
		validToken := request.CheckToken(tokenKey)
		if !validToken {
			go_lib.LogWarnf("Invalid token key '%s' ! Ignore the login request.", tokenKey)
			r.Method = "GET"
			http.Redirect(w, r, r.URL.Path, http.StatusFound)
			return
		} else {
			request.RemoveToken(tokenKey)
		}
		loginName = template.HTMLEscapeString(r.Form.Get(request.LOGIN_NAME_KEY))
		go_lib.LogInfoln("login - loginName:", loginName)
		password := template.HTMLEscapeString(r.Form.Get(request.PASSWORD_KEY))
		go_lib.LogInfoln("login - password:", password)
		rememberMe := r.Form.Get("remember-me")
		go_lib.LogInfoln("login - remember-me:", rememberMe)
		validLogin, err := rights.VerifyUser(loginName, password)
		go_lib.LogInfoln("Verify user:", validLogin)
		redirectPath := "/"
		if err != nil {
			go_lib.LogErrorf("VerifyUserError (loginName=%s): %s\n", loginName, err)
			redirectPath = r.URL.Path
		} else {
			if validLogin {
				longTerm := len(rememberMe) == 0 || rememberMe != "y"
				_, err = session.NewSession(loginName, longTerm, w, r)
				if err != nil {
					go_lib.LogErrorf("SetSessionError (loginName=%s): %s\n", loginName, err)
				}
			} else {
				hint := "Wrong login name or password."
				redirectPath = request.AppendParameter(r.URL.Path, map[string]string{request.HINT_KEY: hint})
			}
		}
		go_lib.LogInfof("RPATH: %s\n", redirectPath)
		http.Redirect(w, r, redirectPath, http.StatusFound)
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
			go_lib.LogInfof("Logout: User '%s' logout. (result=%v)\n", loginName, done)
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

func getAuthCodeForAdmin(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		errorMsg := "The Web Server does not support Hijacking! "
		http.Error(w, errorMsg, http.StatusInternalServerError)
		go_lib.LogErrorf(errorMsg)
		return
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		errorMsg := "Internal error!"
		http.Error(w, errorMsg, http.StatusInternalServerError)
		go_lib.LogErrorf(errorMsg+" Hijacking Error: %s\n", err)
		return
	}
	defer conn.Close()
	r.ParseForm()
	reqType := r.FormValue("type")
	go_lib.LogInfoln(request.GetRequestInfo(r))
	attrMap := request.GenerateBasicAttrMap(w, r)
	loginName := attrMap[request.LOGIN_NAME_KEY]
	groupName := attrMap[request.GROUP_NAME_KEY]
	parameterOutline := fmt.Sprintf("[loginName=%s, groupName=%s, reqType=%s]", loginName, groupName, reqType)
	if groupName != rights.ADMIN_USER_GROUP_NAME {
		errorMsg := "Authentication failed!"
		http.Error(w, errorMsg, http.StatusForbidden)
		go_lib.LogErrorf(errorMsg+" [auth code push handler] %s \n", parameterOutline)
		return
	}
	if reqType != "lp" {
		currentAuthCode, err := request.GetCurrentAuthCode()
		if err != nil {
			go_lib.LogErrorf("GetCurrentAuthCodeError: %s\n", err)
		}
		go_lib.LogInfof("Push current auth code '%s' %s \n", currentAuthCode, parameterOutline)
		done := pushResponse(bufrw, currentAuthCode)
		if !done {
			go_lib.LogErrorf("Pushing current auth code '%s' is failing! %s \n", currentAuthCode, parameterOutline)
		}
	} else {
		nacChan := make(chan string)
		triggerFunc := func(newAuthCode string) {
			nacChan <- newAuthCode
		}
		triggerId := fmt.Sprintf("long-polling|%s|%s|%d", loginName, groupName, time.Now().UnixNano())
		request.AddNewAuthCodeTrigger(triggerId, triggerFunc)
		defer request.DelNewAuthCodeTrigger(triggerId)
		newAuthCode := <-nacChan // wait for new auth code generating
		go_lib.LogInfof("Push new auth code '%s' %s \n", newAuthCode, parameterOutline)
		done := pushResponse(bufrw, newAuthCode)
		if !done {
			go_lib.LogErrorf("Pushing new auth code '%s' is failing! %s \n", newAuthCode, parameterOutline)
		}
	}
	defer go_lib.LogInfof("The auth code push handler will be close. %s \n", parameterOutline)
}

func getUserListForAdmin(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		errorMsg := "The Web Server does not support Hijacking! "
		http.Error(w, errorMsg, http.StatusInternalServerError)
		go_lib.LogErrorf(errorMsg)
		return
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		errorMsg := "Internal error!"
		http.Error(w, errorMsg, http.StatusInternalServerError)
		go_lib.LogErrorf(errorMsg+" Hijacking Error: %s\n", err)
		return
	}
	defer conn.Close()
	r.ParseForm()
	go_lib.LogInfoln(request.GetRequestInfo(r))
	attrMap := request.GenerateBasicAttrMap(w, r)
	loginName := attrMap[request.LOGIN_NAME_KEY]
	groupName := attrMap[request.GROUP_NAME_KEY]
	parameterOutline := fmt.Sprintf("[loginName=%s, groupName=%s]", loginName, groupName)
	if groupName != rights.ADMIN_USER_GROUP_NAME {
		errorMsg := "Authentication failed!"
		http.Error(w, errorMsg, http.StatusForbidden)
		go_lib.LogErrorf(errorMsg+" [user list handler] %s \n", parameterOutline)
		return
	}
	var respBuffer bytes.Buffer
	users, err := rights.FindUser("*")
	if err != nil {
		go_lib.LogErrorf("FindUserError: %s\n", err)
	} else {
		b, err := json.Marshal(users)
		if err != nil {
			go_lib.LogErrorf("JsonMarshalError (source=%v): %s\n", users, err)
		} else {
			respBuffer.WriteString(string(b))
		}
	}
	resp := respBuffer.String()
	done := pushResponse(bufrw, resp)
	if !done {
		go_lib.LogErrorf("Pushing user list '%s' is failing! %s \n", resp, parameterOutline)
	}
}

func pushResponse(bufrw *bufio.ReadWriter, authCode string) bool {
	_, err := bufrw.Write([]byte(authCode))
	if err == nil {
		err = bufrw.Flush()
	}
	if err != nil {
		go_lib.LogErrorf("PushAuthCodeError: %s\n", err)
		return false
	}
	return true
}

func main() {
	flag.Parse()
	fileServer := http.FileServer(http.Dir("web"))
	http.Handle("/css/", fileServer)
	http.Handle("/js/", fileServer)
	http.Handle("/img/", fileServer)
	http.HandleFunc("/", requestDispatcher)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/get-cv", getCv)
	http.HandleFunc("/auth_code", getAuthCodeForAdmin)
	http.HandleFunc("/user_list", getUserListForAdmin)
	go_lib.LogInfof("Starting hypermind http server (port=%d)...\n", serverPort)
	err := http.ListenAndServe(":"+fmt.Sprintf("%d", serverPort), nil)
	if err != nil {
		go_lib.LogFatalln("ListenAndServeError: ", err)
	} else {
		go_lib.LogInfoln("Hypermind http server is started.")
	}
}
