package send

//Sender can send messages.
type Sender interface {
	Type() Type
	Send(m Message, to string) error
}

//Senders are the implicitly available senders, they can be registered with send.Register()
//When a message is sent
var Senders []Sender

//Register registers a new sender.
func Register(s Sender) {
	Senders = append(Senders, s)
}
