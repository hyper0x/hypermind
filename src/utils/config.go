package utils

import (
	"os"
	"bufio"
	"strings"
	"io"
	"strconv"
	"fmt"
	"bytes"
)

const (
	DefaultServerPort = 9090
	ConfigFilePath = "go-web-demo.config"
)

type MyConfig struct {
    ServerPort int
	WorkDir string
    Extras map[string]string
}

var myConfig MyConfig

var loaded bool = false

func ReadConfig(fresh bool) (MyConfig, error) {
	defer func() {
		LogInfof("Loaded config: %v\n", myConfig)
	}()
	if !fresh && loaded {
		return myConfig, nil
	}
	myConfig = *new(MyConfig)
	myConfig.ServerPort = DefaultServerPort
	currentDir, err := os.Getwd()
	if err != nil {
		LogErrorln("GetwdError:", err)
	} else {
		myConfig.WorkDir = currentDir
	}
	myConfig.Extras = make(map[string]string)
	configFile, err := os.OpenFile(ConfigFilePath, os.O_RDONLY, 0666)
	if err != nil {
		switch err.(type) {
		case *os.PathError:
			var warningBuffer bytes.Buffer
			warningBuffer.WriteString("Warning: the config file '")
			warningBuffer.WriteString(ConfigFilePath)
			warningBuffer.WriteString("' is NOT FOUND! ")
			warningBuffer.WriteString("Use DEFAULT config '")
			warningBuffer.WriteString(fmt.Sprintf("%v", myConfig))
			warningBuffer.WriteString("'. ")
			LogWarnln(warningBuffer.String())
			loaded = true
			return myConfig, nil
		default:
			return myConfig, err
		}
	}
	defer configFile.Close()
	configReader := bufio.NewReader(configFile)
	for {
		str, err := configReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// The file end is touched.
				break
			} else {
				return myConfig, err
			}
		}
		str = strings.TrimRight(str, "\r\n")
		if len(str) == 0 {
			continue
		}
		sepIndex := strings.Index(str, "=")
		if sepIndex <= 0 {
			continue
		}
		key := str[0:sepIndex]
		value := str[sepIndex + 1:len(str)]
		switch strings.ToLower(key) {
		case "server_port":
			portNumber , err := strconv.Atoi(value)
			if err == nil {
				myConfig.ServerPort = portNumber
			}
		case "current_dir":
			myConfig.WorkDir = value
		default:
			myConfig.Extras[key] = value
		}
	}
	loaded = true
	return myConfig, nil
}

