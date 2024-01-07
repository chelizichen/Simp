package handlers

import (
	"Simp/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
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

	err := ctx.Engine.Run(ctx.port)
	if err != nil {
		return
	}
}
