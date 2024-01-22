package types

import (
	"database/sql"
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
	Id        int              `json:"id,omitempty" db:"id"`
	Income    float64          `json:"income,omitempty" db:"income"`        // 预计收入
	OutCome   float64          `json:"outcome,omitempty" db:"outcome"`      // 预计支出
	Style     PaymentFrequency `json:"style,omitempty" db:"style"`          // 周期
	Comment   string           `json:"comment,omitempty" db:"comment"`      // 评价
	Sum       int              `json:"sum,omitempty" db:"sum"`              // 几个周期 0为完整周期、其他为指定周期
	StartTime string           `json:"startTime,omitempty" db:"start_time"` // 开始时间 传空为覆盖完整周期
	PlanId    int              `json:"planId,omitempty" db:"plan_id"`       // 开始时间 传空为覆盖完整周期
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
type PlanDAO struct {
	Id        int    `db:"id,omitempty" json:"id,omitempty"`
	Comment   string `db:"comment,omitempty" json:"comment,omitempty"`      // 标注
	Name      string `db:"name,omitempty" json:"name,omitempty"`            // 名称
	StartTime string `db:"start_time,omitempty" json:"startTime,omitempty"` // 开始时间
	EndTime   string `db:"end_time,omitempty" json:"endRime,omitempty"`     // 结束时间
}

func PlanConvert(dto PlanDTO) *PlanDAO {
	return &PlanDAO{
		Comment:   dto.Comment,
		Name:      dto.Name,
		StartTime: dto.StartTime,
		EndTime:   dto.EndTime,
		Id:        dto.Id,
	}
}

type PlanListDAO struct {
	Id         int             `db:"id"`
	Comment    sql.NullString  `db:"comment,omitempty"`
	Name       sql.NullString  `db:"name,omitempty"`
	StartTime  sql.NullString  `db:"start_time,omitempty"`
	EndTime    sql.NullString  `db:"end_time,omitempty"`
	Details    sql.NullString  `db:"details,omitempty"`
	SId        sql.NullInt64   `db:"s_id,omitempty"`
	SComment   sql.NullString  `db:"s_comment,omitempty"`
	Income     sql.NullFloat64 `db:"income,omitempty"`
	Outcome    sql.NullFloat64 `db:"outcome,omitempty"`
	SStartTime sql.NullString  `db:"s_start_time,omitempty"`
	Sum        sql.NullInt32   `db:"sum,omitempty"`
}
