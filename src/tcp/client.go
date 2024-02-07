package tcp

import (
	h "Simp/src/http"
	"fmt"
	"net"
	"os"
)

var Servants = make(map[string]net.Conn)

func SimpTcpWrite(serverName string, buffer []byte) {
	conn := Servants[serverName]
	conn.Write(buffer)
}

func SimpTcpRead(serverName string) {

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
		go func() {
			buffer := make([]byte, 1024)
			_, err := conn.Read(buffer)
			handleBuffer(serverName, buffer)
			if err != nil {
				fmt.Println("Error reading response:", err)
				return
			}
		}()
	}
}

func handleBuffer(serverName string, buf []byte) {

}
