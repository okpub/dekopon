package observer

import (
	"reflect"
	"strings"

	"github.com/okpub/dekopon/utils"
)

func getClassName(caller interface{}) string {
	atype := reflect.TypeOf(caller)
	switch atype.Kind() {
	case reflect.Ptr:
		return atype.Elem().Name()
	case reflect.Struct:
		return atype.Name()
	default:
		return ""
	}
}

func PackValues(args []reflect.Value) (arr []interface{}) {
	for _, v := range args {
		arr = append(arr, v.Interface())
	}
	return
}

func GetValue(args []reflect.Value, die error) (res interface{}, err error) {
	if err = die; err == nil {
		if len(args) > 0 {
			res = args[0].Interface()
		} else {
			err = utils.NilErr
		}
	}
	return
}

func ToRouter(route string) (className string, methodName string) {
	var arr = strings.Split(route, ".")
	if len(arr) > 1 {
		return arr[0], arr[1]
	}
	return arr[0], "Index"
}
