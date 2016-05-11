package communication

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Sirupsen/logrus"
)

//SMSService is an SMS communication channel using Twilio
type SMSService struct {
	AccountSID          string
	AuthToken           string
	MessagingServiceSID string
}

//Send sends an SMS
func (s *SMSService) Send(phonenumber string, message string) (err error) {
	client := &http.Client{}

	data := url.Values{
		"MessagingServiceSid": {s.MessagingServiceSID},
		"To":   {phonenumber},
		"Body": {message},
	}

	req, err := http.NewRequest("POST", "https://api.twilio.com/2010-04-01/Accounts/"+s.AccountSID+"/Messages.json", strings.NewReader(data.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.AccountSID, s.AuthToken)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusAccepted {
		logrus.Error("Problem when sending sms via Twilio: ", resp.StatusCode, "\n", body)
		err = errors.New("Error sending sms")
	}
	return
}
