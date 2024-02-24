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
	result := int16(binary.LittleEndian.Uint16(d.Bytes[d.Position-2 : d.Position]))
	return result
}

func (d *Decode[T]) ReadInt32(tag int) int32 {
	d.Current = int32(tag)
	d.Position += 4
	result := int32(binary.LittleEndian.Uint32(d.Bytes[d.Position-4 : d.Position]))
	return result
}

func (d *Decode[T]) ReadInt64(tag int) int64 {
	d.Current = int32(tag)
	d.Position += 8
	result := int64(binary.LittleEndian.Uint64(d.Bytes[d.Position-8 : d.Position]))
	return result
}

func (d *Decode[T]) ReadString(tag int) string {
	d.Current = int32(tag)
	d.Position += 4
	valueLength := queryStructLen(d.Bytes[d.Position-4 : d.Position])
	value := string(d.Bytes[d.Position : d.Position+valueLength])
	d.Position += valueLength
	return value
}

func (d *Decode[T]) ReadList(tag int, className string) *[]interface{} {
	d.Current = int32(tag)
	ret := new([]interface{})
	return ret
}

func (d *Decode[T]) ReadStruct(tag int, resp interface{}) {
	d.Current = int32(tag)
	t := reflect.TypeOf(resp)
	if t.Kind() != reflect.Ptr {
		panic("Error! WriteStruct expects a pointer to a struct")
	}
	d.Position += 4
	bytes := d.Bytes[d.Position-4 : d.Position]
	valLen := queryStructLen(bytes)
	m, b := t.MethodByName("Decode")
	if !b {
		panic(fmt.Sprintf("Error! Struct %s does not have Method Decode", resp))
	}
	BytesVal := d.Bytes[d.Position : d.Position+valLen]
	callResp := m.Func.Call([]reflect.Value{reflect.ValueOf(resp), reflect.ValueOf(BytesVal)})
	d.Position += valLen
	// 使用类型断言将结果转换为泛型类型 T
	resp = callResp[0]
}
