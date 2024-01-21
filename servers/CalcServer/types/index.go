package types

import (
	"encoding/json"
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

type PaymentFrequency int

const (
	OneTime PaymentFrequency = iota
	PerDay
	Week
	Month
	Year
)

type PlanDetail struct {
	Income    float64          `json:"income,omitempty"`    // 预计收入
	OutCome   float64          `json:"outcome,omitempty"`   // 预计支出
	Style     PaymentFrequency `json:"style,omitempty"`     // 周期
	Comment   string           `json:"comment,omitempty"`   // 评价
	Sum       int              `json:"sum,omitempty"`       // 几个周期 0为完整周期、其他为指定周期
	StartTime string           `json:"startTime,omitempty"` // 开始时间 传空为覆盖完整周期
}

type PlanDTO struct {
	Id        int          `json:"id,omitempty"`
	Details   []PlanDetail `json:"details,omitempty"`   // 计划周期
	Comment   string       `json:"comment,omitempty"`   // 标注
	Name      string       `json:"name,omitempty"`      // 名称
	StartTime string       `json:"startTime,omitempty"` // 开始时间
	EndTime   string       `json:"endTime,omitempty"`   // 结束时间
}

// ST_Plan 数据库中存入的字段
type ST_Plan struct {
	Id        int    `db:"id,omitempty" json:"id,omitempty"`
	Comment   string `db:"comment,omitempty" json:"comment,omitempty"`      // 标注
	Name      string `db:"name,omitempty" json:"name,omitempty"`            // 名称
	StartTime string `db:"start_time,omitempty" json:"startTime,omitempty"` // 开始时间
	EndTime   string `db:"end_time,omitempty" json:"endRime,omitempty"`     // 结束时间
	Details   string `db:"details,omitempty" json:"details,omitempty"`      // 细节
}

func DTO2ST_PLAN(dto PlanDTO) *ST_Plan {
	b, err := json.Marshal(dto.Details)
	if err != nil {
		fmt.Println("Error On PlanDTO_ST ", err.Error())
	}
	return &ST_Plan{
		Comment:   dto.Comment,
		Name:      dto.Name,
		StartTime: dto.StartTime,
		EndTime:   dto.EndTime,
		Details:   string(b),
		Id:        dto.Id,
	}
}
