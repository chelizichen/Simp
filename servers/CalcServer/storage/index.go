package storage

import "Simp/servers/CalcServer/types"

func SavePlan(dto types.PlanDTO) (sql string, args []any) {
	req := types.DTO2ST_PLAN(dto)
	sql = `insert into calc_plan 
		(
			name,comment,start_time,end_time,details
		) 
		values 
		(
			?,	?,	?,	?,	?
		)`
	args = append(args, req.Name, req.Comment, req.StartTime, req.EndTime, req.Details)
	return sql, args
}

func GetList() (sql string) {
	sql = `select * from calc_plan `
	return sql
}

func DeleteById(id int) (sql string, args []any) {
	sql = `delete from calc_plan where id = `
	args = append(args, id)
	return sql, args
}

func UpdatePlan(req types.PlanDTO) (sql string, args []any) {
	s := types.DTO2ST_PLAN(req)
	sql = `update calc_plan set 
	name = ?,
	comment = ?,
	start_time = ?,
	end_time = ?,
	details = ?
	where id = ?
	`
	args = append(args, s.Name, s.Comment, s.StartTime, s.EndTime, s.Details, s.Id)
	return sql, args
}
