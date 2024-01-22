package handlers

import (
	"Simp/config"
	"Simp/utils"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron"
)

type SimpHttpServerCtx struct {
	port        string
	name        string
	Engine      *gin.Engine
	Storage     *sqlx.DB
	isMain      bool
	StoragePath string
	StaticPath  string
}

func Resp(code int, message string, data interface{}) *gin.H {
	return &gin.H{
		"Code":    code,
		"Message": message,
		"Data":    data,
	}
}

func (c *SimpHttpServerCtx) Use(callback func(engine *SimpHttpServerCtx, pre string)) {
	s := c.name
	pre := strings.ToLower(s)
	callback(c, pre)
}

// 主控服务需要做日志系统与监控
func (c *SimpHttpServerCtx) DefineMain() {
	// go func() {
	utils.AutoSetLogWriter()
	go func() {
		c := cron.New()

		// 4小时执行一次，更换日志文件指定目录
		spec := "* * 4 * * *"

		// 添加定时任务
		err := c.AddFunc(spec, func() {
			utils.AutoSetLogWriter()
		})
		if err != nil {
			fmt.Println("AddFuncErr", err)
		}
		// 启动Cron调度器
		go c.Start()
	}()

	// 关闭文件
	// logWriter.file.Close()
	// }()
	// os.
	c.isMain = true
}
func (c *SimpHttpServerCtx) Post(realPath string, handle gin.HandlerFunc) {
	c.Engine.POST(realPath, handle)
}

func (c *SimpHttpServerCtx) Get(realPath string, handle gin.HandlerFunc) {
	c.Engine.GET(realPath, handle)
}

func (c *SimpHttpServerCtx) Static(realPath string) {
	s := c.name
	pre := strings.ToLower(s)
	f := utils.Join(pre)
	target := f(realPath)

	wd, _ := os.Getwd()
	SIMP_PRODUCTION := os.Getenv("SIMP_PRODUCTION")
	var staticPath string
	if SIMP_PRODUCTION == "Yes" {
		SIMP_SERVER_PATH := os.Getenv("SIMP_SERVER_PATH")
		staticPath = filepath.Join(SIMP_SERVER_PATH, c.StaticPath)
	} else {
		staticPath = filepath.Join(wd, c.StaticPath)
	}
	c.Engine.Static(target, staticPath)
}

func NewSimpHttpCtx(path string) (ctx *SimpHttpServerCtx) {
	conf, err := config.NewConfig(path)
	if err != nil {
		fmt.Println("NewConfig Error:", err.Error())
		panic(err.Error())
	}
	fmt.Println("SIMP_PRODUCTION Conf", conf)
	var G *gin.Engine
	G = gin.Default()
	if err != nil {
		fmt.Println("get Config Error :", err.Error())
	}
	// database, err := sqlx.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/stockpool")
	// 主控不提供数据库存储服务，存储服务由子服务提供
	// 子服务生产时数据库链接不上将会panic
	SIMP_PRODUCTION := os.Getenv("SIMP_PRODUCTION")
	var ST *sqlx.DB
	if !conf.Server.Main {
		fmt.Println("conf.Server.Storage ", conf.Server.Storage)
		database, err := sqlx.Open("mysql", conf.Server.Storage)
		if err != nil && SIMP_PRODUCTION == "Yes" {
			panic("init db error" + err.Error())
		} else if err != nil {
			fmt.Println("init db error", err.Error())
		}
		ST = database
		fmt.Println("ctx.Storage", database, database == nil)
		err = database.Ping()
		if err != nil {
			fmt.Println("Error! database ping ", err.Error())
		}
	}

	ctx = &SimpHttpServerCtx{
		name:        conf.Server.Name,
		port:        ":" + strconv.Itoa(conf.Server.Port),
		Engine:      G,
		StoragePath: conf.Server.Storage,
		StaticPath:  conf.Server.StaticPath,
		Storage:     ST,
	}
	return ctx
}

func NewSimpHttpServer(ctx *SimpHttpServerCtx) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, addr := range addrs {
		// Check if the address is not a loopback address
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("IPv4 Address:", ipnet.IP.String())
			} else {
				fmt.Println("IPv6 Address:", ipnet.IP.String())
			}
		}
	}
	SIMP_PRODUCTION := os.Getenv("SIMP_PRODUCTION")
	fmt.Println("SIMP_PRODUCTION", SIMP_PRODUCTION)
	// 子服务生产时需要提供对应的API路由
	if SIMP_PRODUCTION == "Yes" {
		fmt.Println("CreateAPIFile |", ctx.name)
		utils.CreateAPIFile(ctx.Engine, ctx.name)
	}

	err = ctx.Engine.Run(ctx.port)
	if err != nil {
		return
	}
}
