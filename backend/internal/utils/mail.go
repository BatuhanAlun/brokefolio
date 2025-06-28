package utils

import (
	"log"
	"net/smtp"
	"os"
)

func SendEmail(recipientEmail, subject, body string) error {

	appKey := os.Getenv("MAIL_APP_KEY")

	if appKey == "" {
		log.Fatalf("APP KEY ENV NOT FOUND!")
	}

	auth := smtp.PlainAuth("", "brokefolio@gmail.com", appKey, "smtp.gmail.com")

	to := []string{recipientEmail}
	msg := []byte(
		"From: brokefolio@gmail.com\r\n" +
			"To: " + recipientEmail + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n" +
			"\r\n" +
			body + "\r\n")

	err := smtp.SendMail("smtp.gmail.com:587", auth, "brokefolio@gmail.com", to, msg)
	if err != nil {
		return err
	}
	return nil
}
