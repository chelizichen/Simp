package handlers

import (
	"Simp/config"
	"Simp/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
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
	Host       string
	Apis       *Apis
}

type SimpHttpGateway struct {
}

// InitGateway
// 主控根据服务名称直接寻址
// 1 可以更灵活的做限流
// 2 可以统一API网关
// 3 好统一做校验
func (s *SimpHttpGateway) InitGateway() {
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
		servantConfigPath := filepath.Join(utils.PublishPath, serverName, "simpProd.yaml")
		conf, err := config.NewConfig(servantConfigPath)
		if err != nil {
			fmt.Println(NewConfig_Error + servantConfigPath)
			panic(NewConfig_Error + servantConfigPath)
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

type InvokeBody struct {
	ServerName string      `json:"serverName,omitempty"`
	Route      string      `json:"route,omitempty"`
	Token      string      `json:"token,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

// GetTarget 获取调用方
// host := "example.com"
// port := 8080
// route := "/api/data"
func (s *ServantProvider) GetTarget(route string) string {
	url := fmt.Sprintf("http://%s:%d%s", s.Host, s.Port, route)
	return url
}

func (s *SimpHttpGateway) Invoke(header *InvokeBody) *http.Response {
	provider := ServantProviders[header.ServerName]
	route := header.Route
	target := provider.GetTarget(route)
	client := &http.Client{}
	marshal, err := json.Marshal(header.Data)
	if err != nil {
		fmt.Println("JSON stringify Error:", err)
		return nil
	}
	var R io.Reader = bytes.NewBuffer(marshal)
	resp, err := client.Post(target, "application/json", R)

	if err != nil {
		fmt.Println("创建请求失败:", err)
		return nil
	}
	return resp
}

func (s *SimpHttpGateway) GatewayMiddleWare(c *gin.Context) {
	var invokeBody *InvokeBody
	err := c.BindJSON(invokeBody)
	if err != nil {
		fmt.Println("Error To BindJson")
	}
	go func() {
		invoke := s.Invoke(invokeBody)
		c.JSON(http.StatusOK, invoke)
	}()
}
