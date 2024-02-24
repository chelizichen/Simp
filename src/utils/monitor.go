package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

// 日志上报功能
// 根据服务名称 查询日志
type SimpMonitor struct {
	LogPath string
}

func NewSearchLogMonitor(serverName string, fileName string) (s SimpMonitor, e error) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error To GetWd", err.Error())
	}
	path := path.Join(cwd, PublishPath, serverName, fileName)
	s.LogPath = path
	b := IsExist(path)
	if !b {
		return s, errors.New("Error! File is Not Exist")
	}
	return s, nil
}

func NewMainSearchLogMonitor(fileName string) (s SimpMonitor, e error) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error To GetWd", err.Error())
	}
	path := path.Join(cwd, "static/main", fileName)
	s.LogPath = path
	b := IsExist(path)
	if !b {
		return s, errors.New("Error! File is Not Exist")
	}
	return s, nil
}

func NewSimpMonitor(serverName string, date string) (s SimpMonitor, e error) {
	// 判断date是否有值
	// 没有则传当天
	// path = static/serverName/log_date.log
	now := time.Now().Format(time.DateOnly)
	fmt.Println("Today", now)
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
	defer F.Close()
	if err != nil {
		fmt.Println("Error creating log file:", err)
	}
	F.WriteString("// SimpLog Created on " + path + " at " + tdy + "\n")

}

type APIFileVo struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

func CreateAPIFile(c *gin.Engine, serverName string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error To GetWd", err.Error())
	}

	// 获取所有监听的路由
	routes := c.Routes()
	var APIs []APIFileVo

	// 打印路由信息
	for _, route := range routes {
		vo := APIFileVo{
			Method: route.Method,
			Path:   route.Path,
		}
		APIs = append(APIs, vo)
		fmt.Printf("Method: %s, Path: %s  \n", route.Method, route.Path)
	}

	B, err := json.Marshal(APIs)
	if err != nil {
		fmt.Println("Errored On Json Marshal", err.Error())
	}
	path := path.Join(cwd, PublishPath, serverName, "API.json")
	fmt.Println("Create JSON at ", path)
	err = IFExistThenRemove(path)
	if err != nil {
		fmt.Println("IFExistThenRemove Error ", path)
	}
	F, err := os.Create(path)
	defer F.Close()
	if err != nil {
		fmt.Println("Create JSON API File Error", err.Error())
		return
	}
	F.WriteString(string(B))
}

// 主控服务用
type LogWriter struct {
	File *os.File
}

func NewLogWriter(filePath string) (*LogWriter, error) {
	F, err := IFNotExistThenCreate(filePath)
	if err != nil {
		return nil, err
	}

	return &LogWriter{File: &F}, nil
}

func (lw *LogWriter) Write(p []byte) (n int, err error) {
	s := time.Now().Format(time.Stamp) + " " + string(p)
	_, err = lw.File.WriteString(s)
	return len(p), err
}

func GetServerLogName() string {
	tdy := time.Now().Format(time.DateOnly)
	return "server_" + tdy + ".log"
}

func AutoSetLogWriter() {
	cwd, err := os.Getwd()
	if err != nil {
		panic("GetWd Error" + err.Error())
	}
	logFilePath := path.Join(cwd, "static/main", GetServerLogName())

	logWriter, err := NewLogWriter(logFilePath)
	if err != nil {
		fmt.Println("创建文件失败:", err)
		return
	}

	// 重定向标准输出到自定义的写入器
	os.Stdout = logWriter.File
}
