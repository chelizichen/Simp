package service

import (
	handlers "Simp/handlers/http"
	"fmt"
	"os"
	"path/filepath"

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

func Static(ctx *handlers.SimpHttpServerCtx) {
	G := ctx.Engine
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error To GetWd", err.Error())
	}
	staticRoot := filepath.Join(wd, "static")
	webRoot := filepath.Join(wd, "pages")
	G.Use(CORSMiddleware())
	G.Static("/static", staticRoot)
	G.Static("/web", webRoot)
	// G.GET("/", func(ctx *gin.Context) {
	// 	ctx.Redirect(304, "/web")
	// })
}
