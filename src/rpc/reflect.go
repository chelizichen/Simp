package rpc

import "reflect"

var SimpReflect = make(map[string]reflect.Type)

func Get(className string) reflect.Type {
	return SimpReflect[className]
}

func Set(className string, s interface{}) {
	SimpReflect[className] = reflect.TypeOf(s)
}
