package rights

import (
	"encoding/json"
	"errors"
	"go_lib"
	"hypermind/core/dao"
	"hypermind/core/request"
)

type GroupRights struct {
	Dict map[string]string
}

type UserGroup struct {
	Name   string
	Rights GroupRights
}

var normalGroupRightsDict map[string]string = map[string]string{
	request.HOME_PAGE:              "true",
	request.ABOUT_ME_PAGE:          "true",
	request.ABOUT_WEBSITE_PAGE:     "true",
	request.MEETING_KANBAN_PAGE:    "true",
	request.PROJECT_HASH_RING_PAGE: "true",
	request.ADMIN_CV_PAGE:          "false",
}

var adminGroupRightsDict map[string]string = map[string]string{
	request.HOME_PAGE:              "true",
	request.ABOUT_ME_PAGE:          "true",
	request.ABOUT_WEBSITE_PAGE:     "true",
	request.MEETING_KANBAN_PAGE:    "true",
	request.PROJECT_HASH_RING_PAGE: "true",
	request.ADMIN_CV_PAGE:          "true",
}

var userGroupMap map[string]GroupRights = map[string]GroupRights{
	NORMAL_USER_GROUP_NAME: GroupRights{Dict: normalGroupRightsDict},
	ADMIN_USER_GROUP_NAME:  GroupRights{Dict: adminGroupRightsDict},
}

func init() {
	for name, rights := range userGroupMap {
		userGroup, err := GetUserGroup(name)
		if err != nil {
			go_lib.LogErrorf("Get User Group (name=%s) Error: %s\n", name, err)
			continue
		}
		if userGroup != nil {
			err = DeleteUserGroup(name)
			if err != nil {
				go_lib.LogErrorf("Delete User Group (name=%s) Error: %s\n", name, err)
				continue
			}
		}
		userGroup = &UserGroup{Name: name, Rights: rights}
		err = AddUserGroup(userGroup)
		if err != nil {
			go_lib.LogErrorf("Add User Group '%v' Error: %s\n", userGroup, err)
			continue
		}
	}
}

func UnmarshalGroupRights(literals string) (GroupRights, error) {
	var rights GroupRights
	err := json.Unmarshal([]byte(literals), &rights)
	if err != nil {
		return rights, err
	}
	return rights, nil
}

func MarshalGroupRights(rights GroupRights) (string, error) {
	b, err := json.Marshal(rights)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func AddUserGroup(userGroup *UserGroup) error {
	if userGroup == nil || userGroup.Name == "" || len(userGroup.Rights.Dict) == 0 {
		return errors.New("The parameter named userGroup is NOT Ready!")
	}
	groupRightsLiterals, err := MarshalGroupRights(userGroup.Rights)
	if err != nil {
		return err
	}
	conn := dao.RedisPool.Get()
	defer conn.Close()
	err = dao.SetHash(dao.USER_GROUP_KEY, userGroup.Name, groupRightsLiterals)
	if err != nil {
		return err
	}
	return nil
}

func GetUserGroup(groupName string) (*UserGroup, error) {
	if len(groupName) == 0 {
		return nil, errors.New("The parameter named groupName is EMPTY!")
	}
	groupRightsLiterals, err := dao.GetHash(dao.USER_GROUP_KEY, groupName)
	if err != nil {
		return nil, err
	}
	if len(groupRightsLiterals) == 0 {
		return nil, nil
	}
	groupRights, err := UnmarshalGroupRights(groupRightsLiterals)
	if err != nil {
		return nil, err
	}
	userGroup := new(UserGroup)
	userGroup.Name = groupName
	userGroup.Rights = groupRights
	return userGroup, nil
}

func DeleteUserGroup(groupName string) error {
	if len(groupName) == 0 {
		return errors.New("The parameter named groupName is EMPTY!")
	}
	err := dao.DelHashField(dao.USER_GROUP_KEY, groupName)
	if err != nil {
		return err
	}
	return nil
}
