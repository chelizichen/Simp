package handlers

import (
	"Simp/config"
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

// FassServants 最多创建512 个 Fass服务
var FassServants = make(map[string]*SimpFassHandler, 512)

type SimpFassHandler struct {
	CMD    *exec.Cmd
	WRITER io.WriteCloser
	READER io.Reader
	CONF   *config.SimpConfig
	pid    int
	isOpen bool
}

func NewIoHandler(conf config.SimpConfig) (handler *SimpFassHandler, errs error) {
	var wd, err = os.Getwd()
	child := wd + conf.Server.Name
	fmt.Println("standard log | start fass server", child)
	cmd := exec.Command(child)
	// 创建标准输入管道
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return handler, err
	}

	// 创建标准输出管道
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return handler, err
	}

	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return handler, fmt.Errorf("error log | starting child service: %v", err)
	}
	handler.CONF = &conf
	handler.CMD = cmd
	handler.READER = stdoutPipe
	handler.WRITER = stdinPipe
	// 等待一段时间，确保子服务已经启动
	time.Sleep(2 * time.Second)
	handler.isOpen = true
	handler.pid = cmd.Process.Pid
	FassServants[conf.Server.Name] = handler
	return handler, nil
}

func (i *SimpFassHandler) Write(data string) (int, error) {
	return i.WRITER.Write([]byte(data))
}

func (i *SimpFassHandler) Read() {
	scanner := bufio.NewScanner(i.READER)
	for scanner.Scan() {
		outputLine := scanner.Text()
		fmt.Println("Child Process Output:", outputLine)
	}
	err := i.CMD.Wait()
	if err != nil {
		fmt.Println("Error waiting for child service:", err)
	}
}

func (i *SimpFassHandler) ReadStatus() (pid int, isOpen bool) {
	if i.pid != 0 && i.isOpen == true {
		return i.pid, true
	} else {
		return 0, false
	}
}
