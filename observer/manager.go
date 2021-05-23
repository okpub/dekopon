package observer

import (
	"fmt"
	"reflect"
)

type ObserverSet map[string]*Observer

type Manager struct {
	observers ObserverSet
}

func NewManager() *Manager {
	return &Manager{observers: make(ObserverSet)}
}

func (manager *Manager) On(caller interface{}) (err error) {
	var observer = NewObserver(caller)
	manager.observers[observer.Name()] = observer
	return
}

func (manager *Manager) Off(className string) (err error) {
	delete(manager.observers, className)
	return
}

func (manager *Manager) Router(className, method string, args ...interface{}) (values []reflect.Value, err error) {
	if observer, ok := manager.observers[className]; ok {
		values, err = observer.Call(method, args...)
	} else {
		err = fmt.Errorf("can't find class=%s", className)
	}
	return
}

func (manager *Manager) Emit(route string, args ...interface{}) ([]reflect.Value, error) {
	var className, methodName = ToRouter(route)
	return manager.Router(className, methodName, args...)
}

//全局默认管理
var (
	defaultManager = NewManager()
)

func Register(caller interface{}) {
	defaultManager.On(caller)
}

func Unregister(className string) {
	defaultManager.Off(className)
}

func Router(className, method string, args ...interface{}) (values []reflect.Value, err error) {
	return defaultManager.Router(className, method, args...)
}

func Emit(route string, args ...interface{}) (values []reflect.Value, err error) {
	var className, methodName = ToRouter(route)
	return Router(className, methodName, args...)
}
