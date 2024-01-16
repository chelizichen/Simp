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
			ctx.JSON(http.StatusOK, handlers.Resp(0-1, "Error"+err.Error(), nil))
			return
		}
		s := types.DTO2ST_PLAN(requestBody)
		sql, args := types.SavePlan(*s)
		fmt.Println("ctx.ST is Nil", ST == nil)
		fmt.Println("sql", sql, "args", args)
		_, err := ST.DB.Exec(sql, args...) // 保存
		if err != nil {
			fmt.Println("SavePlan Error ", err.Error())
			ctx.JSON(http.StatusOK, handlers.Resp(-1, "saveError", nil))
			return
		}
		ctx.JSON(http.StatusOK, handlers.Resp(0, "ok", nil))
	})
}
