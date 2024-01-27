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
	"strings"

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

func (s *SimpHttpGateway) Use(svr *gin.RouterGroup, ctx *SimpHttpServerCtx) {
	var serverPath string
	serverPath = filepath.Join(ctx.StaticPath, utils.PublishPath)
	// 获取服务目录
	fmt.Println("ServerPath |", serverPath)
	subdirectories, err := utils.GetSubdirectories(serverPath)
	if err != nil {
		fmt.Println(GetSubdirectories_Error)
		panic(GetSubdirectories_Error)
	}
	// 遍历映射服务详情与API
	for _, serverName := range subdirectories {
		fmt.Println("serverName", serverName, " || ctx.name", ctx.Name)
		if strings.ToLower(serverName) == strings.ToLower(ctx.Name) {
			return
		}
		s := &ServantProvider{}
		ServantProviders[serverName] = s
		// 对接生产Path
		servantConfigPath := filepath.Join(ctx.StaticPath, utils.PublishPath, serverName, "simpProd.yaml")
		conf, err := config.NewConfig("", servantConfigPath)
		if err != nil {
			fmt.Println(NewConfig_Error + servantConfigPath)
			panic(NewConfig_Error + servantConfigPath)
		}
		ApisPath := filepath.Join(ctx.StaticPath, utils.PublishPath, serverName, "API.json")

		var apis *Apis
		Content, err := os.ReadFile(ApisPath)
		if err != nil {
			fmt.Println(ReadApi_Error + err.Error())
			panic(ReadApi_Error + err.Error())
		}
		err = json.Unmarshal(Content, &apis)
		if err != nil {
			fmt.Println(ReadApi_Error + err.Error())
			panic(ReadApi_Error + err.Error())
		}
		s.ServerName = serverName
		s.Port = conf.Server.Port
		s.Apis = apis
		s.Host = conf.Server.Host
	}
	for serverName, Provider := range ServantProviders {
		fmt.Println("serverName", serverName, " is load success")
		for _, V := range *Provider.Apis {
			v := V
			if v.Method == "POST" {
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
						fmt.Println("Error reading response body:", err.Error())
						ctx.JSON(http.StatusInternalServerError, gin.H{"error": "读取响应体失败"})
						return
					}
					var anyBody interface{}
					err = json.Unmarshal(responseBody, &anyBody)
					if err != nil {
						fmt.Println("Error To Unmarshal json:", err.Error())
						ctx.JSON(http.StatusInternalServerError, gin.H{"error": "解码响应体失败"})
						return
					}
					ctx.JSON(resp.StatusCode, anyBody)
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
						fmt.Println("Error reading response body:", err.Error())
						ctx.JSON(http.StatusInternalServerError, gin.H{"error": "读取响应体失败"})
						return
					}
					var anyBody interface{}
					err = json.Unmarshal(responseBody, &anyBody)
					if err != nil {
						fmt.Println("Error To Unmarshal json:", err.Error())
						ctx.JSON(http.StatusInternalServerError, gin.H{"error": "解码响应体失败"})
						return
					}
					ctx.JSON(resp.StatusCode, anyBody)
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
	fmt.Println("target ", url)
	return url
}

func UseGateway(ctx *SimpHttpServerCtx, pre string) {
	Engine := ctx.Engine
	Group := Engine.Group(pre)
	GateWay := &SimpHttpGateway{}
	GateWay.Use(Group, ctx)
}
