package rights

import (
	"bytes"
	"fmt"
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
	loginName := "testing_user"
	password := "hypermind"
	email := "freej.cn@gmail.com"
	mobilePhone := "8618610000000"
	group := NORMAL_USER_GROUP_NAME
	remark := "Testing user"
	user0 := &User{
		LoginName:   loginName,
		Password:    password,
		Email:       email,
		MobilePhone: mobilePhone,
		Group:       group,
		Remark:      remark}
	if debugTag {
		t.Logf("User0: %v\n", user0)
	}
	err := AddUser(user0)
	if err != nil {
		t.Errorf("Error: Add User Error: %s\n", err)
		t.FailNow()
	}
	user1, err := GetUser(loginName)
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
		user1.Group != group ||
		user1.Remark != remark {
		t.Errorf("Fail: The user should be %v but %v. (negligible password)\n", user0, user1)
		t.FailNow()
	}
	pass, err := VerifyUser(loginName, password)
	if err != nil {
		t.Errorf("Error: Verify User Error: %s\n", err)
		t.FailNow()
	}
	if !pass {
		t.Errorf("Fail: The password of user (loginName=%s) should equals %s. \n", loginName, password)
		t.FailNow()
	}
	loginNamePattern := loginName + "*"
	users, err := FindUser(loginNamePattern)
	if err != nil {
		t.Errorf("Error: Find User Error: %s\n", err)
		t.FailNow()
	}
	expectUsersLen := 1
	usersLen := len(users)
	if usersLen != expectUsersLen {
		t.Errorf("Fail: The length of user list (pattern=%s) should be %d but %v. (negligible password)\n", loginNamePattern, expectUsersLen, usersLen)
		t.FailNow()
	}
	user2 := users[0]
	if user2.LoginName != loginName ||
		user2.Password != encryptPassword(password) ||
		user2.Email != email ||
		user2.MobilePhone != mobilePhone ||
		user2.Group != group ||
		user2.Remark != remark {
		t.Errorf("Fail: The user should be %v but %v. (negligible password)\n", user0, user2)
		t.FailNow()
	}
	if debugTag {
		var buffer bytes.Buffer
		buffer.WriteString("[")
		for i, u := range users {
			buffer.WriteString(fmt.Sprintf("%s", u))
			if i+1 < usersLen {
				buffer.WriteString(", ")
			}
		}
		buffer.WriteString("]")
		t.Logf("Find user with pattern '%s': %s\n", loginNamePattern, buffer.String())
	}
	err = DeleteUser(loginName)
	if err != nil {
		t.Errorf("Error: Delete User Error: %s\n", err)
		t.FailNow()
	}
	user3, err := GetUser(loginName)
	if err != nil {
		t.Errorf("Error: Get User Error: %s\n", err)
		t.FailNow()
	}
	if user3 != nil {
		t.Errorf("Fail: The user '%v' should be deleted. %s\n", user3)
		t.FailNow()
	}
}
