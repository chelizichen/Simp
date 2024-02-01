package main

import (
	"Simp/servers/SimpProxyServer/svr"
	h "Simp/src/http"
	"fmt"
)

func main() {
	ctx := h.NewSimpHttpCtx("simp.yaml")
	svr.InitizalCacheSvr(ctx)
	fmt.Println("ctx", ctx.StaticPath, ctx.Storage)
	ctx.Use(h.UseGateway)
	ctx.Use(svr.Service)
	test(ctx)
	h.NewSimpHttpServer(ctx)
}

func test(ctx *h.SimpHttpServerCtx) {
	ctx.CacheSvr.Set("test", "1111")
	i, b := ctx.CacheSvr.Get("test")
	if b {
		fmt.Println("test", i)
		ctx.CacheSvr.Del("test")
	}
}
