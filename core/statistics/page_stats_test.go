package statistics

import (
	"fmt"
	"math/rand"
	"runtime/debug"
	"testing"
)

func TestPageAccessRecord(t *testing.T) {
	debugTag := true
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	pageName := "page name for test"
	visitor := "guest"
	parameterInfo := fmt.Sprintf("(pageName=%s, visitor=%s)", pageName, visitor)
	if debugTag {
		t.Logf("Testing page access record (%s)...", parameterInfo)
	}
	done, err := AddPageAccessRecord(pageName, visitor, uint64(rand.Int63n(99)))
	if err != nil {
		t.Errorf("Adding page access record error: %s %s\n", err, parameterInfo)
		t.FailNow()
	}
	if !done {
		t.Fatalf("Adding page access record is failing! %s", parameterInfo)
		t.FailNow()
	}
	visitorAccessRecords, err := GetPageAccessRecords(pageName)
	if err != nil {
		t.Errorf("Getting page access record error: %s %s\n", err, parameterInfo)
		t.FailNow()
	}
	expectedRecordLen := 1
	recordLen := len(visitorAccessRecords)
	if recordLen != expectedRecordLen {
		t.Fatalf("The length of visitor access record should be '%s', but it's '%s'.\n", expectedRecordLen, recordLen)
	}
	_, ok := visitorAccessRecords[visitor]
	if !ok {
		t.Fatalf("The access record of visitor '%s' should be existing, but it's nonexistent.\n", visitor)
	}
	done, err = ClearPageAccessRecord(pageName, visitor)
	if err != nil {
		t.Errorf("Clearing page access record error: %s %s\n", err, parameterInfo)
		t.FailNow()
	}
	if !done {
		t.Fatalf("Clearing page access record is failing! %s", parameterInfo)
		t.FailNow()
	}
}
