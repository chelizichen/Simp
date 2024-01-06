package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func startChildService() (*exec.Cmd, io.WriteCloser, io.Reader, error) {
	var wd, err = os.Getwd()
	child := wd + "/test/fass/child/TestFassChild"
	fmt.Println("FASS", child)
	cmd := exec.Command(child)

	// 创建标准输入管道
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	// 创建标准输出管道
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, nil, nil, fmt.Errorf("error starting child service: %v", err)
	}

	// 等待一段时间，确保子服务已经启动
	time.Sleep(2 * time.Second)
	return cmd, stdinPipe, stdoutPipe, nil
}

func main() {
	childCmd, stdinPipe, stdoutPipe, err := startChildService()
	if err != nil {
		fmt.Println("Error starting child service:", err)
		return
	}
	defer stdinPipe.Close()

	// 模拟手动输入
	inputData := "Some test input\n"
	_, err = stdinPipe.Write([]byte(inputData))
	if err != nil {
		fmt.Println("Error writing to child process:", err)
		return
	}

	// 读取子进程的输出
	scanner := bufio.NewScanner(stdoutPipe)
	for scanner.Scan() {
		outputLine := scanner.Text()
		fmt.Println("Child Process Output:", outputLine)
	}

	// 在这里可以继续处理其他逻辑
	// ...

	// 等待子服务退出
	err = childCmd.Wait()
	if err != nil {
		fmt.Println("Error waiting for child service:", err)
	}
}
