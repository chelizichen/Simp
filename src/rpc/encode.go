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

func (e *Encode[T]) WriteInt8(tag int, value int8) (bool, error) {
	if int(e.Position) >= len(e.Bytes) {
		// 如果字节切片容量不足，可以进行扩容操作
		newBytes := make([]byte, 2*len(e.Bytes))
		copy(newBytes, e.Bytes)
		e.Bytes = newBytes
	}

	e.Current = int32(tag)
	e.Bytes[e.Position] = byte(value)
	e.Position++

	return false, nil
}

func (e *Encode[T]) WriteString(tag int, value string) (bool, error) {
	// 将字符串转换为字节数组
	stringBytes := []byte(value)

	// 检查字节切片容量是否足够
	if int(e.Position)+len(stringBytes)+4 >= len(e.Bytes) {
		// 如果容量不足，进行扩容操作
		newBytes := make([]byte, 2*(len(e.Bytes)+len(stringBytes)+4))
		copy(newBytes, e.Bytes)
		e.Bytes = newBytes
	}

	// 写入字符串长度
	binary.LittleEndian.PutUint32(e.Bytes[e.Position:], uint32(len(stringBytes)))
	e.Position += 4

	// 写入字符串内容
	copy(e.Bytes[e.Position:], stringBytes)
	e.Position += int32(len(stringBytes))

	e.Current = int32(tag)
	return false, nil
}

func (e *Encode[T]) WriteStruct(tag int, value interface{}) (bool, error) {
	t := reflect.TypeOf(value)
	m, b := t.MethodByName("Encode")
	if !b {
		panic(fmt.Sprintf("Error! Struct %s does not have Method Encode", t.Name()))
	}
	m.Func.Call([]reflect.Value{reflect.ValueOf(value)})
	e.Current = int32(tag)
	return false, nil
}

func (d *Encode[T]) WriteList(tag int, className string, value []interface{}) (bool, error) {
	d.Current = int32(tag)
	return false, nil
}
