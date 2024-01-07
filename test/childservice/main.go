package main

import (
	h "Simp/handlers/http"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := h.NewSimpHttpCtx("test/childservice/simp.yaml")
	ctx.DefineMain()
	ctx.Get("/greet", func(context *gin.Context) {
		fmt.Println("hello greet")
		context.JSON(200, h.Resp(0, "ok", nil))
	})
	h.NewSimpHttpServer(ctx)
}
