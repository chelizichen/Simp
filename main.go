package main

import (
	h "Simp/handlers/http"
	"Simp/service"
)

func main() {
	ctx := h.NewSimpHttpCtx("simp.yaml")
	ctx.DefineMain()
	ctx.Use(service.Registry)
	ctx.Use(service.Static)
	h.NewSimpHttpServer(ctx)
}
