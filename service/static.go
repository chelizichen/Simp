package service

import (
	handlers "Simp/handlers/http"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func Static(ctx *handlers.SimpHttpServerCtx) {
	G := ctx.Engine
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error To GetWd", err.Error())
	}
	staticRoot := filepath.Join(wd, "static")
	webRoot := filepath.Join(wd, "pages")
	G.Static("/static", staticRoot)
	G.Static("/web", webRoot)
	G.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(304, "/web")
	})
}
