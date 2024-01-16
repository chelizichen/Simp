package storage

import "Simp/servers/CalcServer/types"

func SavePlan(req types.ST_Plan) (sql string, args []any) {
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
