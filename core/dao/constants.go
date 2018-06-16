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
	AUTH_CODE_KEY           = "auth_code"
	USER_KEY_PREFIX         = "hmuser-"
	USER_GROUP_KEY          = "hmusergroup"
	PAGE_ACCESS_RECORDS_KEY = "page_access_records"
)
