package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hypermind/core/base"
	"hypermind/core/request"
	"hypermind/core/rights"
	"net/http"
	"time"
)

func GetAuthCodeForAdmin(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		errorMsg := "The Web Server does not support Hijacking! "
		http.Error(w, errorMsg, http.StatusInternalServerError)
		base.Logger().Errorf(errorMsg)
		return
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		errorMsg := "Internal error!"
		http.Error(w, errorMsg, http.StatusInternalServerError)
		base.Logger().Errorf(errorMsg+" Hijacking Error: %s\n", err)
		return
	}
	defer conn.Close()
	r.ParseForm()
	reqType := r.FormValue("type")
	base.Logger().Infoln(request.GetRequestInfo(r))
	attrMap := request.GenerateBasicAttrMap(w, r)
	loginName := attrMap[request.LOGIN_NAME_KEY]
	groupName := attrMap[request.GROUP_NAME_KEY]
	parameterOutline := fmt.Sprintf("[loginName=%s, groupName=%s, reqType=%s]", loginName, groupName, reqType)
	if groupName != rights.ADMIN_USER_GROUP_NAME {
		errorMsg := "Authentication failed!"
		http.Error(w, errorMsg, http.StatusForbidden)
		base.Logger().Errorf(errorMsg+" [auth code push handler] %s \n", parameterOutline)
		return
	}
	if reqType != "lp" {
		currentAuthCode, err := request.GetCurrentAuthCode()
		if err != nil {
			base.Logger().Errorf("GetCurrentAuthCodeError: %s\n", err)
		}
		base.Logger().Infof("Push current auth code '%s' %s \n", currentAuthCode, parameterOutline)
		done := pushResponse(bufrw, currentAuthCode)
		if !done {
			base.Logger().Errorf("Pushing current auth code '%s' is failing! %s \n", currentAuthCode, parameterOutline)
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
		base.Logger().Infof("Push new auth code '%s' %s \n", newAuthCode, parameterOutline)
		done := pushResponse(bufrw, newAuthCode)
		if !done {
			base.Logger().Errorf("Pushing new auth code '%s' is failing! %s \n", newAuthCode, parameterOutline)
		}
	}
	defer base.Logger().Infof("The auth code push handler will be close. %s \n", parameterOutline)
}

func GetUserListForAdmin(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		errorMsg := "The Web Server does not support Hijacking! "
		http.Error(w, errorMsg, http.StatusInternalServerError)
		base.Logger().Errorf(errorMsg)
		return
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		errorMsg := "Internal error!"
		http.Error(w, errorMsg, http.StatusInternalServerError)
		base.Logger().Errorf(errorMsg+" Hijacking Error: %s\n", err)
		return
	}
	defer conn.Close()
	r.ParseForm()
	base.Logger().Infoln(request.GetRequestInfo(r))
	attrMap := request.GenerateBasicAttrMap(w, r)
	loginName := attrMap[request.LOGIN_NAME_KEY]
	groupName := attrMap[request.GROUP_NAME_KEY]
	parameterOutline := fmt.Sprintf("[loginName=%s, groupName=%s]", loginName, groupName)
	if groupName != rights.ADMIN_USER_GROUP_NAME {
		errorMsg := "Authentication failed!"
		http.Error(w, errorMsg, http.StatusForbidden)
		base.Logger().Errorf(errorMsg+" [user list handler] %s \n", parameterOutline)
		return
	}
	var respBuffer bytes.Buffer
	users, err := rights.FindUser("*")
	if err != nil {
		base.Logger().Errorf("FindUserError: %s\n", err)
	} else {
		b, err := json.Marshal(users)
		if err != nil {
			base.Logger().Errorf("JsonMarshalError (source=%v): %s\n", users, err)
		} else {
			respBuffer.WriteString(string(b))
		}
	}
	resp := respBuffer.String()
	done := pushResponse(bufrw, resp)
	if !done {
		base.Logger().Errorf("Pushing user list '%s' is failing! %s \n", resp, parameterOutline)
	}
}
