package config

import "time"

//FS mean file service
const (
	FSRedisAddr = "localhost:6379"
	FSRedisPassword = ""
	FSRedisDB = 0
)

//Auth mean Auth service
const (
	AuthRedisAddr = "localhost:6379"
	AuthRedisPassword = ""
	AuthRedisDB = 1
	AuthRedisTokenExpireTime = 72 * time.Hour
)
