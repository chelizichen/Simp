package handlers

import (
	"Simp/config"
	"bufio"
	"fmt"
	"os"
)

type SimpFassServiceHandler struct {
	CONF *config.SimpConfig
}

func (h SimpFassServiceHandler) Handle() {
	reader := bufio.NewReader(os.Stdin)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		// 处理收到的消息（这里简单打印）
		fmt.Println("Received:", message)
	}
}
