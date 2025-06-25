package utils

import (
	"net/smtp"
)

func SendEmail(recipientEmail, subject, body string) error {

	auth := smtp.PlainAuth("", "brokefolio@gmail.com", "osxu kppw anwn ofsc", "smtp.gmail.com")

	to := []string{"batuhanalun1999@hotmail.com"}
	msg := []byte("To: batuhanalun1999@hotmail.com\r\n" +
		"Subject: Hello from Go\r\n" +
		"\r\n" +
		"This is a test email.\r\n")

	err := smtp.SendMail("smtp.gmail.com:587", auth, "brokefolio@gmail.com", to, msg)
	if err != nil {
		return err
	}
	return nil
}
