package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ecommerce-api/pkg/helper"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	_, b, _, _ = runtime.Caller(0)
	BasePath   = filepath.Join(filepath.Dir(b), "../../") + "/"
	Config     config
)

type config struct {
	AppEnv   string
	AppName  string
	AppPort  string
	AppUrl   string
	WebUrl   string
	CertPath string
	LogPath  string
}

func Load(envs ...string) {

	var envFile = BasePath + ".env"
	if os.Getenv("APP_ENV") == "development" {
		for _, env := range envs {
			if strings.ToLower(env) == "test" {
				envFile = BasePath + ".env.test"
			}
		}
	}

	if err := godotenv.Load(envFile); err != nil {
		logrus.Fatal(err)
	}

	//SET VARIABLE
	Config = config{
		AppEnv:   helper.GetEnv("APP_ENV", "local"),
		AppPort:  helper.GetEnv("APP_PORT", "9090"),
		AppName:  helper.GetEnv("APP_NAME", "APP NAME"),
		AppUrl:   helper.GetEnv("APP_URL", "https://ecommorce-api.test"),
		WebUrl:   helper.GetEnv("WEB_URL", "https://ecommorce-api.test"),
		CertPath: BasePath + "var/cert.key",
		LogPath:  BasePath + "var/logs/logger.log",
	}

	//LOAD CONFIG
	LoggerLoad()
	SqlDBLoad()
	RedisLoad()
	MailerLoad()

}
