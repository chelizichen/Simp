package main

import (
	h "Simp/handlers/http"
	"Simp/servers/SimpProxyServer/svr"
	"fmt"
)

func main() {
	ctx := h.NewSimpHttpCtx("simp.yaml")
	fmt.Println("ctx", ctx.StaticPath, ctx.Storage)
	ctx.Use(h.UseGateway)
	svr.InitizalCacheSvr(ctx)
	h.NewSimpHttpServer(ctx)
}
