package dao

// storage server info
const (
	DEFAULT_SERVER_PORT           = 9090
	DEFAULT_REDIS_SERVER_IP       = "127.0.0.1"
	DEFAULT_REDIS_SERVER_PORT     = "6379"
	DEFAULT_REDIS_SERVER_PASSWORD = ""
)

// redis key
const (
	AUTH_CODE_KEY      = "auth_code"
	USER_KEY_PREFIX    = "hmuser-"
	LOGIN_NAME_FIELD   = "login_name"
	PASSWORD_FIELD     = "password"
	EMAIL_FIELD        = "email"
	MOBILE_PHONE_FIELD = "mobile_phone"
	RIGHTS_FIELD       = "rights"
	REMARK_FIELD       = "remark"
	USER_GROUP_KEY     = "hmusergroup"
)

// user group
const (
	NORMAL_USER_GROUP_NAME = "normal"
	ADMIN_USER_GROUP_NAME  = "admin"
)
