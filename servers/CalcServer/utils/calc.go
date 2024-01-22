package utils

import (
	"Simp/servers/CalcServer/types"
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

func KeyBy(arr []types.PlanListDAO, keyFunc func(p types.PlanListDAO) int) map[int]types.PlanListDAO {
	result := make(map[int]types.PlanListDAO)
	for _, item := range arr {
		key := keyFunc(item)
		result[key] = item
	}
	return result
}

func GroupBy[T any, K comparable](arr []T, keyFunc func(p T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range arr {
		key := keyFunc(item)
		result[key] = append(result[key], item)
	}
	return result
}
