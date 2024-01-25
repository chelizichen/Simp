package service

import handlers "Simp/handlers/http"

func Gateway(ctx *handlers.SimpHttpServerCtx, pre string) {
	G := ctx.Engine
	P := G.Group(pre)
	W := &handlers.SimpHttpGateway{}
	W.InitGateway()
	P.POST("/api/*route", W.GatewayMiddleWare)
	G.Use(P.Handlers...)
}
