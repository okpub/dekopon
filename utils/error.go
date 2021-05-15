package utils

import "fmt"

var (
	EOF     = fmt.Errorf("ERROR: buffer closed")
	NilErr  = fmt.Errorf("ERROR: null value")
	TempErr = &PublishErr{}
)

func Die(err error) bool {
	return err != nil
}

func CatchErr(obj interface{}) (err error) {
	switch any := obj.(type) {
	case error:
		err = any
	case nil:
		//no error
	default:
		err = fmt.Errorf("%v", any)
	}
	return
}

//捕获通道错误，如果错误了就EOF
func CatchDie(obj interface{}) (err error) {
	if obj == nil {
		return
	}
	err = EOF
	return
}

//close
func SafeClose(task chan interface{}) (err error) {
	defer func() { err = CatchDie(recover()) }()
	close(task)
	return
}

func SafeDone(done chan<- struct{}) (err error) {
	defer func() { err = CatchDie(recover()) }()
	close(done)
	return
}

//临时错误
type PublishErr struct{}

func (*PublishErr) Error() string   { return "PublishErr" }
func (*PublishErr) String() string  { return "PublishErr" }
func (*PublishErr) Timeout() bool   { return true }
func (*PublishErr) Temporary() bool { return true }
