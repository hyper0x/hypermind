package base

import (
	"go_lib"
)

var hmConfig go_lib.Config

func init() {
	hmConfig = go_lib.Config{Path: CONFIG_FILE_NAME}
	err := hmConfig.ReadConfig(false)
	if err != nil {
		Logger().Errorln("ConfigLoadError: ", err)
	}
}

func GetHmConfig() go_lib.Config {
	return hmConfig
}
