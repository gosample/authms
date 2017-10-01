package smtp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/tomogoma/authms/model"
	"github.com/tomogoma/go-commons/errors"
)

type ConfigStore interface {
	IsNotFoundError(error) bool
	GetSMTPConfig() ([]byte, error)
}

type Mailer struct {
	errors.NotFoundErrCheck
	db    ConfigStore
	tmplt *template.Template
}

func New(cs ConfigStore) (*Mailer, error) {
	if cs == nil {
		return nil, errors.New("ConfigStore was nil")
	}
	tmpl, err := template.New("email").Parse(emailTmplt)
	if err != nil {
		return nil, errors.Newf("parse email template: %v", err)
	}
	return &Mailer{db: cs, tmplt: tmpl}, nil
}

func (m *Mailer) SendEmail(email model.SendMail) error {
	conf, err := m.getConfig()
	if err != nil {
		return err
	}
	msg, err := m.generateMessage(email)
	if err != nil {
		return errors.Newf("generate email message: %v", err)
	}
	return sendMessage(*conf, email.ToEmails, msg)
}

func (m *Mailer) getConfig() (*Config, error) {
	confB, err := m.db.GetSMTPConfig()
	if err != nil {
		if m.db.IsNotFoundError(err) {
			return nil, errors.NewNotFound("SMTP not configured")
		}
		return nil, errors.Newf("get SMTP configuration: %v", err)
	}
	conf := new(Config)
	if err := json.Unmarshal(confB, conf); err != nil {
		return nil, errors.Newf("unmarshal SMTP configuration: %v", err)
	}
	return conf, nil
}

func (m *Mailer) generateMessage(r model.SendMail) ([]byte, error) {
	email := newEmail(r)
	eb := bytes.NewBuffer(make([]byte, 0, 256))
	if err := m.tmplt.Execute(eb, email); err != nil {
		return nil, errors.Newf("execute template: %v", err)
	}
	return eb.Bytes(), nil
}

func sendMessage(conf Config, recipients []string, msg []byte) error {
	auth := smtp.PlainAuth("", conf.Username, conf.Password, conf.ServerAddress)
	addr := fmt.Sprintf("%s:%d", conf.ServerAddress, conf.TLSPort)
	err := smtp.SendMail(addr, auth, conf.FromEmail, recipients, msg)
	if err != nil {
		return errors.Newf("send mail: %v", err)
	}
	return nil
}
