package cluster

import (
	"reflect"

	"github.com/okpub/dekopon/bean/message/rpc"
)

func ValueOf(values []reflect.Value, fail error) (res *rpc.Response, err error) {
	if err = fail; err == nil {
		if values[0].IsNil() {
			//empty ok
			res = &rpc.Response{ErrMsg: "no handler!"}
		} else {
			res = values[0].Interface().(*rpc.Response)
		}

		if values[1].IsNil() {
			//no value
		} else {
			err = values[1].Interface().(error)
		}
	}
	return
}
