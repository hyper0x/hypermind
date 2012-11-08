package myutil

import (
	"os"
	"bufio"
	"strings"
	"io"
	"strconv"
	"fmt"
)

const (
	DefaultServerPort = 9090
	ConfigFilePath = "go-web-demo.config"
)

type MyConfig struct {
    ServerPort int
    Extras map[string]string
}

var myConfig MyConfig

var loaded bool = false

func ReadConfig(fresh bool) (MyConfig, error) {
	defer func() {
		fmt.Printf("Loaded config: %s\n", myConfig)
	}()
	if !fresh && loaded {
		return myConfig, nil
	}
	configFile, err := os.OpenFile(ConfigFilePath, os.O_RDONLY, 0666)
	if err != nil {
		return myConfig, err
	}
	defer configFile.Close()
	configReader := bufio.NewReader(configFile)
	myConfig = *new(MyConfig)
	myConfig.ServerPort = DefaultServerPort
	myConfig.Extras = make(map[string]string)
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
		default:
			myConfig.Extras[key] = value
		}
	}
	loaded = true
	return myConfig, nil
}

