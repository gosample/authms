package africas_talking

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	errors "github.com/tomogoma/go-typed-errors"
)

const (
	atSendURL = "https://api.africastalking.com/restless/send"

	paramUserName = "username"
	paramAPIKey   = "Apikey"
	paramToPhone  = "to"
	paramMessage  = "message"
)

type Config interface {
	Username() string
	APIKey() string
}

type Option func(at *SMSCl)

type atResponse struct {
	SMSMessageData struct {
		Recipients struct {
			Recipient []struct {
				Status struct {
					Val string `xml:",chardata"`
				} `xml:"status"`
			} `xml:"Recipient"`
		} `xml:"Recipients"`
		Message string `xml:"Message"`
	} `xml:"SMSMessageData"`
}

type SMSCl struct {
	atSendURL string
	userName  string
	apiKey    string
}

func SendURL(URL string) func(at *SMSCl) {
	return func(at *SMSCl) {
		at.atSendURL = URL
	}
}

func NewSMSCl(usrName, APIKey string, opts ...Option) (*SMSCl, error) {
	if APIKey == "" {
		return nil, errors.New("API key was empty")
	}
	if usrName == "" {
		return nil, errors.New("API UserName was empty")
	}
	at := &SMSCl{atSendURL: atSendURL, userName: usrName, apiKey: APIKey}
	for _, opt := range opts {
		opt(at)
	}
	if at.atSendURL == "" {
		return nil, errors.New("Send URL was empty")
	}
	return at, nil
}

func (at *SMSCl) SMS(toPhone, message string) error {
	resp, err := at.sendRequest(toPhone, message)
	if err != nil {
		return errors.Newf("error sending request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return errors.Newf("error connecting to API: %s: %s", resp.Status, respBody)
	}
	respBody, err := readRespBody(resp.Body)
	if err != nil {
		return err
	}
	recipients := respBody.SMSMessageData.Recipients.Recipient
	if len(recipients) != 1 {
		return errors.Newf("%d recipients recorded, expecting 1 - got err message(%s)",
			len(recipients), respBody.SMSMessageData.Message)
	}
	if !strings.EqualFold(strings.TrimSpace(recipients[0].Status.Val), "success") {
		return errors.Newf("API reported an error: %v", recipients[0].Status.Val)
	}
	return nil
}

func (at *SMSCl) sendRequest(toPhone, message string) (*http.Response, error) {
	if toPhone == "" {
		return nil, errors.Newf("toPhone was empty")
	}
	if message == "" {
		return nil, errors.Newf("message was empty")
	}
	URL, err := url.Parse(at.atSendURL)
	if err != nil {
		return nil, errors.Newf("error parsing configured send URL: %v", err)
	}
	q := URL.Query()
	q.Add(paramUserName, at.userName)
	q.Add(paramAPIKey, at.apiKey)
	q.Add(paramToPhone, toPhone)
	q.Add(paramMessage, message)
	URL.RawQuery = q.Encode()
	return http.Get(URL.String())
}

func readRespBody(resp io.Reader) (atResponse, error) {
	respBody, err := ioutil.ReadAll(resp)
	if err != nil {
		return atResponse{}, errors.Newf("error reading response body: %v", err)
	}
	respStruct := atResponse{}
	if err := xml.Unmarshal(respBody, &respStruct); err != nil {
		return atResponse{}, errors.Newf("error unmarshalling response body: %v", err)
	}
	return respStruct, nil
}
