package storage

import (
	"Simp/servers/CalcServer/types"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func SavePlan(dto types.PlanDTO) (sql string, args []any) {
	sql = `insert into calc_plan 
		(
			name,comment,start_time,end_time
		) 
		values 
		(
			?,	?,	?,	?
		)`
	args = append(args, dto.Name, dto.Comment, dto.StartTime, dto.EndTime)
	return sql, args
}

func SaveSubPlan(req types.PlanDetail) (sql string, args []any) {
	sql = `insert into calc_sub_plan
		(
			comment, plan_id, income, outcome, start_time, sum
		)
		values
		(
			?,	?,	?,	?,	?,?
		)`
	args = append(args, req.Comment, req.PlanId, req.Income, req.OutCome, req.StartTime, req.Sum)
	return sql, args
}

func GetList() (sql string) {
	sql = " select c.*,s.id as s_id,s.`comment` as s_comment, s.income,s.outcome,s.start_time as s_start_time,s.sum from calc_plan c left join calc_sub_plan s on s.plan_id = c.id"
	return sql
}

func DeleteById(id int) (sql string, args []any) {
	sql = `delete from calc_plan where id = `
	args = append(args, id)
	return sql, args
}

func UpdatePlan(ST *sqlx.DB, req types.PlanDTO) error {
	var args []interface{}
	updateSql := `update calc_plan set 
	name = ?,
	comment = ?,
	start_time = ?,
	end_time = ?
	where id = ?
	`
	args = append(args, req.Name, req.Comment, req.StartTime, req.EndTime, req.Id)
	fmt.Println("UpdatePlan", updateSql, args)
	_, err := ST.DB.Exec(updateSql, args...) // 保存
	if err != nil {
		fmt.Println("UpdateError Error ", err.Error())
		return err
	}
	deleteSql := `delete from calc_sub_plan where plan_id = ?`
	_, err = ST.DB.Exec(deleteSql, req.Id) // 保存

	if err != nil {
		fmt.Println("Delete Error ", err.Error())
		return err
	}

	for _, v := range req.Details {
		v.PlanId = req.Id
		subUpdateSql, args := SaveSubPlan(v)
		_, err = ST.DB.Exec(subUpdateSql, args...) // 保存
		fmt.Println("subUpdateSql", subUpdateSql, args)
		if err != nil {
			fmt.Println("subUpdateSql Error ", err.Error())
			return err
		}
	}
	return nil
}
