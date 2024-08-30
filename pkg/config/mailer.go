package config

import (
	"github.com/ecommerce-api/pkg/helper"
	mail "github.com/xhit/go-simple-mail/v2"
	"time"
)

var (
	MailHost           string
	MailPort           string
	MailUsername       string
	MailPassword       string
	MailFrom           string
	MailEncryption     = mail.EncryptionTLS
	MailConnectTimeout = 10 * time.Second
	MailSendTimeout    = 10 * time.Second
	MailKeepAlive      = false
	MailPriority       = mail.PriorityLow
)

func MailerLoad() {
	MailHost = helper.GetEnv("MAIL_HOST", "smtp")
	MailPort = helper.GetEnv("MAIL_PORT", "578")
	MailUsername = helper.GetEnv("MAIL_USERNAME", "username")
	MailPassword = helper.GetEnv("MAIL_PASSWORD", "password")
	MailFrom = helper.GetEnv("MAIL_FROM", "username@mail.com")
}
