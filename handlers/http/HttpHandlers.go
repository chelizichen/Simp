package handlers

import (
	"Simp/config"
	"Simp/utils"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SimpHttpServerCtx struct {
	port   string
	name   string
	Engine *gin.Engine
	isMain bool
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
		name:   conf.Server.Name,
		port:   ":" + strconv.Itoa(conf.Server.Port),
		Engine: G,
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
	err := ctx.Engine.Run(ctx.port)
	if err != nil {
		return
	}
}
