package service

import (
	handlers "Simp/handlers/http"
	"Simp/servers/CalcServer/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Calc(ctx *handlers.SimpHttpServerCtx) {
	G := ctx.Engine
	G.POST("/calc/FixeIncome", func(ctx *gin.Context) {
		var requestBody utils.CalculateFutureValueDTO
		if err := ctx.BindJSON(&requestBody); err != nil {
			ctx.JSON(http.StatusOK, handlers.Resp(0, "-1", err.Error()))
			return
		}

		futureValue := utils.FixedIncomevalue(
			requestBody.Principal,
			requestBody.Months,
			requestBody.AnnualInterestRate,
		)

		responseBody := utils.CalculateFutureValueVo{
			TotalAmount: futureValue,
		}

		responseBody.OwnAmount = requestBody.Principal
		responseBody.Profit = responseBody.TotalAmount - responseBody.OwnAmount

		ctx.JSON(http.StatusOK, handlers.Resp(0, "ok", responseBody))
	})
	G.POST("/calc/CalculateFutureValue", func(ctx *gin.Context) {
		var requestBody utils.CalculateFutureValueDTO
		if err := ctx.BindJSON(&requestBody); err != nil {
			ctx.JSON(http.StatusOK, handlers.Resp(0, "-1", err.Error()))
			return
		}

		futureValue := utils.CalculateFutureValue(
			requestBody.Principal,
			requestBody.MonthlyDeposit,
			requestBody.AnnualInterestRate,
			requestBody.Months,
		)

		responseBody := utils.CalculateFutureValueVo{
			TotalAmount: futureValue,
		}

		responseBody.OwnAmount = requestBody.Principal + requestBody.MonthlyDeposit*float64(requestBody.Months)
		responseBody.Profit = responseBody.TotalAmount - responseBody.OwnAmount

		ctx.JSON(http.StatusOK, handlers.Resp(0, "ok", responseBody))
	})
}
