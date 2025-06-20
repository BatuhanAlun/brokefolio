package utils

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

func SendEmail(recipientEmail, subject, body string) error {
	smtpUsername := "brokefolio@hotmail.com" // Replace with your actual email
	smtpPassword := "securepass123456"       // Replace with your actual password
	smtpHost := "smtp-mail.outlook.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// Connect to the SMTP server.
	conn, err := net.Dial("tcp", smtpHost+":"+smtpPort)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Create a new SMTP client.
	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer c.Close()

	// Check if the server supports STARTTLS.
	if ok, _ := c.Extension("STARTTLS"); !ok {
		return fmt.Errorf("SMTP server does not support STARTTLS")
	}

	// Start TLS negotiation.
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false, // Verify the server's certificate
		ServerName:         smtpHost,
	}
	if err = c.StartTLS(tlsconfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// Re-authenticate after TLS is established (some servers require this).
	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate after STARTTLS: %w", err)
	}

	// Set the sender and recipient.
	if err = c.Mail(smtpUsername); err != nil {
		return fmt.Errorf("failed to set sender address: %w", err)
	}
	if err = c.Rcpt(recipientEmail); err != nil {
		return fmt.Errorf("failed to set recipient address: %w", err)
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		return fmt.Errorf("failed to open data stream: %w", err)
	}
	defer wc.Close()

	_, err = fmt.Fprintf(wc, "To: %s\r\n", recipientEmail)
	if err != nil {
		return fmt.Errorf("failed to write 'To' header: %w", err)
	}
	_, err = fmt.Fprintf(wc, "Subject: %s\r\n", subject)
	if err != nil {
		return fmt.Errorf("failed to write 'Subject' header: %w", err)
	}
	_, err = fmt.Fprintf(wc, "\r\n%s\r\n", body)
	if err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}

	if err = wc.Close(); err != nil {
		return fmt.Errorf("failed to close data stream: %w", err)
	}

	// Quit the SMTP connection.
	if err = c.Quit(); err != nil {
		return fmt.Errorf("failed to quit SMTP connection: %w", err)
	}

	return nil
}
