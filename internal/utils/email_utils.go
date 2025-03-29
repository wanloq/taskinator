package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// Email configuration
const (
	SMTPServer  = "smtp.gmail.com"
	SMTPPort    = "587"
	linkTimeout = 5
)

func ReadSecretFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LoadEmailConfig returns email configurations (SMTPUsername and SMTPPassword) or an error .
func LoadEmailConfig() (string, string, error) {
	SMTPUsername, err := ReadSecretFile("/run/secrets/smtp_username")
	if err != nil {
		return "", "", fmt.Errorf("failed to load SMTP username: %w", err)
	}
	SMTPPassword, err := ReadSecretFile("/run/secrets/smtp_password")
	if err != nil {
		return "", "", fmt.Errorf("failed to load SMTP password: %w", err)
	}
	return SMTPUsername, SMTPPassword, nil
}

// SendPasswordResetEmail sends an email with a password reset link
func SendPasswordResetEmail(toEmail string, resetToken string) error {
	resetLink := fmt.Sprintf("http://0.0.0.0:8080/user/password-reset/confirm?token=%s", resetToken)
	subject := "Taskinator - Password Reset Request"
	note := fmt.Sprintf("Please note that the link expires in %v minutes", linkTimeout)
	body := fmt.Sprintf("Click the link below to reset your Taskinator password:\n\n%s \n\n%s", resetLink, note)

	log.Println("Email Content!\n", body)
	return sendEmail(toEmail, subject, body)
}

// SendVerificationEmail sends an email verification link
func SendVerificationEmail(toEmail string, verificationToken string) error {
	verificationLink := fmt.Sprintf("http://0.0.0.0:8080/user/email/verify?token=%s", verificationToken)
	subject := "Taskinator - Verify Your Email"
	note := fmt.Sprintf("Please note that the link expires in %v minutes", linkTimeout)
	body := fmt.Sprintf("Click the link below to verify your email on Taskinator:\n\n%s\n\n%s", verificationLink, note)

	log.Println("Email Content!\n", body)
	return sendEmail(toEmail, subject, body)
}

// Helper function to send email
func sendEmail(toEmail, subject, body string) error {
	SMTPUsername, SMTPPassword, err := LoadEmailConfig()
	if err != nil {
		log.Fatalf("Error loading SMTP credentials: %v", err)
		return err
	}
	auth := smtp.PlainAuth("", SMTPUsername, SMTPPassword, SMTPServer)
	msg := []byte("Subject: " + subject + "\r\n\r\n" + body)
	log.Println("Email sent!")
	return smtp.SendMail(SMTPServer+":"+SMTPPort, auth, SMTPUsername, []string{toEmail}, msg)
}
