package utils

import (
	"go_lib"
)

var myConfig go_lib.Config

func init() {
	myConfig = go_lib.Config{Path : CONFIG_FILE_NAME}
	err := myConfig.ReadConfig(false)
	if err != nil {
		go_lib.LogErrorln("ConfigLoadError: ", err)
	}
}
