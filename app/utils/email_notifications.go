package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Notifier interface {
	Send(to, subject, body string) error
}

type EmailNotifier struct{}

func (emailNotifier EmailNotifier) Send(to, subject, body string) error {
	var goMail *gomail.Message = gomail.NewMessage()

	from := GetEnv("MAIL_FROM_ADDRESS")
	if from == "" {
		return fmt.Errorf("MAIL_FROM_ADDRESS no está configurado")
	}

	// Set email headers
	goMail.SetHeader("From", from)
	goMail.SetHeader("To", to)
	goMail.SetHeader("Subject", subject)

	// Set email body
	goMail.SetBody("text/html", body)

	// Convertir el puerto a int
	port, err := strconv.Atoi(GetEnv("MAIL_PORT"))
	if err != nil {
		return fmt.Errorf("MAIL_PORT inválido: %w", err)
	}

	// Set up SMTP server configuration
	var dialer *gomail.Dialer = gomail.NewDialer(
		GetEnv("MAIL_HOST"),
		port,
		GetEnv("MAIL_USERNAME"),
		GetEnv("MAIL_PASSWORD"),
	)

	var sendErr error = dialer.DialAndSend(goMail)
	if sendErr != nil {
		return fmt.Errorf("error sending email: %w", sendErr)
	}

	return nil
}

// generateResetToken genera un hash simple para recuperación
func GenerateResetToken(email string) string {
	h := sha256.New()
	h.Write([]byte(email + ":reset"))
	return hex.EncodeToString(h.Sum(nil))
}
