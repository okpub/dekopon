package actor

type SystemMessage interface {
	SystemPermission()
}

var (
	EVENT_START      = &Started{}
	EVENT_STOP       = &Stopped{}
	EVENT_RESTARTING = &Restarting{}
)

//system event
type (
	Started    struct{}
	Stopped    struct{}
	Restarting struct{}
)

func (Started) String() string       { return "start" }
func (Started) SystemPermission()    {}
func (Stopped) String() string       { return "stop" }
func (Stopped) SystemPermission()    {}
func (Restarting) String() string    { return "restarting" }
func (Restarting) SystemPermission() {}

//func request
type Function func()

func (fn Function) Done()          { fn() }
func (fn Function) String() string { return "function" }

//watch
type (
	Terminated struct{ Who PID }
	Watch      struct{ Who PID }
	Unwatch    struct{ Who PID }
)

func (this *Terminated) GetWho() PID  { return this.Who }
func (*Terminated) SystemPermission() {}
func (*Terminated) String() string    { return "terminated" }

func (this *Watch) GetWho() PID  { return this.Who }
func (*Watch) SystemPermission() {}
func (*Watch) String() string    { return "watch" }

func (this *Unwatch) GetWho() PID  { return this.Who }
func (*Unwatch) SystemPermission() {}
func (*Unwatch) String() string    { return "unwatch" }

//fail
type RanError struct {
	Who    PID
	Err    error
	Result interface{}
}

func (this *RanError) GetWho() PID         { return this.Who }
func (this *RanError) GetWhy() interface{} { return this.Result }
func (this *RanError) Error() error        { return this.Err }
func (this *RanError) String() string      { return "failure" }
func (this *RanError) SystemPermission()   {}
