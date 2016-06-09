package communication

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-gomail/gomail"
)

//EmailService defines an email communication channel
type EmailService interface {
	Send(email string, subject string, message string) (err error)
}

//DevEmailService is the implementation of an EmailService suitable for use in local development environments
type DevEmailService struct{}

//Send sends an Email
func (s *DevEmailService) Send(email string, subject string, message string) (err error) {
	log.Infof("In production an email would be sent to %s with the following content:\n%s", email, message)
	return
}

//SMTPEmailService implements an email service using plain old SMTP
type SMTPEmailService struct {
	dialer *gomail.Dialer
}

//NewSMTPEmailService creates a nes SMTPEmailService
func NewSMTPEmailService(host string, port int, user string, password string) (service *SMTPEmailService) {
	dialer := gomail.NewDialer(host, port, user, password)
	service = &SMTPEmailService{dialer: dialer}
	return
}

//Send sends an Email
func (s *SMTPEmailService) Send(email string, subject string, message string) (err error) {
	gomsg := gomail.NewMessage()
	gomsg.SetHeader("Subject", subject)
	gomsg.SetHeader("From", "noreply@itsyou.online")
	gomsg.SetHeader("To", email)
	gomsg.SetBody("text/plain", message)
	err = s.dialer.DialAndSend(gomsg)
	if err != nil {
		log.Error("Failed to send email ", err)
	}
	return
}
