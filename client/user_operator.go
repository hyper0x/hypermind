package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"hypermind/core/rights"
	"regexp"
)

var op string
var loginName string
var password string
var email string
var mobilePhone string
var groupName string
var remark string
var pattern string

func init() {
	flag.StringVar(&op, "op", "", "operation type")
	flag.StringVar(&loginName, "name", "", "user login name")
	flag.StringVar(&password, "password", "", "user password")
	flag.StringVar(&mobilePhone, "mphone", "", "user mobile phone")
	flag.StringVar(&groupName, "group", "", "user group name")
	flag.StringVar(&remark, "remark", "", "remark for user")
	flag.StringVar(&pattern, "pattern", "", "user login name pattern")
}

// operation type
const (
	OP_ADD    = "add"
	OP_DELETE = "del"
	OP_GET    = "get"
	OP_FIND   = "find"
)

// group name
const (
	GROUP_ADMIN  = "admin"
	GROUP_NORMAL = "normal"
)

// parameter rule
const (
	NAME_PATTERN     = "^[a-zA-z0-9_]{6}[a-zA-z0-9_]*$"
	PASSWORD_PATTERN = "^[a-zA-z0-9_\\.]{8}[a-zA-z0-9_\\.]*$"
	PATTERN_PATTERN  = "[a-zA-z0-9_\\[\\]\\*\\?]+$"
)

func printUsinghelp() {
	separator := "#####################"
	var buffer bytes.Buffer
	buffer.WriteString(separator)
	buffer.WriteString("\n\n")
	buffer.WriteString("User Operator Usage:\n")
	buffer.WriteString("\n")
	buffer.WriteString("  Examples:\n")
	buffer.WriteString("    - Add User:    go run <file name> -op=add -name=user1 -password=abcefghi -group=normal \n")
	buffer.WriteString("    - Delete User: go run <file name> -op=del -name=user1 \n")
	buffer.WriteString("    - Get User:    go run <file name> -op=get -name=user1 \n")
	buffer.WriteString("    - Find User:   go run <file name> -op=find -pattern=user* \n")
	buffer.WriteString("\n")
	buffer.WriteString("  Required parameters:\n")
	buffer.WriteString("    * op:\n")
	buffer.WriteString("        operation type. Its value should be 'add' or 'del' or 'get' or 'find'. \n")
	buffer.WriteString("    * name:\n")
	buffer.WriteString("        user login name. It's required when op not equals 'find'. \n")
	buffer.WriteString("        Its value should be at least six english alphabet or digital or '_'. \n")
	buffer.WriteString("    * passward:\n")
	buffer.WriteString("        user password. It's required when op equals 'add'. \n")
	buffer.WriteString("        Its value should be at least eight english alphabet or digital or '_' or '.'. \n")
	buffer.WriteString("    * group:\n")
	buffer.WriteString("        user group. It's required when op equals 'add'. Its value should be 'admin' or 'normal'. \n")
	buffer.WriteString("    * pattern:\n")
	buffer.WriteString("        user login name pattern. It's required when op equals 'find'.\n")
	buffer.WriteString("        See also the description of command 'KEYS' in redis document. \n")
	buffer.WriteString("\n")
	buffer.WriteString("  Optional parameters:\n")
	buffer.WriteString("    + email:\n")
	buffer.WriteString("        user email. The default value is \"\". It's useless when op not equals 'add'. \n")
	buffer.WriteString("    + mphone:\n")
	buffer.WriteString("        user mobile phone. The default value is \"\". It's useless when op not equals 'add'. \n")
	buffer.WriteString("    + remark:\n")
	buffer.WriteString("        The remark for user. The default value is \"\". It's useless when op not equals 'add'. \n")
	buffer.WriteString("\n")
	buffer.WriteString(separator)
	buffer.WriteString("\n")
	fmt.Println(buffer.String())
}

func matchString(content string, pattern string) bool {
	pass, err := regexp.MatchString(pattern, content)
	if err != nil {
		fmt.Printf("RegexpMatchError (content=%s, pattern=%s): %s\n", content, pattern, err)
		return false
	}
	return pass
}

func addUser() {
	user := &rights.User{
		LoginName:   loginName,
		Password:    password,
		Email:       email,
		MobilePhone: mobilePhone,
		Group:       groupName,
		Remark:      remark}
	err := rights.AddUser(user)
	if err != nil {
		fmt.Printf("AddUserError: %s\n", err)
	}
	fmt.Printf("The user (loginName=%s) was added. \n", loginName)
}

func deleteUser() {
	err := rights.DeleteUser(loginName)
	if err != nil {
		fmt.Printf("DeleteUserError: %s\n", err)
	}
	fmt.Printf("The user (loginName=%s) was deleted. \n", loginName)
}

func getUser() {
	user, err := rights.GetUser(loginName)
	if err != nil {
		fmt.Printf("GetUserError: %s\n", err)
	}
	var buffer bytes.Buffer
	buffer.WriteString("Get User:\n")
	buffer.WriteString(getUserLiterals(user))
	buffer.WriteString("\n")
	fmt.Println(buffer.String())
}

func findUser() {
	users, err := rights.FindUser(pattern)
	if err != nil {
		fmt.Printf("FindUserError: %s\n", err)
	}
	var buffer bytes.Buffer
	buffer.WriteString("Find User:\n")
	buffer.WriteString(getUserLiterals(users...))
	buffer.WriteString("\n")
	fmt.Println(buffer.String())
}

func getUserLiterals(users ...*rights.User) string {
	length := len(users)
	if length == 0 {
		return "<Nil>"
	}
	separator := "---------"
	var buffer bytes.Buffer
	if length > 1 {
		buffer.WriteString(separator)
		buffer.WriteString("\n")
	}
	for i, v := range users {
		if length > 1 {
			buffer.WriteString("User [" + fmt.Sprintf("%d", i) + "]:\n")
		}
		b, err := json.Marshal(v)
		if err != nil {
			buffer.WriteString("<JsonMarshalError: " + fmt.Sprintf("%s", err) + ">")
		} else {
			buffer.WriteString(string(b))
		}
		buffer.WriteString("\n")
		if length > 1 {
			buffer.WriteString(separator)
			buffer.WriteString("\n")
		}
	}
	return buffer.String()
}

func checkParameter() (bool, string) {
	var msg string
	switch op {
	case OP_ADD:
		validLoginName := matchString(loginName, NAME_PATTERN)
		if !validLoginName {
			msg = fmt.Sprintf("Invalid name '%s'!\n", loginName)
			break
		}
		validPassword := matchString(password, PASSWORD_PATTERN)
		if !validPassword {
			msg = fmt.Sprintf("Invalid password '%s'!\n", password)
			break
		}
		validGroupName := groupName == GROUP_ADMIN || groupName == GROUP_NORMAL
		if !validGroupName {
			msg = fmt.Sprintf("Invalid group '%s'!\n", groupName)
			break
		}
	case OP_DELETE:
		validLoginName := len(loginName) > 0
		if !validLoginName {
			msg = fmt.Sprintf("Invalid name '%s'!\n", loginName)
			break
		}
	case OP_GET:
		validLoginName := len(loginName) > 0
		if !validLoginName {
			msg = fmt.Sprintf("Invalid name '%s'!\n", loginName)
			break
		}
	case OP_FIND:
		validPattern := matchString(pattern, PATTERN_PATTERN)
		if !validPattern {
			msg = fmt.Sprintf("Invalid pattern '%s'!\n", pattern)
			break
		}
	default:
		msg = fmt.Sprintf("Invalid op '%s'!\n", op)
	}
	pass := false
	if len(msg) == 0 {
		pass = true
	}
	return pass, msg
}

func main() {
	fmt.Printf("Start User Operater ...\n\n")
	flag.Parse()
	pass, msg := checkParameter()
	if !pass {
		fmt.Printf("ParameterError: %s\n", msg)
		printUsinghelp()
		return
	}
	switch op {
	case OP_ADD:
		addUser()
	case OP_DELETE:
		deleteUser()
	case OP_GET:
		getUser()
	case OP_FIND:
		findUser()
	default:
		fmt.Printf("OperationError: Invalid operation '%s'!\n", op)
	}
}
