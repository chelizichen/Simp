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
	if tag != -1 {
		e.Current = int32(tag)
	}
	e.Bytes[e.Position] = byte(value)
	e.Position += 1
	return true, nil
}

func (e *Encode[T]) WriteInt16(tag int, value int16) (bool, error) {
	e.Dilatation(2)
	if tag != -1 {
		e.Current = int32(tag)
	}
	binary.LittleEndian.PutUint16(e.Bytes[e.Position:], uint16(value))
	e.Position += 2
	return true, nil
}

func (e *Encode[T]) WriteInt32(tag int, value int32) (bool, error) {
	e.Dilatation(2)
	if tag != -1 {
		e.Current = int32(tag)
	}
	binary.LittleEndian.PutUint32(e.Bytes[e.Position:], uint32(value))
	e.Position += 4
	return true, nil
}

func (e *Encode[T]) WriteInt64(tag int, value int64) (bool, error) {
	e.Dilatation(2)
	if tag != -1 {
		e.Current = int32(tag)
	}
	binary.LittleEndian.PutUint64(e.Bytes[e.Position:], uint64(value))
	e.Position += 8
	return true, nil
}

func (e *Encode[T]) WriteString(tag int, value string) (bool, error) {
	// 将字符串转换为字节数组
	stringBytes := []byte(value)
	e.Dilatation(len(stringBytes) + 4)
	if tag != -1 {
		e.Current = int32(tag)
	}
	// 写入字符串长度
	binary.LittleEndian.PutUint32(e.Bytes[e.Position:], uint32(len(stringBytes)))
	e.Position += 4

	// 写入字符串内容
	copy(e.Bytes[e.Position:e.Position+int32(len(stringBytes))], stringBytes)
	e.Position += int32(len(stringBytes))

	return true, nil
}

func (e *Encode[T]) WriteStruct(tag int, value reflect.Value) (bool, error) {
	fmt.Println("value", value.Kind())
	if tag != -1 {
		e.Current = int32(tag)
	}
	var encodeFunc reflect.Value
	if value.Kind() == reflect.Ptr {
		encodeFunc = value.MethodByName("Encode")
	} else {
		encodeFunc = value.Addr().MethodByName("Encode")
	}
	callResp := encodeFunc.Call([]reflect.Value{})
	v := callResp[0].Elem()
	position := v.FieldByName("Position").Interface().(int32)
	bytes := v.FieldByName("Bytes").Interface().([]byte)
	binary.LittleEndian.PutUint32(e.Bytes[e.Position:], uint32(position))
	e.Position += 4
	copy(e.Bytes[e.Position:], bytes)
	e.Position += position
	return true, nil
}

func (e *Encode[T]) WriteList(tag int, value interface{}) (bool, error) {
	e.Current = int32(tag)
	beforeEncodePosition := e.Position
	e.Position += 4
	// 使用类型开关检查 value 的类型并进行相应处理
	switch v := value.(type) {
	case []int8:
		for _, item := range v {
			e.WriteInt8(-1, item)
		}
	case []int16:
		for _, item := range v {
			e.WriteInt16(-1, item)
		}
	case []int32:
		for _, item := range v {
			e.WriteInt32(-1, item)
		}
	case []int:
		for _, item := range v {
			e.WriteInt32(-1, int32(item))
		}
	case []int64:
		for _, item := range v {
			e.WriteInt64(-1, item)
		}

	case []string:
		for _, item := range v {
			e.WriteString(-1, item)
		}
	default:
		length := reflect.ValueOf(v).Len()
		val := reflect.ValueOf(v)
		if length != 0 {
			for i := 0; i < length; i++ {
				v := val.Index(i)
				fmt.Println("v", v)
				e.WriteStruct(-1, v)
			}
		}
	}
	listBytes := e.Position - beforeEncodePosition
	binary.LittleEndian.PutUint32(e.Bytes[beforeEncodePosition:], uint32(listBytes))
	return false, nil
}

func (e *Encode[T]) WriteAny(value interface{}) {
	t := reflect.TypeOf(value)
	fmt.Println(t.Kind())
	switch t.Kind() {
	case reflect.String:
		e.WriteString(-1, value.(string))
	case reflect.Int16:
		e.WriteInt16(-1, value.(int16))
	case reflect.Int32:
		e.WriteInt32(-1, value.(int32))
	case reflect.Int:
		e.WriteInt32(-1, value.(int32))
	case reflect.Int64:
		e.WriteInt64(-1, value.(int64))
	case reflect.Slice:
		e.WriteList(-1, value)
	case reflect.Struct:
		e.WriteStruct(-1, reflect.ValueOf(value))
	default:
	}
}

func (e *Encode[T]) WriteMap(tag int, value reflect.Value) {
	if tag != -1 {
		e.Current = int32(tag)
	}
	beforeEncodePosition := e.Position
	e.Position += 4
	fmt.Println("value | ", value.MapRange())
	iter := value.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		e.WriteString(-1, k.String())
		e.WriteAny(v.Interface())
	}
	currentPosition := e.Position - beforeEncodePosition
	binary.LittleEndian.PutUint32(e.Bytes[beforeEncodePosition:], uint32(currentPosition))
}
