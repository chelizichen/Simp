package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/robfig/cron"
)

const (
	PublishPath         = "static/publish/"
	GoEntry             = "service_go"
	NodeJsEntry         = "app.js"
	DevConfEntry        = "simp.yaml"
	ProdConfEntry       = "simpProd.yaml"
	RELEASE_CLUSTER     = "RELEASE_CLUSTER"
	RELEASE_SINGLENODE  = "RELEASE_SINGLENODE"
	RELEASE_TYPE_GO     = "go-http"
	RELEASE_TYPE_NODEJS = "node-http"
	RELEASE_TYPE_JAVA   = "java-http"
)

func GetFilePath(cwd, serverName, fileName string) string {
	return filepath.Join(cwd, PublishPath, serverName, fileName)
}

type Servants struct {
	Language   string
	Pid        int
	ServerName string
	Port       string
	Process    *exec.Cmd
	ExitSignal atomic.Bool
	Cron       *cron.Cron
}

func NewServant() (svr *Servants) {
	return &Servants{}
}

func GetServant(ServerName string, port string) Servants {
	svr := ServerName + port
	servant := SubServants[svr]
	servant.ServerName = ServerName
	servant.Port = port
	return servant
}

func (s *Servants) GetContextName() string {
	return s.ServerName + s.Port
}

func (s *Servants) ServantMonitor() (string, bool) {
	time.Sleep(time.Second * 5)
	b := IsPidAlive(s.Pid)
	if !b {
		fmt.Println(s.GetContextName(), "isPidAlive ", s.Pid, "|", b)
		return "", false
	}
	pInfo := GetProcessMemoryInfo(s.Pid)
	cpuPercent, _ := pInfo.CPUPercent()
	cpuAffinity, _ := pInfo.CPUAffinity()
	createTime, _ := pInfo.CreateTime()
	Status, _ := pInfo.Status()
	if Status == "Z" {
		s.StopServant()
		return "", false
	}
	pid := s.Pid
	MemoryPercent, _ := pInfo.MemoryPercent()
	MemoryInfo, _ := pInfo.MemoryInfo()
	info := make(map[string]interface{})
	info["pid"] = pid
	info["MemoryInfo"] = MemoryInfo
	info["MemoryPercent"] = MemoryPercent
	info["CpuPercent"] = cpuPercent
	info["CpuAffinity"] = cpuAffinity
	info["CreateTime"] = createTime
	info["Status"] = Status
	info["ServerName"] = s.ServerName
	pInfoContent, err := json.Marshal(info)
	if err != nil {
		return "", false
	}
	t := time.Now().Format(time.DateTime)
	content := t + " ServerName " + s.ServerName + " || " + string(pInfoContent) + "\n"
	return content, true
}

func (s *Servants) StopServant() error {
	if s.Pid == 0 {
		return nil
	}
	s.Cron.Stop()
	s.ExitSignal.Store(true)
	err := s.Process.Process.Kill() // 终止进程
	if err != nil {
		fmt.Println("Killed Error", err.Error())
	}
	return err
}

func CopyProdYml(storageYmlEPath, storageYmlProdPath string) (err error) {
	_, err = os.Stat(storageYmlProdPath)
	if err != nil {
		fmt.Println("os.Stat ", err.Error())
	}
	// 如果没有该文件，则将simp.yaml拷贝一份成simpProd.yaml
	if os.IsNotExist(err) {
		err = CopyFile(storageYmlEPath, storageYmlProdPath)
		if err != nil {
			fmt.Println("utils.CopyFile ", storageYmlEPath, err.Error())
		}
	}
	return err
}

var SubServants = make(map[string]Servants)
