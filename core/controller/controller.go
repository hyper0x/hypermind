package controller

import (
	"bytes"
	"fmt"
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

func RequestDispatcher(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	base.Logger().Infoln(request.GetRequestInfo(r))
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
		base.Logger().Errorln("ParseFilesErr:", err)
	}
	attrMap["currentPage"] = currentPage
	err = t.ExecuteTemplate(w, "page", attrMap)
	if err != nil {
		base.Logger().Errorln("ExecuteTemplateErr:", err)
	}
	recordPageAccessInfo(currentPage, attrMap[request.LOGIN_NAME_KEY], uint64(1))
}

func GetCv(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	attrMap := request.GenerateBasicAttrMap(w, r)
	loginName := attrMap[request.LOGIN_NAME_KEY]
	base.Logger().Infoln(request.GetRequestInfo(r))
	auth_code := r.FormValue(request.AUTH_CODE)
	base.Logger().Infof("Getting CV by user '%s' with input '%s'...\n", loginName, auth_code)
	pass, err := request.VerifyAuthCode(auth_code)
	if err != nil {
		base.Logger().Errorf("Occur error when verify auth code: %s\n", err)
		// w.WriteHeader(500)
		fmt.Fprintln(w, "Error: Something wrong when verify auth code!")
		return
	}
	if !pass {
		base.Logger().Warnf("Unauthorized CV getting by user '%s' with input '%s'.\n", loginName, auth_code)
		// w.WriteHeader(401)
		fmt.Fprintln(w, "FAIL: Wrong authorization code.")
		return
	}
	cvContent, err := base.GetCvContent()
	if err != nil {
		base.Logger().Errorf("Occur error when get cv content: %s.\n", err)
		// w.WriteHeader(500)
		fmt.Fprintln(w, "Error: Something wrong when get CV content!")
		return
	}
	fmt.Fprintln(w, cvContent)
	base.Logger().Infof("The CV had taken by user '%s' with input '%s'.\n", loginName, auth_code)
}

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	base.Logger().Infoln(request.GetRequestInfo(r))
	attrMap := request.GenerateBasicAttrMap(w, r)
	loginName := attrMap[request.LOGIN_NAME_KEY]
	if r.Method == "GET" {
		token := request.GenerateToken(r, loginName)
		base.Logger().Infof("Token: %v\n", token)
		request.SaveToken(token)
		attrMap := request.GenerateBasicAttrMap(w, r)
		attrMap[request.TOKEN_KEY] = token.Key
		hint := r.FormValue(request.HINT_KEY)
		if len(hint) > 0 {
			attrMap[request.HINT_KEY] = hint
		}
		t, err := template.ParseFiles(request.GeneratePagePath("login"), request.GeneratePagePath("common"))
		if err != nil {
			base.Logger().Errorln("TemplateParseErr:", err)
		}
		err = t.ExecuteTemplate(w, "page", attrMap)
		if err != nil {
			base.Logger().Errorln("PageWriteErr:", err)
		}
	} else {
		r.ParseForm()
		tokenKey := r.Form.Get(request.TOKEN_KEY)
		base.Logger().Infoln("Token Key:", tokenKey)
		validToken := request.CheckToken(tokenKey)
		if !validToken {
			base.Logger().Warnf("Invalid token key '%s' ! Ignore the login request.", tokenKey)
			r.Method = "GET"
			http.Redirect(w, r, r.URL.Path, http.StatusFound)
			return
		} else {
			request.RemoveToken(tokenKey)
		}
		loginName = template.HTMLEscapeString(r.Form.Get(request.LOGIN_NAME_KEY))
		base.Logger().Infoln("login - loginName:", loginName)
		password := template.HTMLEscapeString(r.Form.Get(request.PASSWORD_KEY))
		base.Logger().Infoln("login - password:", password)
		rememberMe := r.Form.Get("remember-me")
		base.Logger().Infoln("login - remember-me:", rememberMe)
		validLogin, err := rights.VerifyUser(loginName, password)
		base.Logger().Infoln("Verify user:", validLogin)
		redirectPath := "/"
		if err != nil {
			base.Logger().Errorf("VerifyUserError (loginName=%s): %s\n", loginName, err)
			redirectPath = r.URL.Path
		} else {
			if validLogin {
				longTerm := len(rememberMe) == 0 || rememberMe != "y"
				_, err = session.NewSession(loginName, longTerm, w, r)
				if err != nil {
					base.Logger().Errorf("SetSessionError (loginName=%s): %s\n", loginName, err)
				}
			} else {
				hint := "Wrong login name or password."
				redirectPath = request.AppendParameter(r.URL.Path, map[string]string{request.HINT_KEY: hint})
			}
		}
		base.Logger().Infof("RPATH: %s\n", redirectPath)
		http.Redirect(w, r, redirectPath, http.StatusFound)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	base.Logger().Infoln(request.GetRequestInfo(r))
	hmSession, err := session.GetMatchedSession(w, r)
	if err != nil {
		base.Logger().Errorln("GetSessionErr:", err)
	}
	if hmSession != nil {
		loginName, err := hmSession.Get(session.SESSION_GRANTORS_KEY)
		if err != nil {
			base.Logger().Errorln("GetLoginNameErr:", err)
		}
		done, err := hmSession.Destroy()
		if err != nil {
			base.Logger().Errorln("DestroySessionErr:", err)
		} else {
			base.Logger().Infof("Logout: User '%s' logout. (result=%v)\n", loginName, done)
		}
	} else {
		base.Logger().Infoln("Logout: Current visitor has yet login.\n")
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	base.Logger().Infoln(request.GetRequestInfo(r))
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
			base.Logger().Errorln("PageWriteErr:", err)
		}
	} else {
		fieldMap, invalidFields := request.VerifyRegisterForm(r)
		base.Logger().Infoln("The field map:", fieldMap)
		if len(invalidFields) > 0 {
			hint := fmt.Sprintln("There are some invalid fields of '':", invalidFields, ".")
			base.Logger().Infoln(hint)
			encodedHint := request.UrlEncoding(hint)
			redirectUrl := "/register?hint=" + encodedHint
			http.Redirect(w, r, redirectUrl, http.StatusFound)
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func Upload(w http.ResponseWriter, r *http.Request) {
	base.Logger().Infoln(request.GetRequestInfo(r))
	if r.Method == "GET" {
		token := r.Form.Get("token")
		t, _ := template.ParseFiles(request.GeneratePagePath("upload"))
		err := t.Execute(w, token)
		if err != nil {
			base.Logger().Errorln("PageWriteErr:", err)
		}
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			base.Logger().Errorln("UploadFileParsError:", err)
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
			base.Logger().Errorln(err)
			return
		}
		defer f.Close()
		base.Logger().Infoln("Receive a file & save to %s ...\n", tempFilePath)
		io.Copy(f, file)
		go request.DeleteTempFile(time.Duration(time.Minute*5), tempFilePath)
	}
}
