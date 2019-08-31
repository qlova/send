package send

import "errors"

//Message is a string that can be sent.
type Message struct {
	header string
	string

	attachments []Attachment
}

//SetHeader sets the message header.
func (m *Message) SetHeader(header string) {
	m.header = header
}

//Set sets the message body.
func (m *Message) Set(body string) {
	m.string = body
}

func (m Message) String() string {
	return m.string
}

//Header returns the message's header.
func (m Message) Header() string {
	return m.header
}

//New returns a new message.
func New(body string) Message {
	return Message{"", body, nil}
}

//SendTo tries to send the message and returns nil if the message was delivered.
func (m Message) SendTo(to string) (err error) {
	err = errors.New("could not send message: ")
	for _, sender := range Senders {
		if sender.Type().Detect(to) {
			e := sender.Send(m, to)
			if e == nil {
				return nil
			}
			err = errors.New(err.Error() + e.Error())
		}
	}
	return
}
