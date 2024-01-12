package utils

import (
	"fmt"
)

type CalculateFutureValueDTO struct {
	Principal          float64 `json:"principal"`
	MonthlyDeposit     float64 `json:"monthlyDeposit"`
	AnnualInterestRate float64 `json:"annualInterestRate"`
	Months             int     `json:"months"`
}

type CalculateFutureValueVo struct {
	TotalAmount float64 `json:"totalAmount"`
	Profit      float64 `json:"profit"`
	OwnAmount   float64 `json:"ownAmount"`
}

// CalculateFutureValue 计算未来价值和利息
func CalculateFutureValue(principal, monthlyDeposit, annualInterestRate float64, months int) float64 {
	monthlyInterestRate := annualInterestRate / 12 / 100
	fmt.Println("args...", principal, monthlyDeposit, annualInterestRate, months, monthlyInterestRate)
	var futureValue float64 = principal

	for i := 0; i < months; i++ {
		futureValue += monthlyDeposit
		futureValue = (1 + monthlyInterestRate) * futureValue
	}

	return futureValue
}

func FixedIncomevalue(principal float64, months int, annualInterestRate float64) float64 {
	monthlyInterestRate := annualInterestRate / 12 / 100
	var futureValue float64 = principal
	for i := 0; i < months; i++ {
		futureValue = (1 + monthlyInterestRate) * futureValue
	}
	return futureValue
}
