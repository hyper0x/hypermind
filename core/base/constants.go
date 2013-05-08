package base

import (
	"go_lib/logging"
)

var logger logging.Logger = logging.GetSimpleLogger()

func Logger() logging.Logger {
	return logger
}

// config
const (
	CONFIG_FILE_NAME = "hypermind.config"
)

// page parameter
const (
	HOME_PAGE_KEY           = "homePage"
	HOME_PAGE               = "home"
	ABOUT_ME_PAGE_KEY       = "aboutMePage"
	ABOUT_ME_PAGE           = "about_me"
	ABOUT_WEBSITE_PAGE_KEY  = "aboutWebsitePage"
	ABOUT_WEBSITE_PAGE      = "about_website"
	MEETING_KANBAN_PAGE_KEY = "meetingKanbanPage"
	MEETING_KANBAN_PAGE     = "meeting_kanban"
	PROJECT_HASH_RING_KEY   = "projectHashRingPage"
	PROJECT_HASH_RING_PAGE  = "project_hash_ring"
	ADMIN_AUTH_CODE_KEY     = "adminAuthCodePage"
	ADMIN_AUTH_CODE_PAGE    = "admin_auth_code"
	ADMIN_USER_LIST_KEY     = "adminUserListPage"
	ADMIN_USER_LIST_PAGE    = "admin_user_list"
)
