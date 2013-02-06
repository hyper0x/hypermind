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
	CV_FILE_NAME        = "resume.html"
)

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
