package sendgrid

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/qlova/send"
	"github.com/qlova/send/with/email"
	"github.com/sendgrid/sendgrid-go"

	helper "github.com/sendgrid/sendgrid-go/helpers/mail"
)

//Options includes the sendgrid key and the from address for the emails.
type Options struct {
	Key, From string
}

//Sender contains the sendgrid Client and a mutable from address.
type Sender struct {
	*sendgrid.Client
	From string
}

//New returns a new sendgrid Sender.
func New(o Options) Sender {
	return Sender{sendgrid.NewSendClient(o.Key), o.From}
}

//Error is an internal struct for detecting errors.
type Error struct {
	Message string `json:"message"`
}

//Response is an internal struct for detecting errors.
type Response struct {
	Errors []Error `json:"errors"`
}

//Type returns the Email send type.
func (sender Sender) Type() send.Type {
	return email.Email{}
}

//Send attempts to send the message.
func (sender Sender) Send(m send.Message, to string) error {
	email := helper.NewV3Mail()

	var p = helper.NewPersonalization()
	p.AddTos(helper.NewEmail("", to))

	email.AddPersonalizations(p)
	email.SetFrom(helper.NewEmail("", sender.From))
	email.Subject = m.Header()
	email.AddContent(&helper.Content{
		Type:  "text/plain",
		Value: m.String(),
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
