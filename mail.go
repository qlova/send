package mail

import (
	"errors"
	"fmt"
	"net"
	"strings"
)
import "github.com/nilslice/email"
import "github.com/mileusna/spf"
import "github.com/glendc/go-external-ip"

func New() Mail {
	return Mail{}
}

type Mail struct {
	To, From, Subject, Body string
}

func (m *Mail) SetTo(to string) {
	m.To = to
}
func (m *Mail) SetSubject(subject string) {
	m.Subject = subject
}
func (m *Mail) SetBody(body string) {
	m.Body = body
}
func (m *Mail) SetFrom(from string) {
	m.From = from
}

func (m Mail) Send() error {
	return DefaultSender.Send(m)
}

type Sender interface {
	Send(Mail) error
}

var PublicAddress net.IP
var DomainCache = make(map[string]bool)

func WeAreNotAllowedToSend(mail Mail) bool {

	if PublicAddress == nil {
		var err error
		PublicAddress, err = externalip.DefaultConsensus(nil, nil).ExternalIP()
		if err != nil {
			fmt.Println(err)
			return true
		}
	}

	components := strings.Split(mail.From, "@")
	_, domain := components[0], components[1]

	if cached, ok := DomainCache[domain]; ok {
		return !cached
	}

	var result = spf.CheckHost(PublicAddress, domain, mail.From, "")

	DomainCache[domain] = result == spf.Pass
	return !(result == spf.Pass)
}

var DefaultSender Sender = RawSender{}

type RawSender struct{}

func (raw RawSender) Send(m Mail) error {
	if WeAreNotAllowedToSend(m) {
		return errors.New("Your IP Address (" + PublicAddress.String() + ") is not configured to send emails from " + m.From)
	}

	var e = email.Message{
		To:      m.To,
		From:    m.From,
		Subject: m.Subject,
		Body:    m.Body,
	}

	return e.Send()
}
