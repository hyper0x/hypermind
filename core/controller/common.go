package controller

import (
	"bufio"
	"hypermind/core/base"
	"hypermind/core/statistics"
)

func pushResponse(bufrw *bufio.ReadWriter, authCode string) bool {
	_, err := bufrw.Write([]byte(authCode))
	if err == nil {
		err = bufrw.Flush()
	}
	if err != nil {
		base.Logger().Errorf("PushAuthCodeError: %s\n", err)
		return false
	}
	return true
}

func recordPageAccessInfo(pageName string, visitor string, number uint64) bool {
	var result bool
	done, err := statistics.AddPageAccessRecord(pageName, visitor, number)
	if err != nil {
		base.Logger().Errorf("Adding page access record error: %s (pageName=%s, visitor=%s, number=%d)\n", err, pageName, visitor, number)
		result = false
	} else {
		result = done
	}
	return result
}
