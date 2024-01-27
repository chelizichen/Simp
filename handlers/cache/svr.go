package cache

import (
	h "Simp/handlers/http"
	"database/sql"
	"time"
)

const (
	DB_NAME = "simpcache_"
)

type CacheSvr struct {
	SVR      ICache
	InitTime string
	CTX      *h.SimpHttpServerCtx
}

type SimpCacheTableStruct struct {
	// 定义表的结构，字段应该和数据库表的字段对应
	ID          sql.NullInt32  `db:"id"`
	Key         sql.NullString `db:"name"`
	Value       sql.NullByte   `db:"value"`
	CreatetTime sql.NullString `db:"create_time"`
}

func InitizalCacheSvr(ctx *h.SimpHttpServerCtx) (cacheSvr *CacheSvr) {
	servant := NewMemCache()
	t := time.Now().Format(time.DateTime)
	cacheSvr.SVR = servant
	cacheSvr.InitTime = t
	cacheSvr.CTX = ctx
	return
}

type CacheWithExpired func(k string, v interface{}) error
