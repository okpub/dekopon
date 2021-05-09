package observer

import (
	"fmt"
	"reflect"
)

//对象观察者(用于调度)
type Observer struct {
	caller    reflect.Value
	scheduler reflect.Type
	className string
}

func NewObserver(caller interface{}) *Observer {
	return &Observer{caller: reflect.ValueOf(caller), scheduler: reflect.TypeOf(caller), className: getClassName(caller)}
}

func (b *Observer) Call(name string, args ...interface{}) (values []reflect.Value, err error) {
	var req []reflect.Value
	req = append(req, b.caller)
	for _, v := range args {
		req = append(req, reflect.ValueOf(v))
	}

	if method, ok := b.scheduler.MethodByName(name); ok {
		values = method.Func.Call(req)
	} else {
		err = fmt.Errorf("can't find class=%s method=%s", b.className, name)
	}
	return
}

func (b *Observer) Name() string {
	return b.className
}
