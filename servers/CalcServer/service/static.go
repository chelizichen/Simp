package service

import (
	handlers "Simp/handlers/http"
	"os"
	"path/filepath"
)

func Static(ctx *handlers.SimpHttpServerCtx) {
	wd, _ := os.Getwd()
	SIMP_PRODUCTION := os.Getenv("SIMP_PRODUCTION")
	var staticPath string
	if SIMP_PRODUCTION == "Yes" {
		SIMP_SERVER_PATH := os.Getenv("SIMP_SERVER_PATH")
		staticPath = filepath.Join(SIMP_SERVER_PATH, ctx.StaticPath)
	} else {
		staticPath = filepath.Join(wd, ctx.StaticPath)
	}
	ctx.Engine.Static("/web", staticPath)
}
