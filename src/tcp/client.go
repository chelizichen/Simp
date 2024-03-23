package tcp

import (
	"Simp/src/rpc"
	"net"
	"os"
	"reflect"
)

type TarsusStruct[T any] struct {
	TraceIds []string
	Module   string
	Method   string
	Request  string
	Body     *T
}

func (r *TarsusStruct[T]) Decode(Bytes []byte) *TarsusStruct[T] {
	d := new(rpc.Decode[TarsusStruct[T]])
	d.ClassName = "TarsusStruct"
	d.Bytes = Bytes
	r.Module = d.ReadString(1)
	r.Method = d.ReadString(2)
	r.Request = d.ReadString(3)
	r.Body = d.ReadStruct(4, reflect.ValueOf(new(T))).Interface().(*T)
	return r
}

func (r *TarsusStruct[T]) Encode() *rpc.Encode[TarsusStruct[T]] {
	d := new(rpc.Encode[TarsusStruct[T]])
	d.ClassName = "TarsusStruct"
	d.Bytes = make([]byte, 1024)
	d.WriteString(1, r.Module)
	d.WriteString(2, r.Method)
	d.WriteString(3, r.Request)
	d.WriteStruct(4, reflect.ValueOf(r.Body))
	return d
}

var Servants = make(map[string]net.Conn)
var InvokeCallback = make(map[int]func())

func SimpTcpWrite(serverName string, buffer []byte, callback func()) {
	conn := Servants[serverName]
	conn.Write(buffer)
	eid := os.Getuid()
	InvokeCallback[eid] = callback
}
