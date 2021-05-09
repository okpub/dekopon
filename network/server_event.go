package network

var (
	EVENT_START = &EventStart{}
	EVENT_STOP  = &EventStop{}
)

type EventStart struct{}

func (*EventStart) String() string {
	return "EventStart"
}

type EventStop struct{}

func (*EventStop) String() string {
	return "EventStop"
}
