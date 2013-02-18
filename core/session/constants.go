package session

// cookie
const (
	COOKIE_KEY_PREFIX  = "hm-"
	SESSION_COOKIE_KEY = "sessionid"
	COOKIE_MAX_AGE     = 60 * 60 * 2 // cookie live seconds
)

// session
const (
	SESSION_KEY_PREFIX      = "hmsession-"
	SESSION_GRANTORS_KEY    = "grantors"
	SESSION_GROUP_KEY       = "group"
	SESSION_COOKIE_SIGN_KEY = "cookie_sign"
)
