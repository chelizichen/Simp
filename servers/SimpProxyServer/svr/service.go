package svr

import (
	h "Simp/src/http"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	BindJsonError = "BindJsonError"
)

type CacheItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func Service(engine *h.SimpHttpServerCtx, pre string) {
	EG := engine.Engine
	G := EG.Group(pre)
	G.POST("/set", func(ctx *gin.Context) {
		var i CacheItem
		err := ctx.BindJSON(&i)
		if err != nil {
			fmt.Println(BindJsonError + err.Error())
			ctx.JSON(http.StatusOK, h.Resp(-1, BindJsonError+err.Error(), nil))
			return
		}
	})
}
