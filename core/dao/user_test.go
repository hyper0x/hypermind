package dao

import (
	"runtime/debug"
	"testing"
)

func TestUser(t *testing.T) {
	debugTag := true
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	loginName := "root0"
	password := "hypermind"
	email := "freej.cn@gmail.com"
	mobilePhone := "8618610000000"
	rights := ROOT_RIGHTS
	remark := "Testing root user"
	user0 := &User{
		LoginName:   loginName,
		Password:    password,
		Email:       email,
		MobilePhone: mobilePhone,
		Rights:      rights,
		Remark:      remark}
	if debugTag {
		t.Logf("User0: %v\n", user0)
	}
	err := AddUserToDb(user0)
	if err != nil {
		t.Errorf("Error: Add User Error: %s\n", err)
		t.FailNow()
	}
	user1, err := GetUserFromDb(loginName)
	if err != nil {
		t.Errorf("Error: Get User Error: %s\n", err)
		t.FailNow()
	}
	if debugTag {
		t.Logf("User1: %v\n", user1)
	}
	if user1.LoginName != loginName ||
		user1.Password != encryptPassword(password) ||
		user1.Email != email ||
		user1.MobilePhone != mobilePhone ||
		user1.Rights != rights ||
		user1.Remark != remark {
		t.Errorf("Error: The user should be %v but %v. (negligible password)\n", user0, user1)
		t.FailNow()
	}
	pass, err := VerifyUser(loginName, password)
	if err != nil {
		t.Errorf("Error: Verify User Error: %s\n", err)
		t.FailNow()
	}
	if !pass {
		t.Errorf("Error: The password of user (loginName=%s) should equals %s. \n", loginName, password)
		t.FailNow()
	}
	err = DeleteUserFromDb(loginName)
	if err != nil {
		t.Errorf("Error: Delete User Error: %s\n", err)
		t.FailNow()
	}
}
