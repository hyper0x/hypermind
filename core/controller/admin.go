package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go_lib"
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

func GetUserListForAdmin(w http.ResponseWriter, r *http.Request) {
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
