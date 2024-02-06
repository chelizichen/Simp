// Simp Cache Servant
// 1. expire delete 将会从cache 和 mysql中同步更改状态 0 expire 1 delete 2 online
// 2. 热key在保存至数据库的同时 也会长期存在于内存中
// 3. cache长时间未访问，超过24h？12h？的会从内存中移除。下次访问时会从mysql中取出。
// 整个 servant 相当于 LRU + MySql的形式
// 	usage
//	func main() {
//		ctx := h.NewSimpHttpCtx("simp.yaml")
//		svr.InitizalCacheSvr(ctx)
//		h.NewSimpHttpServer(ctx)
//	}

package http

import (
	"Simp/src/cache"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron"
)

type SimpCacheHandleFunc func(ctx *SimpHttpServerCtx) cache.ExpiredCallback

const (
	PrepareError = "PrepareError "
)

// 检查表是否存在
func TableExists(db *sqlx.DB, tableName string) bool {
	query := fmt.Sprintf("SHOW TABLES LIKE '%s'", tableName)
	var result string
	err := db.Get(&result, query)
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		fmt.Println("err", err.Error())
	}
	return result == tableName
}

func GetCache(db *sqlx.DB, k string) (interface{}, bool) {
	var sci cache.DBCacheKey // Initialize the pointer
	query := "select t,k from simp_caches_set  where k = ?"
	err := db.Get(&sci, query, k)
	fmt.Println("")
	if err != nil {
		fmt.Println("Error to Get ", err.Error())
	}
	fmt.Println("sci ", sci.Key, sci.Table)
	before := "select v from [t] st where st.k = ? and st.s = ？ limit 1"
	queryString := strings.Replace(before, "[t]", sci.Table, 0)
	var row cache.SimpCacheItem // Initialize the pointer
	err = db.Get(&row, queryString, k, cache.ITEM_STATUS_DEFAULT)
	fmt.Println("Value ", row)
	if err != nil {
		return row.Value, true
	}
	return nil, false
}

// 创建集合表
func CreateCachesSet(db *sqlx.DB) error {
	query := " 	CREATE TABLE IF NOT EXISTS simp_caches_set ( " +
		"`id` 	INT NOT NULL  PRIMARY KEY AUTO_INCREMENT, \n" +
		"`k` 	VARCHAR(255) NOT NULL, \n" +
		"`t` 	VARCHAR(255) NOT NULL) \n"
	fmt.Println("prepare sql \n", query)
	sql, err := db.Prepare(query)
	if err != nil {
		fmt.Println("prepare error", err.Error())
		return err
	}
	_, err = sql.Exec()
	return err
}

// 创建表
func CreateTable(db *sqlx.DB, tableName string) error {
	t := strings.ToLower(tableName)
	query := " CREATE TABLE IF NOT EXISTS `" + t + "` ( " +
		"`id` INT NOT NULL  PRIMARY KEY AUTO_INCREMENT, \n" +
		"`s` INT NOT NULL, \n" +
		"`k` VARCHAR(255) NOT NULL, \n" +
		"`v` LONGBLOB NOT NULL, \n" +
		"`t` VARCHAR(255) NOT NULL) \n"
	fmt.Println("prepare sql \n", query)
	sql, err := db.Prepare(query)
	if err != nil {
		fmt.Println("prepare error", err.Error())
		return err
	}
	_, err = sql.Exec()
	return err
}

// 插入数据
func InsertData(db *sqlx.DB, tableName, key string, value []byte, status int) error {
	tableName = strings.ToLower(tableName)
	now := time.Now().Format(time.DateTime)
	query := "INSERT INTO `" + tableName + "` (k, v, s, t) VALUES (?, ?, ?, ?)"
	fmt.Println("query", query, key, value, tableName)
	fmt.Println("status ", status, " statusDefault ", cache.ITEM_STATUS_DEFAULT)
	if status == cache.ITEM_STATUS_DEFAULT {
		cache.InsertKeySet(db, key, tableName)
	}
	_, err := db.Exec(query, key, value, status, now)
	return err
}

// 超过20万条就新增一个新表
func IsRowCountsTooBig(db *sqlx.DB, tableName string) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := db.Get(&count, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func InitizalCacheSvr(ctx *SimpHttpServerCtx) {
	sch := SimpCacheHookImpl(ctx)
	servant := cache.NewMemCache(
		cache.WithDeleteCallback(sch.Delete),
		cache.WithExpiredCallback(sch.Exipred),
		cache.WithDefaultCallback(sch.Default),
		cache.WithGetWhenExipred(sch.GetWhenExpired),
	)
	fmt.Println("svr point ", servant, &servant)
	ctx.CacheSvr = servant
}

func SimpCacheHookImpl(ctx *SimpHttpServerCtx) *cache.SimpCacheHook {
	err := CreateCachesSet(ctx.Storage)
	if err != nil {
		ES := fmt.Sprintf("Error to create caches set %s", err.Error())
		panic(ES)
	}
	var idx int = 1
	var latest_cache_table string = ""
	for {
		tableName := cache.DB_NAME + ctx.Name + string(rune(idx))
		if !TableExists(ctx.Storage, tableName) {
			// 如果表不存在，则创建表
			err := CreateTable(ctx.Storage, tableName)
			if err != nil {
				ES := fmt.Sprintf("Error To CreateTable not exist %s", err.Error())
				panic(ES)
			}
			latest_cache_table = tableName
			break
		} else {
			count, err := IsRowCountsTooBig(ctx.Storage, tableName)
			if err != nil {
				ES := fmt.Sprintf("Error To CreateTable too big %s", err.Error())
				panic(ES)
			}
			if count > 200000 {
				idx = idx + 1
			} else {
				// 没大于20万条 则用该表
				latest_cache_table = tableName
				break
			}
		}
	}
	fmt.Println("latest_cache_table", latest_cache_table)
	go func() {
		c := cron.New()

		// 4小时执行一次，更换日志文件指定目录
		spec := "* * 4 * * *"

		// 添加定时任务
		err := c.AddFunc(spec, func() {
			count, err := IsRowCountsTooBig(ctx.Storage, latest_cache_table)
			if err != nil {
				fmt.Println("Error To CreateTable", err.Error())
			}
			if count > 200000 {
				idx = idx + 1
			}
			latest_cache_table = cache.DB_NAME + ctx.Name + string(rune(idx))
			err = CreateTable(ctx.Storage, latest_cache_table)
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
			InsertData(ctx.Storage, latest_cache_table, k, bV, cache.ITEM_STATUS_EXPIRE)
			return nil
		},
		Delete: func(k string, v interface{}) error {
			fmt.Println("执行删除callback")
			bV, err := json.Marshal(v)
			if err != nil {
				fmt.Println("Error Marshal", err.Error())
				return nil
			}
			err = InsertData(ctx.Storage, latest_cache_table, k, bV, cache.ITEM_STATUS_DELETE)
			if err != nil {
				fmt.Println("insert delete data error", err.Error())
				return err
			}
			return nil
		},
		Default: func(k string, v interface{}) error {
			bV, err := json.Marshal(v)
			if err != nil {
				fmt.Println("Error Marshal", err.Error())
				return nil
			}
			InsertData(ctx.Storage, latest_cache_table, k, bV, cache.ITEM_STATUS_DEFAULT)
			return nil
		},
		GetWhenExpired: func(k string) (value interface{}, exist bool) {
			return GetCache(ctx.Storage, k)
		},
	}
}
