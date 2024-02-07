package tcp

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func ServerListen() {
	// 监听本地的8080端口
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Listening on :8080...")

	for {
		// 接受新的连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// 处理连接
		go handleRequest(conn)
	}
}

// 处理连接请求
func handleRequest(conn net.Conn) {
	defer conn.Close()

	// 创建一个新的读取器
	reader := bufio.NewReader(conn)

	for {
		// 读取数据
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}
		// 打印接收到的数据
		fmt.Print("Message received: ", string(message))

		// 发送回应
		conn.Write([]byte("Message received.\n"))
	}
}
