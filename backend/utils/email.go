package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to string, subject string, body string) error {
	from := os.Getenv("GMAIL_ADDRESS")
	password := os.Getenv("GMAIL_APP_PASSWORD")
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	message := []byte("Subject: " + subject + "\r\n" + "\r\n" + body)

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}

func SendVerificationEmail(userEmail, token string) error {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "https://goanddocker.onrender.com"
	}
	verificationURL := fmt.Sprintf("%s/api/verify-email?token=%s", baseURL, token)
	subject := "Please verify your email address"
	body := fmt.Sprintf("Click this link to verify your email address: %s", verificationURL)
	return SendEmail(userEmail, subject, body)
}
