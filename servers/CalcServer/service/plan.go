package service

import (
	"Simp/servers/CalcServer/storage"
	"Simp/servers/CalcServer/types"
	"Simp/servers/CalcServer/utils"
	handlers "Simp/src/http"
	"fmt"
	"net/http"
	"sort"
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
		resp, err := ST.DB.Exec(sql, args...) // 保存
		if err != nil {
			fmt.Println("SavePlan Error ", err.Error())
			ctx.JSON(http.StatusOK, handlers.Resp(-1, "saveError", nil))
			return
		}
		id, err := resp.LastInsertId()
		if err != nil {
			fmt.Println("LastInsertId Error ", err.Error())
			ctx.JSON(http.StatusOK, handlers.Resp(-1, "saveError", nil))
			return
		}

		for _, sub := range requestBody.Details {
			sub.Id = int(id)
			saveSql, args := storage.SaveSubPlan(sub)
			_, err := ST.DB.Exec(saveSql, args...) // 保存
			if err != nil {
				fmt.Println("insert sub plan Error ", err.Error())
				ctx.JSON(http.StatusOK, handlers.Resp(-2, "saveError", nil))
				return
			}
		}

		ctx.JSON(http.StatusOK, handlers.Resp(0, "ok", nil))
	})

	GROUP.POST("/plan/list", func(ctx *gin.Context) {
		querySql := storage.GetList()
		var Resp []types.PlanListDAO
		var RespVo []types.PlanDTO
		err := ST.Select(&Resp, querySql) // 保存
		if err != nil {
			fmt.Println("Get List Error", err.Error())
			ctx.JSON(http.StatusOK, handlers.Resp(-1, "Get List Error", nil))
			return
		}
		groupResp := utils.GroupBy(Resp, func(p types.PlanListDAO) int {
			return p.Id
		})
		for k, v := range groupResp {
			ret := &types.PlanDTO{}
			ret.Id = k
			ret.Comment = v[0].Comment.String
			ret.Name = v[0].Name.String
			ret.StartTime = v[0].StartTime.String
			ret.EndTime = v[0].EndTime.String
			for _, detail := range v {
				if !detail.SId.Valid {
					continue
				}
				d := &types.PlanDetail{}
				d.Comment = detail.SComment.String
				d.Id = int(detail.SId.Int64)
				d.OutCome = float64(detail.Outcome.Float64)
				d.Income = float64(detail.Income.Float64)
				d.StartTime = detail.SStartTime.String
				d.Sum = int(detail.Sum.Int32)
				ret.Details = append(ret.Details, *d)
			}
			RespVo = append(RespVo, *ret)
		}
		sort.Slice(RespVo, func(i, j int) bool {
			return RespVo[i].Id < RespVo[j].Id
		})
		ctx.JSON(http.StatusOK, handlers.Resp(0, "ok", RespVo))
	})

	GROUP.POST("/plan/update", func(ctx *gin.Context) {
		var requestBody types.PlanDTO
		if err := ctx.BindJSON(&requestBody); err != nil {
			ctx.JSON(http.StatusOK, handlers.Resp(-1, "Error"+err.Error(), nil))
			return
		}
		fmt.Println("update requestBody", requestBody)
		err := storage.UpdatePlan(ST, requestBody)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error"+err.Error(), nil))
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
