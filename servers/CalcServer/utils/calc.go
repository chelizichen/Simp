package utils

import (
	"fmt"
)

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
