package service

import handlers "Simp/handlers/http"

func Gateway(ctx *handlers.SimpHttpServerCtx, pre string) {
	Engine := ctx.Engine
	Group := Engine.Group(pre)
	GateWay := &handlers.SimpHttpGateway{}
	GateWay.Use(Group)
	Engine.Use(Group.Handlers...)
}
