package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"
)

// 日志上报功能
// 根据服务名称 查询日志
type SimpMonitor struct {
	LogPath string
}

func NewSimpMonitor(serverName string, date string) (s SimpMonitor, e error) {
	// 判断date是否有值
	// 没有则传当天
	// path = static/serverName/log_date.log
	now := time.Now().Format(time.DateOnly)

	if date == "" {
		date = now
	} else {
		t, err2 := time.Parse(time.DateOnly, date)
		if err2 != nil {
			fmt.Println("time parse error", err2.Error())
		}
		date = t.Format(time.DateOnly)
	}

	fileName := "log_" + date + ".log"
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error To GetWd", err.Error())
	}

	path := path.Join(cwd, PublishPath, serverName, fileName)
	fmt.Println("Logger Path", path)
	s.LogPath = path
	b := IsExist(path)
	if !b {
		F, err := os.Create(path)
		if err != nil {
			fmt.Println("Error creating log file:", err)
		}
		F.WriteString("// SimpLog Created on " + path + " at " + now + "\n")
	}
	return s, nil
}

// 返回一个Logger
func (s *SimpMonitor) GetLogger(pattern string) (string, error) {
	cmd := exec.Command("grep", pattern, s.LogPath)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stderr.String(), err
	}

	return stdout.String(), nil
}

// 返回一个Logger
func (s *SimpMonitor) AppendLogger(content string) {
	file, err := os.OpenFile(s.LogPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(content + "\n"); err != nil {
		fmt.Println("Error writing to log file:", err)
	}
}

func AutoCreateLoggerFile(serverName string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error To GetWd", err.Error())
	}
	tdy := time.Now().Format(time.DateOnly)
	fileName := "log_" + tdy + ".log"
	path := path.Join(cwd, PublishPath, serverName, fileName)
	b := IsExist(path)
	if b {
		return
	}
	F, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating log file:", err)
	}
	F.WriteString("// SimpLog Created on " + path + " at " + tdy + "\n")

}
