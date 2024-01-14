package service

import (
	handlers "Simp/handlers/http"
	"Simp/servers/CalcServer/types"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Plan(ctx *handlers.SimpHttpServerCtx) {
	G := ctx.Engine
	ST := ctx.Storage
	G.POST("/plan/create", func(ctx *gin.Context) {
		var requestBody types.PlanDTO
		if err := ctx.BindJSON(&requestBody); err != nil {
			ctx.JSON(http.StatusOK, handlers.Resp(0, "-1", err.Error()))
			return
		}
		s := types.DTO2ST_PLAN(requestBody)
		sql, args := types.SavePlan(*s)
		_, err := ST.Exec(sql, args)
		if err != nil {
			fmt.Println("SavePlan Error ", err.Error())
		}
		ctx.JSON(http.StatusOK, handlers.Resp(0, "ok", nil))
	})
}
