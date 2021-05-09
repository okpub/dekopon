package actor

type ActorMessage struct {
	sender  PID
	message interface{}
}

func MSG(message interface{}) MessageEnvelope {
	return REQ(message, nil)
}

func REQ(message interface{}, pid PID) MessageEnvelope {
	return &ActorMessage{message: message, sender: pid}
}

func (this *ActorMessage) Sender() PID          { return this.sender }
func (this *ActorMessage) Message() interface{} { return this.message }

//static func
func WrapEnvelope(any interface{}) MessageEnvelope {
	if message, ok := any.(MessageEnvelope); ok {
		return message
	}
	return MSG(any)
}

func GetSender(any interface{}) (sender PID) {
	if message, ok := any.(MessageEnvelope); ok {
		sender = message.Sender()
	}
	return
}

func GetMessage(any interface{}) interface{} {
	if message, ok := any.(MessageEnvelope); ok {
		return message.Message()
	}
	return any
}
