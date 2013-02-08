package utils

// config
const (
	CONFIG_FILE_NAME              = "hypermind.config"
	DEFAULT_SERVER_PORT           = 9090
	DEFAULT_REDIS_SERVER_IP       = "127.0.0.1"
	DEFAULT_REDIS_SERVER_PORT     = "6379"
	DEFAULT_REDIS_SERVER_PASSWORD = ""
)

// web
const (
	COOKIE_LIFE_CYCLE_MINUTES int = 60
)

// request parameter
const (
	LOGIN_NAME_KEY   string = "loginName"
	PASSWORD_KEY     string = "password"
	CN_NAME_KEY      string = "cnName"
	EMAIL_KEY        string = "email"
	MOBILE_PHONE_KEY string = "mobilePhone"
	ROOT_USER_NAME          = "root"
	AUTH_CODE               = "auth_code"
)

// page parameter
const (
	HOME_PAGE_KEY           = "homePage"
	HOME_PAGE               = "home"
	ABOUT_ME_PAGE_KEY       = "aboutMePage"
	ABOUT_ME_PAGE           = "about-me"
	ABOUT_WEBSITE_PAGE_KEY  = "aboutWebsitePage"
	ABOUT_WEBSITE_PAGE      = "about-website"
	MEETING_KANBAN_PAGE_KEY = "meetingKanbanPage"
	MEETING_KANBAN_PAGE     = "meeting-kanban"
	PROJECT_HASH_RING_KEY   = "projectHashRingPage"
	PROJECTS_HASH_RING_PAGE = "project-hash-ring"
)

// redis
const (
	USER_KEY      = "user"
	AUTH_CODE_KEY = "auth_code"
)
