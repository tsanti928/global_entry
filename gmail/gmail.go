package gmail

import (
	"errors"
	"net/smtp"
	"os"
	"strings"
)

// SendEmail sends an email using smtp.gmail.com.
func SendEmail(subject string, body []string, to []string) error {
	user := os.Getenv("GE_USER")
	if user == "" {
		return errors.New("GE_USER environment variable not set")
	}
	appPW := os.Getenv("GE_APP_PW")
	if appPW == "" {
		return errors.New("GE_APP_PW environment variable not set")
	}
	if len(to) == 0 || to[0] == "" {
		return errors.New("To: users not set")
	}

	auth := smtp.PlainAuth("", user, appPW, "smtp.gmail.com")
	msg := make([]string, 0, 1+len(body))
	msg = append(msg, "SUBJECT: "+subject)
	msg = append(msg, body...)
	asBytes := []byte(strings.Join(msg, "\n"))

	return smtp.SendMail("smtp.gmail.com:587", auth, user, to, asBytes)
}
