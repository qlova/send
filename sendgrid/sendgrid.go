package sendgrid

import "errors"
import "encoding/json"
import "strings"

import "github.com/qlova/mail"
import "github.com/sendgrid/sendgrid-go"
import helper "github.com/sendgrid/sendgrid-go/helpers/mail"

type Sender struct {
	*sendgrid.Client
}

func NewSender(key string) Sender {
	return Sender{sendgrid.NewSendClient(key)}
}

type Error struct {
	Message string `json:"message"`
}

type Response struct {
	Errors []Error `json:"errors"`
}

func (sender Sender) Send(m mail.Mail) error {
	email := helper.NewV3Mail()

	var p = helper.NewPersonalization()
	p.AddTos(helper.NewEmail("", m.To))

	email.AddPersonalizations(p)
	email.SetFrom(helper.NewEmail("", m.From))
	email.Subject = m.Subject
	email.AddContent(&helper.Content{
		Type:  "text/plain",
		Value: m.Body,
	})

	response, err := sender.Client.Send(email)
	if err != nil {
		return err
	}

	var r Response
	json.NewDecoder(strings.NewReader(response.Body)).Decode(&r)

	if len(r.Errors) > 0 {
		return errors.New(r.Errors[0].Message)
	}

	return err
}
