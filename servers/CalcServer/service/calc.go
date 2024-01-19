package service

import (
	handlers "Simp/handlers/http"
	"Simp/servers/CalcServer/types"
	"Simp/servers/CalcServer/utils"
	su "Simp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Calc(ctx *handlers.SimpHttpServerCtx, pre string) {
	G := ctx.Engine
	f := su.Join(pre)

	G.POST(f("/calc/FixeIncome"), func(ctx *gin.Context) {
		var requestBody types.CalculateFutureValueDTO
		if err := ctx.BindJSON(&requestBody); err != nil {
			ctx.JSON(http.StatusOK, handlers.Resp(0, "-1", err.Error()))
			return
		}

		futureValue := utils.FixedIncomevalue(
			requestBody.Principal,
			requestBody.Months,
			requestBody.AnnualInterestRate,
		)

		responseBody := types.CalculateFutureValueVo{
			TotalAmount: futureValue,
		}

		responseBody.OwnAmount = requestBody.Principal
		responseBody.Profit = responseBody.TotalAmount - responseBody.OwnAmount

		ctx.JSON(http.StatusOK, handlers.Resp(0, "ok", responseBody))
	})
	G.POST(f("/calc/CalculateFutureValue"), func(ctx *gin.Context) {
		var requestBody types.CalculateFutureValueDTO
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

		responseBody := types.CalculateFutureValueVo{
			TotalAmount: futureValue,
		}

		responseBody.OwnAmount = requestBody.Principal + requestBody.MonthlyDeposit*float64(requestBody.Months)
		responseBody.Profit = responseBody.TotalAmount - responseBody.OwnAmount

		ctx.JSON(http.StatusOK, handlers.Resp(0, "ok", responseBody))
	})
}
