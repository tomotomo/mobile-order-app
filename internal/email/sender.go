package email

import (
	"fmt"
	"log"
	"net/smtp"
)

type Sender interface {
	Send(to string, subject string, body string) error
}

// LogSender is the "Trap". It prints emails to stdout instead of sending.
type LogSender struct{}

func (s *LogSender) Send(to string, subject string, body string) error {
	log.Println("========== EMAIL TRAP ==========")
	log.Printf("To: %s\n", to)
	log.Printf("Subject: %s\n", subject)
	log.Printf("Body:\n%s\n", body)
	log.Println("================================")
	return nil
}

// SmtpSender sends actual emails.
type SmtpSender struct {
	Host     string
	Port     string
	Username string
	Password string
}

func (s *SmtpSender) Send(to string, subject string, body string) error {
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, subject, body))
	
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	if err := smtp.SendMail(addr, auth, s.Username, []string{to}, msg); err != nil {
		return err
	}
	log.Printf("Email sent to %s", to)
	return nil
}
