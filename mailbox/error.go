package mailbox

var (
	TempErr = &PublishErr{}
)

//临时错误
type PublishErr struct{}

func (*PublishErr) Error() string   { return "PublishErr" }
func (*PublishErr) String() string  { return "PublishErr" }
func (*PublishErr) Timeout() bool   { return true }
func (*PublishErr) Temporary() bool { return true }
