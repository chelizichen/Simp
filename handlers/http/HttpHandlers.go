package handlers

import (
	"Simp/config"
	"Simp/utils"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type SimpHttpServerCtx struct {
	port        string
	name        string
	Engine      *gin.Engine
	Storage     *sqlx.DB
	isMain      bool
	storagePath string
}

func Resp(code int, message string, data interface{}) *gin.H {
	return &gin.H{
		"Code":    code,
		"Message": message,
		"Data":    data,
	}
}

func (c *SimpHttpServerCtx) Use(callback func(engine *SimpHttpServerCtx)) {
	callback(c)
}

func (c *SimpHttpServerCtx) DefineMain() {
	c.isMain = true
}
func (c *SimpHttpServerCtx) Post(realPath string, handle gin.HandlerFunc) {
	c.Engine.POST(realPath, handle)
}

func (c *SimpHttpServerCtx) Get(realPath string, handle gin.HandlerFunc) {
	c.Engine.GET(realPath, handle)
}

func NewSimpHttpCtx(path string) (ctx *SimpHttpServerCtx) {
	conf, err := config.NewConfig(path)
	var G *gin.Engine
	G = gin.Default()
	if err != nil {
		fmt.Println("get Config Error :", err.Error())
	}
	ctx = &SimpHttpServerCtx{
		name:        conf.Server.Name,
		port:        ":" + strconv.Itoa(conf.Server.Port),
		Engine:      G,
		storagePath: conf.Server.Storage,
	}
	return ctx
}

func NewSimpHttpServer(ctx *SimpHttpServerCtx) {
	SIMP_PRODUCTION := os.Getenv("SIMP_PRODUCTION")
	fmt.Println("SIMP_PRODUCTION", SIMP_PRODUCTION)
	// 子服务生产时需要提供对应的API路由
	if SIMP_PRODUCTION == "Yes" {
		fmt.Println("CreateAPIFile |", ctx.name)
		utils.CreateAPIFile(ctx.Engine, ctx.name)
	}
	// database, err := sqlx.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/stockpool")
	// 主控不提供数据库存储服务，存储服务由子服务提供
	// 子服务生产时数据库链接不上将会panic
	if !ctx.isMain {
		fmt.Println("ctx.storagePath ", ctx.storagePath)
		database, err := sqlx.Open("mysql", ctx.storagePath)
		fmt.Println("database", err)
		if err != nil && SIMP_PRODUCTION == "Yes" {
			panic("init db error" + err.Error())
		} else if err != nil {
			fmt.Println("init db error", err.Error())
		}
		ctx.Storage = database
	}

	err := ctx.Engine.Run(ctx.port)
	if err != nil {
		return
	}
}
