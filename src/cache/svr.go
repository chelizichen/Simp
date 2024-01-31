package cache

import (
	h "Simp/src/http"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	DB_NAME = "simp_cache_"
)

type CacheSvr struct {
	SVR      ICache
	InitTime string
	CTX      *h.SimpHttpServerCtx
}

type SimpCacheTableStruct struct {
	// 定义表的结构，字段应该和数据库表的字段对应
	ID          sql.NullInt32  `db:"id"`
	Key         sql.NullString `db:"k"`
	Value       sql.NullByte   `db:"v"`
	CreatetTime sql.NullString `db:"t"`
	Status      sql.NullInt16  `db:"s"`
}

type SimpCacheHook struct {
	Exipred        ExpiredCallback
	Delete         ExpiredCallback
	Default        ExpiredCallback
	GetWhenExpired GetWhenExpiredFunc
}

func InsertKeySet(db *sqlx.DB, key string, table string) error {
	query := fmt.Sprintf("INSERT INTO simp_cache_keys (key, table) VALUES (?, ?)")
	_, err := db.Exec(query, key, table)
	return err
}
