package main

import (
	h "Simp/handlers/http"
	"Simp/servers/CalcServer/service"
)

func main() {
	ctx := h.NewSimpHttpCtx("simp.yaml")
	ctx.Use(service.Calc)
	ctx.Use(service.Plan)
	// static
	ctx.Use(service.Static)
	h.NewSimpHttpServer(ctx)
}
