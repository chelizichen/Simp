package main

import (
	h "Simp/handlers/http"
	"Simp/servers/CalcServer/service"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	ctx := h.NewSimpHttpCtx("simp.yaml")
	ctx.Use(service.Calc)
	ctx.Use(service.Plan)
	// static
	ctx.Use(func(engine *h.SimpHttpServerCtx) {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		webRoot := filepath.Join(wd, "web")
		engine.Engine.Static("/web", webRoot)
	})
	h.NewSimpHttpServer(ctx)
}
