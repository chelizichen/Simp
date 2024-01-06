package handlers

import (
	"Simp/config"
	"fmt"
	"net"
)

type SimpTcpClientHandlers struct {
	CONF   *config.SimpConfig
	pid    int
	isOpen bool
	Conn   net.Conn
}

func (s SimpTcpClientHandlers) Write(message string) {
	s.Conn.Write([]byte(message))
}

func NewSimpTcpClient(CONF *config.SimpConfig) (handle SimpTcpClientHandlers, err error) {
	// 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8080")
	handle.CONF = CONF
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	handle.Conn = conn

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Close connecting Error:", err)
		}
	}(conn)

	// 读取服务器的响应
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	response := buffer[:n]
	fmt.Println("Server response:", string(response))
	return handle, nil
}
