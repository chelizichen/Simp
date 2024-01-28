package main

import (
	"Simp/servers/CalcServer/service"
	h "Simp/src/http"
)

func main() {
	ctx := h.NewSimpHttpCtx("simp.yaml")
	ctx.Use(service.Calc)
	ctx.Use(service.Plan)
	// static
	ctx.Static("/web")
	h.NewSimpHttpServer(ctx)
}
