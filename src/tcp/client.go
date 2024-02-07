package tcp

import (
	h "Simp/src/http"
	"Simp/src/rpc"
	"fmt"
	"net"
	"os"
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
	_body := new(T)
	d.ReadStruct(4, _body)
	r.Body = _body
	return r
}

func (r *TarsusStruct[T]) Encode() *rpc.Encode[TarsusStruct[T]] {
	d := new(rpc.Encode[TarsusStruct[T]])
	d.ClassName = "TarsusStruct"
	d.Bytes = make([]byte, 1024)
	d.WriteString(1, r.Module)
	d.WriteString(2, r.Method)
	d.WriteString(3, r.Request)
	d.WriteStruct(4, r.Body)
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

// 代理服务端链接
func ClientListen(ctx *h.SimpHttpServerCtx) {
	for _, v := range *ctx.Proxy {
		// 连接到服务器
		conn, err := net.Dial("tcp", ":"+v.Server.Port)
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			os.Exit(1)
		}
		defer conn.Close()
		Servants[v.Server.Name] = conn
	}
	for serverName, conn := range Servants {
		// 读取服务器的响应内容
		// 启动一个 goroutine 读取数据并发送到通道
		CONN := conn
		SERVER_NAME := serverName
		go func() {
			buffer := make([]byte, 1024)
			_, err := CONN.Read(buffer)
			handleBuffer(SERVER_NAME, &buffer)
			if err != nil {
				fmt.Println("Error reading response:", err)
				return
			}
		}()
	}
}

func handleBuffer(serverName string, buf *[]byte) (err error) {
	ts := new(TarsusStruct[any])
	err = readHead[any](ts, buf)
	err = readRespBody[any](ts, buf)
	return err
}

func readHead[T any](tarsus *TarsusStruct[T], buf *[]byte) (err error) {
	return
}

func readRespBody[T any](tarsus *TarsusStruct[T], buf *[]byte) (err error) {
	return
}
