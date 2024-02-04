package rpc

import (
	"fmt"
	"reflect"
)

type Decode[T any] struct {
	Position  int32
	Current   int32
	ClassName string
	Bytes     []byte
}

func (d *Decode[T]) ReadInt8(tag int) int8 {
	d.Current = int32(tag)
	d.Position += 1
	return 1
}

func (d *Decode[T]) ReadString(tag int) string {
	var value string
	d.Current = int32(tag)
	d.Position += int32(len(value))
	return ""
}

func (d *Decode[T]) ReadList(tag int, className string) *[]interface{} {
	d.Current = int32(tag)
	ret := new([]interface{})
	return ret
}

func (d *Decode[T]) ReadStruct(tag int, className string) interface{} {
	clazz := Get(className)
	instance := reflect.New(clazz).Interface()
	name, b := clazz.MethodByName("Decode")
	valLen := d.Bytes[d.Position : d.Position+4]
	val := d.Bytes[d.Position+4 : len(valLen)]
	d.Position += int32(4 + len(valLen))
	if !b {
		panic(fmt.Sprintf("Error! Struct %s does not have Method Decode", className))
	}
	call := name.Func.Call([]reflect.Value{reflect.ValueOf(instance), reflect.ValueOf(val)})
	// 使用类型断言将结果转换为泛型类型 T
	result := call[0]
	return result
}
