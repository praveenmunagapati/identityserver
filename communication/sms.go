package communication

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
)

//SMSService defines an sms communication channel
type SMSService interface {
	Send(phonenumber string, message string) (err error)
}

//TwilioSMSService is an SMS communication channel using Twilio
type TwilioSMSService struct {
	AccountSID          string
	AuthToken           string
	MessagingServiceSID string
}

//Send sends an SMS
func (s *TwilioSMSService) Send(phonenumber string, message string) (err error) {
	client := &http.Client{}

	data := url.Values{
		"MessagingServiceSid": {s.MessagingServiceSID},
		"To":   {phonenumber},
		"Body": {message},
	}

	req, err := http.NewRequest("POST", "https://api.twilio.com/2010-04-01/Accounts/"+s.AccountSID+"/Messages.json", strings.NewReader(data.Encode()))
	if err != nil {
		log.Error("Error creating sms request: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.AccountSID, s.AuthToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error sending sms via Twilio: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		log.Error("Problem when sending sms via Twilio: ", resp.StatusCode, "\n", string(body))
		err = errors.New("Error sending sms")
	}
	log.Infof("SMS: sms send to %s", phonenumber)
	return
}

// CmTelecomSMSService is an SMS communication channel using cmtelecom
type CmTelecomSMSService struct {
	ProductToken string
}

func (s *CmTelecomSMSService) Send(phonenumber string, message string) (err error) {

	client := &http.Client{}

	type MsgBody struct {
		Content string `json:"content"`
	}

	type Number struct {
		Number string `json:"number"`
	}

	type Msg struct {
		From string   `json:"from"`
		To   []Number `json:"to"`
		Body MsgBody  `json:"body"`
	}

	type Authentication struct {
		ProductToken string `json:"producttoken"`
	}

	type Msgs struct {
		Authentication Authentication `json:"authentication"`
		Message        []Msg          `json:"msg"`
	}

	type SMSData struct {
		Messages Msgs `json:"messages"`
	}

	data := SMSData{
		Messages: Msgs{
			Authentication: Authentication{
				ProductToken: s.ProductToken,
			},
			Message: []Msg{{
				From: "ItsYou.online",
				To: []Number{{
					Number: phonenumber,
				},
				},
				Body: MsgBody{
					Content: message,
				},
			},
			},
		},
	}

	bodyBuf := bytes.NewBuffer(nil)
	err = json.NewEncoder(bodyBuf).Encode(data)
	if err != nil {
		log.Error("Failed to serialize request body: ", err)
		return
	}

	log.Debug(bodyBuf)

	req, err := http.NewRequest("POST", "https://gw.cmtelecom.com/v1.0/message", bodyBuf)
	if err != nil {
		log.Error("Error creating sms request: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error sending sms via CmTelecom: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Error("Problem when sending sms via CmTelecom: ", resp.StatusCode, "\n", string(body))
		err = errors.New("Error sending sms")
	} else {
		log.Debug("CmTelecom response: ", resp.StatusCode, "\n", string(body))
	}
	log.Infof("SMS: sms send to %s", phonenumber)
	return

}
