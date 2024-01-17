package main

import (
	h "Simp/handlers/http"
	"Simp/servers/BlogServer/service"
)

func main() {
	ctx := h.NewSimpHttpCtx("simp.yaml")
	ctx.Use(service.BlogService)
	ctx.Static("/web")
	h.NewSimpHttpServer(ctx)
}
