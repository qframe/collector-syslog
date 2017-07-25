package qtypes_messages


type Message struct {
	Base
	Message string
}

func NewMessage(b Base, msg string) Message {
	return Message{
		Base: b,
		Message: msg,
	}
}
