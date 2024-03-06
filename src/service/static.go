package service

import (
	handlers "Simp/src/http"
	"Simp/src/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func Static(ctx *handlers.SimpHttpServerCtx, pre string) {
	f := utils.Join(pre)
	G := ctx.Engine
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error To GetWd", err.Error())
	}
	G.Use(func(c *gin.Context) {
		if strings.Index(c.Request.URL.Path, "simpserver/static") > -1 {
			// 对特定路由的静态资源设置 Cache-Control
			c.Writer.Header().Set("Cache-Control", "public, max-age=3600")
		}
		c.Next()
	})
	staticRoot := filepath.Join(wd, "static")
	webRoot := filepath.Join(wd, "dist")
	// webSPAAssets := filepath.Join(wd, "dist/assets")
	G.Use(CORSMiddleware())
	G.Static(f("/static"), staticRoot)
	G.GET(f("/web/*path"), func(ctx *gin.Context) {
		// 获取请求的路径
		requestPath := ctx.Param("path")

		// 拼接请求的文件路径
		filePath := filepath.Join(webRoot, requestPath)

		// 检查文件是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// 文件不存在，返回 404
			// ctx.Status(http.StatusNotFound)
			filePath = filepath.Join(webRoot, "index.html")
			ctx.File(filePath)
			return
		}
		// 返回文件
		ctx.File(filePath)
	})
}
