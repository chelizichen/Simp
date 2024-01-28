// Simp Cache Servant
// 1. expire delete 将会从cache 和 mysql中同步更改状态 0 expire 1 delete 2 online
// 2. 热key在保存至数据库的同时 也会长期存在于内存中
// 3. cache长时间未访问，超过24h？12h？的会从内存中移除。下次访问时会从mysql中取出。
// 整个 servant 相当于 LRU + MySql的形式
package svr

import (
	"Simp/handlers/cache"
	h "Simp/handlers/http"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron"
)

type SimpCacheHandleFunc func(ctx *h.SimpHttpServerCtx) cache.ExpiredCallback

// 检查表是否存在
func tableExists(db *sqlx.DB, tableName string) bool {
	query := fmt.Sprintf("SHOW TABLES LIKE '%s'", tableName)
	var result string
	err := db.Get(&result, query)
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		fmt.Println("err", err.Error())
	}
	return result == tableName
}

// 创建表
func createTable(db *sqlx.DB, tableName string) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INT PRIMARY KEY AUTO_INCREMENT,
		status INT PRIMARY,
		name VARCHAR(255) NOT NULL,
		value LONGBLOB,
		create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`, tableName)
	_, err := db.Exec(query)
	return err
}

// 插入数据
func insertData(db *sqlx.DB, tableName, key string, value []byte) error {
	query := fmt.Sprintf("INSERT INTO %s (name, value) VALUES (?, ?)", tableName)
	_, err := db.Exec(query, key, value)
	return err
}

// 超过20万条就新增一个新表
func isRowCountsTooBig(db *sqlx.DB, tableName string) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := db.Get(&count, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func InitizalCacheSvr(ctx *h.SimpHttpServerCtx) (cacheSvr *cache.CacheSvr) {
	sch := SimpCacheHookImpl(ctx)
	servant := cache.NewMemCache(cache.WithDeleteCallback(sch.Delete), cache.WithExpiredCallback(sch.Exipred))
	t := time.Now().Format(time.DateTime)
	fmt.Println("svr point ", servant, &servant)
	fmt.Println("cacheSvr", cacheSvr, &cacheSvr)
	cacheSvr = &cache.CacheSvr{
		SVR:      servant,
		InitTime: t,
		CTX:      ctx,
	}
	return
}

func SimpCacheHookImpl(ctx *h.SimpHttpServerCtx) *cache.SimpCacheHook {
	var idx int = 1
	var latest_cache_table *string
	for {
		tableName := cache.DB_NAME + ctx.Name + string(rune(idx))
		if !tableExists(ctx.Storage, tableName) {
			// 如果表不存在，则创建表
			err := createTable(ctx.Storage, tableName)
			if err != nil {
				fmt.Println("Error To CreateTable", err.Error())
			}
			*latest_cache_table = tableName
			break
		} else {
			count, err := isRowCountsTooBig(ctx.Storage, tableName)
			if err != nil {
				fmt.Println("Error To CreateTable", err.Error())
			}
			if count > 200000 {
				idx = idx + 1
			} else {
				// 没大于20万条 则用该表
				*latest_cache_table = tableName
				break
			}
		}
	}

	go func() {
		c := cron.New()

		// 4小时执行一次，更换日志文件指定目录
		spec := "* * 4 * * *"

		// 添加定时任务
		err := c.AddFunc(spec, func() {
			count, err := isRowCountsTooBig(ctx.Storage, *latest_cache_table)
			if err != nil {
				fmt.Println("Error To CreateTable", err.Error())
			}
			if count > 200000 {
				idx = idx + 1
			}
			*latest_cache_table = cache.DB_NAME + ctx.Name + string(rune(idx))
			err = createTable(ctx.Storage, *latest_cache_table)
			if err != nil {
				fmt.Println("Error To CreateTable", err.Error())
			}
		})
		if err != nil {
			fmt.Println("AddFuncErr", err)
		}
		// 启动Cron调度器
		go c.Start()
	}()
	return &cache.SimpCacheHook{
		Exipred: func(k string, v interface{}) error {
			bV, err := json.Marshal(v)
			if err != nil {
				fmt.Println("Error Marshal", err.Error())
				return nil
			}
			insertData(ctx.Storage, *latest_cache_table, k, bV)
			return nil
		},
		Delete: func(k string, v interface{}) error {
			return nil
		},
	}
}
