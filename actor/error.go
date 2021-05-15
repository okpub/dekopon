package actor

import "errors"

var (
	NilErr = errors.New("the pid is nil")
)

//推送超时错误(为临时错误)
type PublishError struct{}

func (*PublishError) Error() string   { return "send fail with timeout" }
func (*PublishError) Timeout() bool   { return true }
func (*PublishError) Temporary() bool { return true }
