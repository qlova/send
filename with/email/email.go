package email

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/mileusna/spf"
	"github.com/nilslice/email"
	"github.com/qlova/send"

	externalip "github.com/glendc/go-external-ip"
)

//Email is the send type.
type Email struct{}

//Name returns the send type name.
func (email Email) Name() string {
	return "email"
}

//Detect detects if the to address is an email.
func (email Email) Detect(to string) bool {
	return strings.Contains(to, "@")
}

//Sender is a naive email sender that sends emails directly.
type Sender struct {
	From string
}

//New returns a new email sender.
func New(from string) Sender {
	return Sender{
		From: from,
	}
}

//Type returns the Email send type.
func (sender Sender) Type() send.Type {
	return Email{}
}

//Send sends a message as an email.
func (sender Sender) Send(m send.Message, to string) error {
	if WeAreNotAllowedToSend(sender.From) {
		return errors.New("Your IP Address (" + publicAddress.String() + ") is not configured to send emails from " + sender.From)
	}

	var e = email.Message{
		To:      to,
		From:    sender.From,
		Subject: m.Header(),
		Body:    m.String(),
	}

	return e.Send()
}

var publicAddress net.IP
var domainCache = make(map[string]bool)

//WeAreAllowedToSend returns true if emails are permitted to be sent from 'from' on this server.
func WeAreAllowedToSend(from string) bool {

	if publicAddress == nil {
		var err error
		publicAddress, err = externalip.DefaultConsensus(nil, nil).ExternalIP()
		if err != nil {
			fmt.Println(err)
			return true
		}
	}

	components := strings.Split(from, "@")
	_, domain := components[0], components[1]

	if cached, ok := domainCache[domain]; ok {
		return !cached
	}

	var result = spf.CheckHost(publicAddress, domain, from, "")

	domainCache[domain] = result == spf.Pass
	return (result == spf.Pass)
}

//WeAreNotAllowedToSend returns true if emails are permitted to be sent from 'from' on this server.
func WeAreNotAllowedToSend(from string) bool {
	return !WeAreAllowedToSend(from)
}
