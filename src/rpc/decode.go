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
	if tag != -1 {
		d.Current = int32(tag)
	}
	d.Current = int32(tag)
	d.Position += 1
	bytes := d.Bytes[d.Position-1 : d.Position]
	i := int8(bytes[0])
	return i
}

func (d *Decode[T]) ReadInt16(tag int) int16 {
	if tag != -1 {
		d.Current = int32(tag)
	}
	d.Current = int32(tag)
	d.Position += 2
	result := int16(binary.LittleEndian.Uint16(d.Bytes[d.Position-2 : d.Position]))
	return result
}

func (d *Decode[T]) ReadInt32(tag int) int32 {
	if tag != -1 {
		d.Current = int32(tag)
	}
	d.Current = int32(tag)
	d.Position += 4
	result := int32(binary.LittleEndian.Uint32(d.Bytes[d.Position-4 : d.Position]))
	fmt.Println("Result", result)
	return result
}

func (d *Decode[T]) ReadInt64(tag int) int64 {
	if tag != -1 {
		d.Current = int32(tag)
	}
	d.Current = int32(tag)
	d.Position += 8
	result := int64(binary.LittleEndian.Uint64(d.Bytes[d.Position-8 : d.Position]))
	return result
}

func (d *Decode[T]) ReadString(tag int) string {
	if tag != -1 {
		d.Current = int32(tag)
	}
	d.Position += 4
	valueLength := queryStructLen(d.Bytes[d.Position-4 : d.Position])
	value := string(d.Bytes[d.Position : d.Position+valueLength])
	d.Position += valueLength
	return value
}

func (d *Decode[T]) ReadStruct(tag int, resp reflect.Value) reflect.Value {
	if tag != -1 {
		d.Current = int32(tag)
	}
	d.Position += 4
	bytes := d.Bytes[d.Position-4 : d.Position]
	valLen := queryStructLen(bytes)
	m := resp.MethodByName("Decode")
	BytesVal := d.Bytes[d.Position : d.Position+valLen]
	callResp := m.Call([]reflect.Value{reflect.ValueOf(BytesVal)})
	d.Position += valLen
	// 使用类型断言将结果转换为泛型类型 T
	resp = callResp[0]
	return resp
}

func (d *Decode[T]) ReadList(tag int, value interface{}) interface{} {
	if tag != -1 {
		d.Current = int32(tag)
	}
	d.Position += 4
	valueLength := queryStructLen(d.Bytes[d.Position-4 : d.Position])
	currPosition := d.Position
	switch v := value.(type) {
	case []int8:
		length := int(valueLength/1) - 1
		for i := 0; i < length; i++ {
			r := d.ReadInt8(-1)
			v = append(v, r)
		}
		return v
	case []int16:
		length := int(valueLength/2) - 1
		for i := 0; i < length; i++ {
			r := d.ReadInt16(-1)
			v = append(v, r)
		}
		return v
	case []int32:
		length := int(valueLength/4) - 1
		for i := 0; i < length; i++ {
			r := d.ReadInt32(-1)
			v = append(v, r)
		}
		return v
	case []int64:
		length := int(valueLength/8) - 1
		for i := 0; i < length; i++ {
			r := d.ReadInt64(-1)
			v = append(v, r)
		}
		return v
	case []string:
		for {
			if d.Position == currPosition+valueLength {
				break
			}
			s := d.ReadString(-1)
			v = append(v, s)
		}
		return v
	default:
		target := reflect.ValueOf(v).Type()
		resp := reflect.MakeSlice(target, 0, 0)
		for {
			if d.Position == currPosition+valueLength-4 {
				break
			}
			t := reflect.New(reflect.ValueOf(v).Type().Elem())
			s := d.ReadStruct(-1, t)
			resp = reflect.Append(resp, s.Elem())
		}
		return resp.Interface()
	}
}

func (d *Decode[T]) ReadAny() reflect.Value {

	return reflect.ValueOf(nil)
}

func (d *Decode[T]) ReadMap(tag int) map[string]interface{} {
	if tag != -1 {
		d.Current = int32(tag)
	}
	d.Position += 4
	bytes := d.Bytes[d.Position-4 : d.Position]
	valLen := queryStructLen(bytes)
	BytesVal := d.Bytes[d.Position : d.Position+valLen]
	resp := make(map[string]interface{})
	for true {
		k := d.ReadString(-1)
		v := d.ReadAny()
		resp[k] = v.Interface()
	}
	fmt.Println(BytesVal)
	return nil
}
