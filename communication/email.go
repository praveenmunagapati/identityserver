package communication

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-gomail/gomail"
)

type EmailService interface {
	Send(email string, subject string, message string) (err error)
}

type DevEmailService struct{}

//Send sends an Email
func (s *DevEmailService) Send(email string, subject string, message string) (err error) {
	log.Infof("In production an email would be sent to %s with the following content:\n%s", email, message)
	return
}


type SMTPEmailService struct{
	dialer *gomail.Dialer
}

func NewSMTPEmailService (host string, port int, user string, password string) (service *SMTPEmailService) {
	dialer := gomail.NewDialer(host, port, user, password)
	service = &SMTPEmailService{dialer: dialer}
	return
}

//Send sends an Email
func (s *SMTPEmailService) Send(email string, subject string, message string) (err error) {
	gomsg := gomail.NewMessage()
	gomsg.SetHeader("Subject", subject)
	gomsg.SetHeader("From", "support@itsyou.online")
	gomsg.SetHeader("To", email)
	gomsg.SetBody("text/plain", message)
	err = s.dialer.DialAndSend(gomsg)
	if err != nil {
		log.Error("Failed to send email ", err)
	}
	return
}

