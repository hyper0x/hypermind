package utils

import (
	"fmt"
)

func SimpleEqual(args ...interface{}) bool {
	if len(args) == 0 {
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
