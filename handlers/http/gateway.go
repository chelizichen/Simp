package handlers

import (
	"Simp/config"
	"Simp/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const (
	CWD_ERROR               = "Error to Getwd"
	GetSubdirectories_Error = "Error to GetSubdirectories"
	NewConfig_Error         = "Get Config Error"
	ReadApi_Error           = "Read Api Json Error"
)

var ServantProviders = make(map[string]*ServantProvider, 512)

type Apis []struct {
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
func (s *SimpHttpGateway) Use(P *gin.RouterGroup) {

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(CWD_ERROR)
		panic(CWD_ERROR)
	}
	serverPath := filepath.Join(cwd, utils.PublishPath)
	// 获取服务目录
	subdirectories, err := utils.GetSubdirectories(serverPath)
	if err != nil {
		fmt.Println(GetSubdirectories_Error)
		panic(GetSubdirectories_Error)
	}
	// 遍历映射服务详情与API
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
		err = json.Unmarshal(Content, &apis)
		if err != nil {
			fmt.Println(ReadApi_Error + err.Error())
			panic(ReadApi_Error + err.Error())
		}
		s.ServerName = serverName
		s.Port = conf.Server.Port
		s.Apis = apis
	}
	fmt.Println("ServantProviders", ServantProviders)
	for serverName, Provider := range ServantProviders {
		svr := P.Group(serverName)
		fmt.Println("serverName", serverName)
		for _, v := range *Provider.Apis {
			if v.Method == "POST" {
				fmt.Println("v.Path", v.Path, "v.Method", v.Method)
				svr.POST(v.Path, func(ctx *gin.Context) {
					route := v.Path
					target := Provider.GetTarget(route) // 调用 proxy
					client := &http.Client{}
					// 读取请求体
					bodyBytes, err := io.ReadAll(ctx.Request.Body)
					if err != nil {
						// 处理错误
						return
					}
					// 创建请求体可读对象
					requestBody := io.NopCloser(bytes.NewReader(bodyBytes))
					// 发送请求
					resp, err := client.Post(target, ctx.ContentType(), requestBody)
					if err != nil {
						// 处理错误
						ctx.JSON(http.StatusInternalServerError, gin.H{"error": "发送请求失败"})
						return
					}
					defer resp.Body.Close()

					// 处理响应
					// 读取响应体，处理状态码等
					responseBody, err := io.ReadAll(resp.Body)
					if err != nil {
						// 处理错误
						ctx.JSON(http.StatusInternalServerError, gin.H{"error": "读取响应体失败"})
						return
					}

					// 假设你想将响应体直接返回给客户端
					ctx.String(resp.StatusCode, string(responseBody))
				})
			}
			if v.Method == "GET" {
				svr.GET(v.Path, func(ctx *gin.Context) {
					route := v.Path
					target := Provider.GetTarget(route) // 调用 proxy
					client := &http.Client{}

					// 读取请求的 query 参数
					queryParams := ctx.Request.URL.Query()
					queryString := queryParams.Encode()

					// 发送 GET 请求
					resp, err := client.Get(target + "?" + queryString)
					if err != nil {
						// 处理错误
						ctx.JSON(http.StatusInternalServerError, gin.H{"error": "发送请求失败"})
						return
					}
					defer resp.Body.Close()

					// 处理响应
					// 读取响应体，处理状态码等
					responseBody, err := io.ReadAll(resp.Body)
					if err != nil {
						// 处理错误
						ctx.JSON(http.StatusInternalServerError, gin.H{"error": "读取响应体失败"})
						return
					}

					// 假设你想将响应体直接返回给客户端
					ctx.String(resp.StatusCode, string(responseBody))
				})
			}
		}
	}
}

// GetTarget 获取调用方
// host := "example.com"
// port := 8080
// route := "/api/data"
func (s *ServantProvider) GetTarget(route string) string {
	url := fmt.Sprintf("http://%s:%d%s", s.Host, s.Port, route)
	return url
}
