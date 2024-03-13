package http

import (
	"Simp/src/cache"
	"Simp/src/config"
	"Simp/src/utils"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron"
)

type SimpHttpServerCtx struct {
	Port        string
	Name        string
	Engine      *gin.Engine
	Storage     *sqlx.DB
	isMain      bool
	StoragePath string
	StaticPath  string
	CacheSvr    cache.ICache
	Proxy       *[]struct {
		Server struct {
			Type string `yaml:"type"`
			Name string `yaml:"name"`
			Port string `yaml:"port"`
		} `yaml:"server"`
	} `yaml:"proxy"`
}

func Resp(code int, message string, data interface{}) *gin.H {
	return &gin.H{
		"Code":    code,
		"Message": message,
		"Data":    data,
	}
}

func (c *SimpHttpServerCtx) Use(callback func(engine *SimpHttpServerCtx, pre string)) {
	s := c.Name
	pre := strings.ToLower(s)
	callback(c, pre)
}

// 主控服务需要做日志系统与监控
func (c *SimpHttpServerCtx) DefineMain() {
	cwd, err := os.Getwd()
	if err != nil {
		Err_Message := "Error To Get Wd" + err.Error()
		fmt.Println(Err_Message)
		panic(Err_Message)
	}
	mainPath := filepath.Join(cwd, "static/main")
	b := utils.IsExist(filepath.Join(cwd, "static/main"))
	if !b {
		err = os.Mkdir(mainPath, os.ModePerm)
		if err != nil {
			Err_Message := "Error To Mkdir" + err.Error()
			fmt.Println(Err_Message)
			panic(Err_Message)
		}
	}
	utils.AutoSetLogWriter()
	go func() {
		c := cron.New()

		// 4小时执行一次，更换日志文件指定目录
		spec := "* * */4 * * *"

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

func (c *SimpHttpServerCtx) UseSPA(path string, root string) {
	wd, _ := os.Getwd()
	SIMP_PRODUCTION := os.Getenv("SIMP_PRODUCTION")
	s := c.Name
	pre := strings.ToLower(s)
	f := utils.Join(pre)

	// 设置缓存
	setCacheHeaders := func(ctx *gin.Context, fileInfo os.FileInfo) {
		ctx.Header("Cache-Control", "public, max-age=2592000")
		expires := time.Now().Add(time.Hour * 24 * 30)
		ctx.Header("Expires", expires.Format(time.RFC1123))
		lastModified := fileInfo.ModTime()
		ctx.Header("Last-Modified", lastModified.Format(time.RFC1123))
	}

	c.Engine.GET(f(path)+"/*path", func(ctx *gin.Context) {
		requestPath := ctx.Param("path")
		var webRoot, targetPath string

		if SIMP_PRODUCTION == "Yes" {
			SIMP_SERVER_PATH := os.Getenv("SIMP_SERVER_PATH")
			webRoot = filepath.Join(SIMP_SERVER_PATH, root)
			targetPath = filepath.Join(SIMP_SERVER_PATH, root, requestPath)
		} else {
			webRoot = filepath.Join(wd, root)
			targetPath = filepath.Join(wd, root, requestPath)
		}

		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			targetPath = filepath.Join(webRoot, "index.html")
		} else if fileInfo, err := os.Stat(targetPath); err == nil {
			if strings.HasSuffix(targetPath, ".js") || strings.HasSuffix(targetPath, ".css") {
				setCacheHeaders(ctx, fileInfo)
			}
		}

		ctx.File(targetPath)
	})
}

func (c *SimpHttpServerCtx) Static(realPath string, args ...string) {
	s := c.Name
	pre := strings.ToLower(s)
	f := utils.Join(pre)
	target := f(realPath)

	wd, _ := os.Getwd()
	SIMP_PRODUCTION := os.Getenv("SIMP_PRODUCTION")
	var staticPath string
	if SIMP_PRODUCTION == "Yes" {
		SIMP_SERVER_PATH := os.Getenv("SIMP_SERVER_PATH")
		if len(args) > 0 {
			otherPath := filepath.Join(args...)
			staticPath = filepath.Join(SIMP_SERVER_PATH, otherPath)
		} else {
			staticPath = filepath.Join(SIMP_SERVER_PATH, c.StaticPath)
		}
	} else {
		if len(args) > 0 {
			otherPath := filepath.Join(args...)
			staticPath = filepath.Join(wd, otherPath)
		} else {
			staticPath = filepath.Join(wd, c.StaticPath)
		}
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
		Name:        conf.Server.Name,
		Port:        ":" + strconv.Itoa(conf.Server.Port),
		Engine:      G,
		StoragePath: conf.Server.Storage,
		StaticPath:  conf.Server.StaticPath,
		Storage:     ST,
		Proxy:       &conf.Server.Proxy,
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
		fmt.Println("CreateAPIFile |", ctx.Name)
		utils.CreateAPIFile(ctx.Engine, ctx.Name)
	}

	err = ctx.Engine.Run(ctx.Port)
	if err != nil {
		return
	}
}
