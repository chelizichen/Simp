package main

import (
	h "Simp/src/http"
	service2 "Simp/src/service"
)

func main() {
	ctx := h.NewSimpHttpCtx("simp.yaml")
	ctx.DefineMain()
	ctx.Use(service2.Registry)
	ctx.Use(service2.Static)
	h.NewSimpHttpServer(ctx)
}
