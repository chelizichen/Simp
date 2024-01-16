package service

import (
	handlers "Simp/handlers/http"
	"Simp/servers/CalcServer/storage"
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
		sql, args := storage.SavePlan(*s)
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

	G.POST("/plan/list", func(ctx *gin.Context) {
		var Resp []types.ST_Plan
		sql := storage.GetList()
		err := ST.Select(&Resp, sql) // 保存
		if err != nil {
			fmt.Println("Get List Error", err.Error())
			ctx.JSON(http.StatusOK, handlers.Resp(-1, "Get List Error", nil))
			return
		}
		ctx.JSON(http.StatusOK, handlers.Resp(0, "ok", Resp))
	})

}
