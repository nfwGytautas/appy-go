package appy

import (
	"reflect"
	"runtime"
)

func ReflectFunctionName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}
