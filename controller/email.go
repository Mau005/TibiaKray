package controller

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"sync"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/jordan-wright/email"
)

var (
	mu sync.Mutex // Mutex para sincronización
)

type EmailController struct{}

func (ec *EmailController) GenerateEmailSend(subject, body string) *email.Email {
	e := &email.Email{
		Subject: subject,
		HTML:    []byte(body),
	}
	return e
}

func (ec *EmailController) SendEmail(userEmailSend string, e *email.Email) {
	mu.Lock()
	defer mu.Unlock()
	mail := configuration.Config.Email

	auth := smtp.PlainAuth("", mail.SmtpUsername, mail.Password, mail.SmtpServer)

	e.From = mail.SmtpUsername
	e.To = []string{userEmailSend}
	err := e.SendWithTLS(fmt.Sprintf("%s:%s", mail.SmtpServer, mail.SmptPort), auth, &tls.Config{
		InsecureSkipVerify: true, // Solo si no tienes un certificado SSL/TLS válido
		ServerName:         mail.SmtpServer,
	})
	if err != nil {
		fmt.Println("Error sending email:", err)
	} else {
		fmt.Println("Email sent successfully to:", userEmailSend)
	}
}
