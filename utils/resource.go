package utils

import (
	"bufio"
	"bytes"
	"go_lib"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const (
	AUTH_CODE_FILE_NAME = "auth_code.txt"
	CV_FILE_NAME        = "resume.html"
)

func VerifyAuthCode(authCode string) (bool, error) {
	if len(authCode) == 0 {
		return false, nil
	}
	filePath := GenerateResourceFilePath(AUTH_CODE_FILE_NAME)
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	defer file.Close()
	if err != nil {
		return false, err
	}
	currentAuthCode, err := GetCurrentAuthCode()
	if err != nil {
		return false, err
	}
	if len(currentAuthCode) == 0 {
		warningMsg := "Warning: No any auth code! Please initialize a auth code."
		go_lib.LogWarnln(warningMsg)
		return false, nil
	}
	if currentAuthCode != strings.TrimSpace(authCode) {
		return false, nil
	}
	return true, nil
}

func GetCvContent() (string, error) {
	filePath := GenerateResourceFilePath(CV_FILE_NAME)
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		return "", err
	}
	content, err := ReadFileLines(file, 0, 655356)
	if err != nil {
		return "", err
	}
	return content, err
}

func GetCurrentAuthCode() (string, error) {
	filePath := GenerateResourceFilePath(AUTH_CODE_FILE_NAME)
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		return "", err
	}
	currentAuthCode, err := ReadFileLines(file, 0, 1)
	if err != nil {
		return currentAuthCode, err
	}
	currentAuthCode = strings.TrimSpace(currentAuthCode)
	return currentAuthCode, nil
}

func NewAuthCode() (string, error) {
	filePath := GenerateResourceFilePath(AUTH_CODE_FILE_NAME)
	file, err := os.OpenFile(filePath, os.O_RDWR, 0)
	defer file.Close()
	if err != nil {
		return "", err
	}
	newAuthCode := GenerateAuthCode()
	var buffer bytes.Buffer
	currentAuthCodes, err := ReadFileLines(file, 0, 655356)
	if err != nil {
		return "", err
	}
	buffer.WriteString(newAuthCode)
	buffer.WriteString("\n")
	buffer.WriteString(currentAuthCodes)
	_, err = file.WriteAt([]byte(buffer.String()), int64(0))
	if err != nil {
		return "", err
	}
	return newAuthCode, nil
}

func ReadFileLines(file *os.File, beginLine uint64, endLine uint64) (string, error) {
	var buffer bytes.Buffer
	r := bufio.NewReaderSize(file, 4096)
	var currentLine uint64
	for currentLine = 0; currentLine < endLine; currentLine++ {
		line, isPrefix, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return buffer.String(), err
			}
		}
		if currentLine < beginLine {
			continue
		}
		buffer.Write(line)
		if !isPrefix {
			buffer.WriteString("\n")
		}
	}
	return buffer.String(), nil
}

func GenerateResourceFilePath(fileName string) string {
	return "resource/" + fileName
}

func GenerateAuthCode() string {
	var limit int64 = 65535
	var buffer bytes.Buffer
	var temp string
	for {
		temp = strconv.FormatInt(rand.Int63n(limit), 16)
		buffer.WriteString(temp)
		if buffer.Len() >= 6 {
			break
		}
	}
	code := buffer.String()
	if len(code) > 6 {
		code = code[:6]
	}
	return code
}
