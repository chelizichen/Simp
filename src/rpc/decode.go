package rpc

import (
	"encoding/binary"
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
	bytes := d.Bytes[d.Position-1 : d.Position]
	i := int8(bytes[0])
	return i
}

func (d *Decode[T]) ReadInt16(tag int) int16 {
	d.Current = int32(tag)
	d.Position += 2

	// 检查字节切片长度是否足够
	if int(d.Position)+1 > len(d.Bytes) {
		// 处理越界情况，这里可能需要根据你的具体需求进行处理
		return 0
	}

	// 将字节切片转换为int16
	result := int16(binary.LittleEndian.Uint16(d.Bytes[d.Position-2 : d.Position]))
	return result
}

func (d *Decode[T]) ReadString(tag int) string {
	d.Current = int32(tag)
	d.Position += 4
	valueLength := int32(binary.LittleEndian.Uint32(d.Bytes[d.Position-4 : d.Position]))
	value := string(d.Bytes[d.Position : d.Position+valueLength])
	return value
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
