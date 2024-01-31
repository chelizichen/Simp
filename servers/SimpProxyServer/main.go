package main

import (
	"Simp/servers/SimpProxyServer/svr"
	h "Simp/src/http"
	"fmt"
)

func main() {
	ctx := h.NewSimpHttpCtx("simp.yaml")
	fmt.Println("ctx", ctx.StaticPath, ctx.Storage)
	ctx.Use(h.UseGateway)
	ctx.Use(svr.Service)
	svr.InitizalCacheSvr(ctx)
	h.NewSimpHttpServer(ctx)
}
