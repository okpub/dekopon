package utils

import "fmt"

var (
	EOF        = fmt.Errorf("ERROR: buffer closed")
	NilErr     = fmt.Errorf("ERROR: null value")
	TimeoutErr = &tempErr{}
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

//临时错误
type tempErr struct{}

func (*tempErr) Error() string   { return "SendTempErr" }
func (*tempErr) String() string  { return "SendTempErr" }
func (*tempErr) Timeout() bool   { return true }
func (*tempErr) Temporary() bool { return true }
