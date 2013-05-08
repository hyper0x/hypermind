package rights

import (
	"encoding/json"
	"errors"
	"hypermind/core/base"
	"hypermind/core/dao"
	"runtime/debug"
)

type GroupRights struct {
	PageRights map[string]string
}

type UserGroup struct {
	Name   string
	Rights GroupRights
}

var guestGroupRightsDict map[string]string = map[string]string{
	base.HOME_PAGE:              "true",
	base.ABOUT_ME_PAGE:          "true",
	base.ABOUT_WEBSITE_PAGE:     "true",
	base.MEETING_KANBAN_PAGE:    "false",
	base.PROJECT_HASH_RING_PAGE: "true",
	base.ADMIN_AUTH_CODE_PAGE:   "false",
	base.ADMIN_USER_LIST_PAGE:   "false",
}

var normalGroupRightsDict map[string]string = map[string]string{
	base.HOME_PAGE:              "true",
	base.ABOUT_ME_PAGE:          "true",
	base.ABOUT_WEBSITE_PAGE:     "true",
	base.MEETING_KANBAN_PAGE:    "true",
	base.PROJECT_HASH_RING_PAGE: "true",
	base.ADMIN_AUTH_CODE_PAGE:   "false",
	base.ADMIN_USER_LIST_PAGE:   "false",
}

var adminGroupRightsDict map[string]string = map[string]string{
	base.HOME_PAGE:              "true",
	base.ABOUT_ME_PAGE:          "true",
	base.ABOUT_WEBSITE_PAGE:     "true",
	base.MEETING_KANBAN_PAGE:    "true",
	base.PROJECT_HASH_RING_PAGE: "true",
	base.ADMIN_AUTH_CODE_PAGE:   "true",
	base.ADMIN_USER_LIST_PAGE:   "true",
}

var userGroupMap map[string]GroupRights = map[string]GroupRights{
	NORMAL_USER_GROUP_NAME: GroupRights{PageRights: normalGroupRightsDict},
	ADMIN_USER_GROUP_NAME:  GroupRights{PageRights: adminGroupRightsDict},
}

func init() {
	for name, rights := range userGroupMap {
		userGroup, err := GetUserGroup(name)
		if err != nil {
			base.Logger().Errorf("Get User Group (name=%s) Error: %s\n", name, err)
			continue
		}
		if userGroup != nil {
			err = DeleteUserGroup(name)
			if err != nil {
				base.Logger().Errorf("Delete User Group (name=%s) Error: %s\n", name, err)
				continue
			}
		}
		userGroup = &UserGroup{Name: name, Rights: rights}
		err = AddUserGroup(userGroup)
		if err != nil {
			base.Logger().Errorf("Add User Group '%v' Error: %s\n", userGroup, err)
			continue
		}
	}
}

func GetGuestPageRights() map[string]string {
	copy := make(map[string]string)
	for k, v := range guestGroupRightsDict {
		copy[k] = v
	}
	return copy
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
	if userGroup == nil || userGroup.Name == "" || len(userGroup.Rights.PageRights) == 0 {
		return errors.New("The parameter named userGroup is NOT Ready!")
	}
	groupRightsLiterals, err := MarshalGroupRights(userGroup.Rights)
	if err != nil {
		return err
	}
	conn := dao.RedisPool.Get()
	defer conn.Close()
	_, err = dao.SetHash(dao.USER_GROUP_KEY, userGroup.Name, groupRightsLiterals)
	if err != nil {
		return err
	}
	return nil
}

func GetUserGroup(groupName string) (*UserGroup, error) {
	if len(groupName) == 0 {
		debug.PrintStack()
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
	_, err := dao.DelHashField(dao.USER_GROUP_KEY, groupName)
	if err != nil {
		return err
	}
	return nil
}
