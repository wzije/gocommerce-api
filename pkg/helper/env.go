package helper

import (
	"os"
)

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func IsEnv(env string) bool {
	return os.Getenv("APP_ENV") == env
}
