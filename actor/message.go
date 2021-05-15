package actor

//class
type ActorMessage struct {
	sender  PID
	message interface{}
}

func (env *ActorMessage) Sender() PID          { return env.sender }
func (env *ActorMessage) Message() interface{} { return env.message }

//static message func
func MSG(message interface{}) MessageEnvelope {
	return WrapMessage(message, nil)
}

func WrapMessage(message interface{}, pid PID) MessageEnvelope {
	return &ActorMessage{message: message, sender: pid}
}

//static envelope func
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
