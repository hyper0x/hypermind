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
	done, err := SetPageAccessRecord(pageName, visitor, uint64(rand.Int63n(99)))
	if err != nil {
		t.Errorf("Setting page access record error: %s %s\n", err, parameterInfo)
		t.FailNow()
	}
	if !done {
		t.Logf("Setting page access record is failing! %s", parameterInfo)
		t.FailNow()
	}
	done, err = ClearPageAccessRecord(pageName, visitor)
	if err != nil {
		t.Errorf("Clearing page access record error: %s %s\n", err, parameterInfo)
		t.FailNow()
	}
	if !done {
		t.Logf("Clearing page access record is failing! %s", parameterInfo)
		t.FailNow()
	}
}
