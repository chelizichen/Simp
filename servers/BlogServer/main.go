package main

import (
	"Simp/servers/BlogServer/service"
	h "Simp/src/http"
)

func main() {
	ctx := h.NewSimpHttpCtx("simp.yaml")
	ctx.Use(service.BlogService)
	ctx.Static("/web")
	h.NewSimpHttpServer(ctx)
}
