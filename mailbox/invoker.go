package mailbox

type InvokerFunc func(interface{})

func (fn InvokerFunc) InvokeUserMessage(message interface{}) {
	fn(message)
}

func (fn InvokerFunc) InvokeSystemMessage(message interface{}) {
	fn(message)
}

func (fn InvokerFunc) EscalateFailure(err error, message interface{}) {
	fn(err)
}
