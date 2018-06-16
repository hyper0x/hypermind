package rights

import (
	"runtime/debug"
	"testing"
)

func TestGroupRights(t *testing.T) {
	debugTag := true
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	rightsDict := make(map[string]string)
	rightsDict["1"] = "A"
	rightsDict["2"] = "B"
	rightsDict["3"] = "C"
	groupRights := GroupRights{PageRights: rightsDict}
	groupRightsLiterals, err := MarshalGroupRights(groupRights)
	if err != nil {
		t.Errorf("Error: Marshal Group Rights Error: %s\n", err)
		t.FailNow()
	}
	expectGroupRightsLiterals := "{\"PageRights\":{\"1\":\"A\",\"2\":\"B\",\"3\":\"C\"}}"
	if groupRightsLiterals != expectGroupRightsLiterals {
		t.Errorf("Fail: The group rights literals should be %s but %s. \n", expectGroupRightsLiterals, groupRightsLiterals)
		t.FailNow()
	}
	if debugTag {
		t.Logf("Group rights literals: %s \n", groupRightsLiterals)
	}
	groupRightsCopy, err := UnmarshalGroupRights(groupRightsLiterals)
	if err != nil {
		t.Errorf("Error: Unmarshal Group Rights Error: %s\n", err)
		t.FailNow()
	}
	groupRightsCopyLen := len(groupRightsCopy.PageRights)
	groupRightsLen := len(groupRights.PageRights)
	if groupRightsCopyLen != groupRightsLen {
		t.Errorf("Fail: The length of group rights dict copy should be %d but %d. \n\n", groupRightsLen, groupRightsCopyLen)
		t.FailNow()
	}
	if debugTag {
		t.Logf("The length of group rights literals copy is %d. \n", groupRightsCopyLen)
	}
	for k, v := range groupRightsCopy.PageRights {
		ev := groupRights.PageRights[k]
		if len(v) == 0 || v != ev {
			t.Errorf("Fail: The value of key '%s' in group rights copy should be %s but %s. \n\n", k, ev, v)
			t.FailNow()
		}
	}
	if debugTag {
		t.Logf("Group rights copy: %v \n", groupRightsCopy)
	}
}

func TestUserGroup(t *testing.T) {
	debugTag := true
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	rightsDict := make(map[string]string)
	rightsDict["1"] = "A"
	rightsDict["2"] = "B"
	rightsDict["3"] = "C"
	groupRights := GroupRights{PageRights: rightsDict}
	groupName := "testing_group"
	userGroup0 := &UserGroup{
		Name:   groupName,
		Rights: groupRights}
	if debugTag {
		t.Logf("UserGroup0: %v\n", userGroup0)
	}
	err := AddUserGroup(userGroup0)
	if err != nil {
		t.Errorf("Error: Add User Group Error: %s\n", err)
		t.FailNow()
	}
	userGroup1, err := GetUserGroup(groupName)
	if err != nil {
		t.Errorf("Error: Get User Group Error: %s\n", err)
		t.FailNow()
	}
	if debugTag {
		t.Logf("UserGroup1: %v\n", userGroup1)
	}
	userGroup0RightsLiterals, err := MarshalGroupRights(userGroup0.Rights)
	if err != nil {
		t.Errorf("Error: Marshal Group Rights Error: %s\n", err)
		t.FailNow()
	}
	userGroup1RightsLiterals, err := MarshalGroupRights(userGroup1.Rights)
	if err != nil {
		t.Errorf("Error: Marshal Group Rights Error: %s\n", err)
		t.FailNow()
	}
	if userGroup1.Name != groupName ||
		userGroup1RightsLiterals != userGroup0RightsLiterals {
		t.Errorf("Fail: The user group should be %v but %v. \n", userGroup0, userGroup1)
		t.FailNow()
	}
	err = DeleteUserGroup(groupName)
	if err != nil {
		t.Errorf("Error: Delete User Group Error: %s\n", err)
		t.FailNow()
	}
	userGroup3, err := GetUserGroup(groupName)
	if err != nil {
		t.Errorf("Error: Get User Group Error: %s\n", err)
		t.FailNow()
	}
	if userGroup3 != nil {
		t.Errorf("Fail: The user group '%v' should be deleted. %s\n", userGroup3)
		t.FailNow()
	}
}
