package rpc

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

type Encode[T any] struct {
	Position  int32
	Current   int32
	ClassName string
	Bytes     []byte
}

func (e *Encode[T]) Dilatation(size int) {
	if int(e.Position)+size >= len(e.Bytes) {
		// 如果字节切片容量不足，可以进行扩容操作
		newBytes := make([]byte, 2*len(e.Bytes)+size)
		copy(newBytes, e.Bytes)
		e.Bytes = newBytes
	}
}

func (e *Encode[T]) WriteInt8(tag int, value int8) (bool, error) {
	e.Dilatation(1)
	e.Current = int32(tag)
	e.Bytes[e.Position] = byte(value)
	e.Position += 1
	return true, nil
}

func (e *Encode[T]) WriteInt16(tag int, value int16) (bool, error) {
	e.Dilatation(2)
	e.Current = int32(tag)
	binary.LittleEndian.PutUint16(e.Bytes[e.Position:], uint16(value))
	e.Position += 2
	return true, nil
}

func (e *Encode[T]) WriteString(tag int, value string) (bool, error) {
	// 将字符串转换为字节数组
	stringBytes := []byte(value)
	e.Dilatation(len(stringBytes) + 4)
	// 写入字符串长度
	binary.LittleEndian.PutUint32(e.Bytes[e.Position:], uint32(len(stringBytes)))
	e.Position += 4

	// 写入字符串内容
	copy(e.Bytes[e.Position:], stringBytes)
	e.Position += int32(len(stringBytes))

	e.Current = int32(tag)
	return true, nil
}

func (e *Encode[T]) WriteStruct(tag int, value interface{}) (bool, error) {
	t := reflect.TypeOf(value) // Use Elem() to get the element type (remove the pointer)
	if t.Kind() != reflect.Ptr {
		panic("Error! WriteStruct expects a pointer to a struct")
	}
	m, b := t.MethodByName("Encode")
	if !b {
		panic(fmt.Sprintf("Error! Struct %s does not have Method Encode", t.Name()))
	}
	callResp := m.Func.Call([]reflect.Value{reflect.ValueOf(value)})
	v := callResp[0].Elem()
	position := v.FieldByName("Position").Interface().(int32)
	bytes := v.FieldByName("Bytes").Interface().([]byte)
	e.Current = int32(tag)
	binary.LittleEndian.PutUint32(e.Bytes[e.Position:], uint32(position))
	e.Position += 4
	copy(e.Bytes[e.Position:], bytes)
	e.Position += int32(position)
	return true, nil
}

func (d *Encode[T]) WriteList(tag int, className string, value []interface{}) (bool, error) {
	d.Current = int32(tag)
	return false, nil
}
