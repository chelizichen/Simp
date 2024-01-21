package service

import (
	handlers "Simp/handlers/http"
	"Simp/servers/CalcServer/storage"
	"Simp/servers/CalcServer/types"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Plan(ctx *handlers.SimpHttpServerCtx, pre string) {
	GROUP := ctx.Engine.Group(pre)
	// f := su.Join(pre)

	G := ctx.Engine
	ST := ctx.Storage
	GROUP.POST("/plan/create", func(ctx *gin.Context) {
		var requestBody types.PlanDTO
		if err := ctx.BindJSON(&requestBody); err != nil {
			ctx.JSON(http.StatusOK, handlers.Resp(0-1, "Error"+err.Error(), nil))
			return
		}
		sql, args := storage.SavePlan(requestBody)
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

	GROUP.POST("/plan/list", func(ctx *gin.Context) {
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

	GROUP.POST("/plan/update", func(ctx *gin.Context) {
		var requestBody types.PlanDTO
		if err := ctx.BindJSON(&requestBody); err != nil {
			ctx.JSON(http.StatusOK, handlers.Resp(0-1, "Error"+err.Error(), nil))
			return
		}
		fmt.Println("update requestBody", requestBody)
		sql, args := storage.UpdatePlan(requestBody)
		fmt.Println("UpdatePlan", sql, args)
		_, err := ST.DB.Exec(sql, args...) // 保存
		if err != nil {
			fmt.Println("UpdateError Error ", err.Error())
			ctx.JSON(http.StatusOK, handlers.Resp(-1, "Update Error", nil))
			return
		}
		ctx.JSON(http.StatusOK, handlers.Resp(0, "ok", nil))
	})

	GROUP.GET("/plan/detail", func(ctx *gin.Context) {
		id := ctx.Query("id")
		if id == "" {
			ctx.JSON(http.StatusOK, handlers.Resp(-1, "Id  Error", nil))
			return
		}
		i, err := strconv.Atoi(id)
		if err != nil {
			ctx.JSON(http.StatusOK, handlers.Resp(-1, "strconv Atoi", nil))
			return
		}
		sql, args := storage.DeleteById(i)
		_, err = ST.DB.Exec(sql, args...) // 保存
		if err != nil {
			ctx.JSON(http.StatusOK, handlers.Resp(-1, "Delete Error", nil))
			return
		}
		ctx.JSON(http.StatusOK, handlers.Resp(0, "ok", nil))
	})

	G.Use(GROUP.Handlers...)
}
