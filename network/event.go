package network

//actor系统消息，不允许转发
var (
	EVENT_OPEN   = &EventOpen{}
	EVENT_CLOSED = &EventClose{}
)

//event open
type EventOpen struct{}

func (*EventOpen) String() string   { return "EventOpen" }
func (EventOpen) SystemPermission() {}

//event close
type EventClose struct{}

func (*EventClose) String() string   { return "EventClose" }
func (EventClose) SystemPermission() {}

//event open error
type DialError struct {
	Err error
}

func (DialError) SystemPermission()   {}
func (err *DialError) String() string { return err.Err.Error() }

//临时错误(读超时)
type TempErr struct {
	Err error
}

func (TempErr) SystemPermission()   {}
func (err *TempErr) String() string { return err.Err.Error() }
