package statistics

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_lib"
	"hypermind/core/dao"
)

var signMap map[string]*go_lib.Sign = make(map[string]*go_lib.Sign)

func AddPageAccessRecord(pageName string, visitor string, number uint64) (bool, error) {
	if len(pageName) == 0 {
		return false, errors.New("The parameter named pageName is EMPTY!")
	}
	sign := getSignForPage(pageName)
	sign.Set()
	defer sign.Unset()
	parameterInfo := fmt.Sprintf("(pageName=%s, visitor=%s, number=%d)", pageName, visitor, number)
	var result bool
	conn := dao.RedisPool.Get()
	defer conn.Close()
	value, err := dao.GetHash(dao.PAGE_ACCESS_RECORDS_KEY, pageName)
	if err != nil {
		return false, err
	}
	var visitorAccessRecords map[string]uint64
	if len(value) > 0 {
		visitorAccessRecords, err = parseVisitorAccessRecords(value)
		if err != nil {
			go_lib.LogErrorf("Parsing visitor access records error: %s %s\n", err, parameterInfo)
		}
	}
	if visitorAccessRecords == nil {
		visitorAccessRecords = make(map[string]uint64)
	}
	visitorAccessRecords[visitor] = visitorAccessRecords[visitor] + uint64(number)
	literals, err := formatVisitorAccessRecords(visitorAccessRecords)
	if err != nil {
		go_lib.LogErrorf("Formating visitor access records error: %s %s\n", err, parameterInfo)
	}
	if len(literals) > 0 {
		result, err = dao.SetHash(dao.PAGE_ACCESS_RECORDS_KEY, pageName, literals)
		if err != nil {
			return result, err
		}
	}
	if result {
		go_lib.LogInfof("The page access info has been recorded. %s\n", parameterInfo)
	} else {
		go_lib.LogWarnf("The page access info failed to record. %s\n", parameterInfo)
	}
	return result, nil
}

func ClearPageAccessRecord(pageName string, visitor string) (bool, error) {
	if len(pageName) == 0 {
		return false, errors.New("The parameter named pageName is EMPTY!")
	}
	sign := getSignForPage(pageName)
	sign.Set()
	defer sign.Unset()
	parameterInfo := fmt.Sprintf("(pageName=%s, visitor=%s)", pageName, visitor)
	var result bool
	conn := dao.RedisPool.Get()
	defer conn.Close()
	value, err := dao.GetHash(dao.PAGE_ACCESS_RECORDS_KEY, pageName)
	if err != nil {
		return false, err
	}
	visitorAccessRecords, err := parseVisitorAccessRecords(value)
	if err != nil {
		go_lib.LogErrorf("Parsing visitor access records error: %s %s\n", err, parameterInfo)
	}
	if visitorAccessRecords != nil {
		_, ok := visitorAccessRecords[visitor]
		if ok {
			delete(visitorAccessRecords, visitor)
			var literals string
			validliterals := true
			if len(visitorAccessRecords) > 0 {
				literals, err = formatVisitorAccessRecords(visitorAccessRecords)
				if err != nil {
					go_lib.LogErrorf("Formating visitor access records error: %s %s\n", err, parameterInfo)
					validliterals = false
				}
			}
			if validliterals {
				result, err = dao.SetHash(dao.PAGE_ACCESS_RECORDS_KEY, pageName, literals)
				if err != nil {
					return false, err
				}
			}
		}
	}
	if result {
		go_lib.LogInfof("The page access info has been cleared. %s\n", parameterInfo)
	} else {
		go_lib.LogWarnf("The page access info failed to clear. %s\n", parameterInfo)
	}
	return result, nil
}

func parseVisitorAccessRecords(literals string) (map[string]uint64, error) {
	if len(literals) == 0 {
		errorMsg := fmt.Sprintf("The parameter named literals is EMPTY! IGNORE the unmarshal operation.")
		return nil, errors.New(errorMsg)
	}
	var records map[string]uint64
	err := json.Unmarshal([]byte(literals), &records)
	if err != nil {
		errorMsg := fmt.Sprintf("Json unmarshal error (source=%v): %s\n", literals, err)
		return nil, errors.New(errorMsg)
	}
	return records, nil
}

func formatVisitorAccessRecords(records map[string]uint64) (string, error) {
	if len(records) == 0 {
		errorMsg := fmt.Sprintf("The parameter named records is EMPTY! IGNORE the marshal operation.")
		return "", errors.New(errorMsg)
	}
	var literals string
	b, err := json.Marshal(records)
	if err != nil {
		errorMsg := fmt.Sprintf("Json marshal error (source=%v): %s\n", records, err)
		return "", errors.New(errorMsg)
	}
	literals = string(b)
	return literals, nil
}

func getSignForPage(pageName string) *go_lib.Sign {
	sign := signMap[pageName]
	if sign == nil {
		sign = go_lib.NewSign()
		signMap[pageName] = sign
	}
	return sign
}
