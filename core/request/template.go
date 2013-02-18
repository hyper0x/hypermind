package request

import (
	"fmt"
	"go_lib"
	"regexp"
)

func SimpleEqual(args ...interface{}) bool {
	if len(args) < 2 {
		return false
	}
	prevValue := ""
	currentValue := ""
	start := false
	result := true
	for _, v := range args {
		currentValue = fmt.Sprintf("%v", v)
		if start && currentValue != prevValue {
			result = false
			break
		}
		prevValue = currentValue
		start = true
	}
	return result
}

func MatchString(args ...interface{}) bool {
	if len(args) < 2 {
		return false
	}
	target := args[0].(string)
	for _, v := range args[1:len(args)] {
		pattern := v.(string)
		pass, err := regexp.MatchString(pattern, target)
		if err != nil {
			go_lib.LogErrorf("RegexpMatchError (target=%s, pattern=%s): %s\n", target, pattern, err)
			return false
		}
		if !pass {
			return false
		}
	}
	return true
}

func AllTrue(literals ...string) bool {
	if len(literals) == 0 {
		return false
	}
	result := true
	for _, v := range literals {
		if v != "true" && v != "y" {
			result = false
			break
		}
	}
	return result
}
