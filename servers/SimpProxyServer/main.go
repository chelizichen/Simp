package main

import (
	h "Simp/handlers/http"
	"fmt"
)

func main() {
	ctx := h.NewSimpHttpCtx("simp.yaml")
	fmt.Println("ctx", ctx.StaticPath, ctx.Storage)
	ctx.Use(h.UseGateway)
	h.NewSimpHttpServer(ctx)
}
