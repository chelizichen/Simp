package handlers

import (
	"Simp/config"
	"Simp/utils"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	CWD_ERROR               = "Error to Getwd"
	GetSubdirectories_Error = "Error to GetSubdirectories"
	NewConfig_Error         = "Get Config Error"
	ReadApi_Error           = "Read Api Json Error"
)

var ServantProviders = make(map[string]*ServantProvider, 512)

type Apis map[string]struct {
	Method string
	Path   string
}
type ServantProvider struct {
	Port       int
	ServerName string
	Apis       *Apis
}

type SimpHttpGateway struct {
}

// 主控根据服务名称直接寻址
// 1 可以更灵活的做限流
// 2 可以统一API网关
// 3 好统一做校验
func (s *SimpHttpGateway) Init() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(CWD_ERROR)
		panic(CWD_ERROR)
	}
	serverPath := filepath.Join(cwd, utils.PublishPath)
	subdirectories, err := utils.GetSubdirectories(serverPath)
	if err != nil {
		fmt.Println(GetSubdirectories_Error)
		panic(GetSubdirectories_Error)
	}
	for _, serverName := range subdirectories {
		s := &ServantProvider{}
		ServantProviders[serverName] = s
		// 对接生产Path
		servantConfigPath := filepath.Join(cwd, utils.PublishPath, serverName, "simpProd.yaml")
		conf, err := config.NewConfig(servantConfigPath)
		if err != nil {
			fmt.Println(NewConfig_Error)
			panic(NewConfig_Error)
		}
		ApisPath := filepath.Join(cwd, utils.PublishPath, serverName, "Api.json")

		var apis *Apis
		Content, err := os.ReadFile(ApisPath)
		if err != nil {
			fmt.Println(ReadApi_Error)
			panic(ReadApi_Error)
		}
		json.Unmarshal(Content, &apis)

		s.ServerName = serverName
		s.Port = conf.Server.Port
		s.Apis = apis
	}
}

type InvokeHeader struct {
	ServerName string
	Route      string
	Token      string
}

func (s *SimpHttpGateway) Invoke(header InvokeHeader, body any) {

}
