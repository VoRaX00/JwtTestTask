package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/smtp"
	"os"
)

func NewErrorResponse(c *gin.Context, code int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(code, gin.H{
		message: message,
	})
}

func SendMessageEmail(email, message string) error {
	if email == "" {
		return errors.New("email is empty")
	}

	fromEmail := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	msg := []byte(message)

	auth := smtp.PlainAuth("", fromEmail, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, []string{email}, msg)
	if err != nil {
		return err
	}
	return nil
}
