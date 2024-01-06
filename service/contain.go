package service

import (
	handlers "Simp/handlers/fass/simp"
	"github.com/gin-gonic/gin"
	"io"
)

var G *gin.Engine

func init() {
	G.POST("/invoke/fass", func(context *gin.Context) {
		servant := context.Request.Header.Get("Servant")
		body, err := io.ReadAll(context.Request.Body)
		if err != nil {
			context.JSON(500, gin.H{"error": "Failed to read request body"})
			return
		}
		handler := handlers.FassServants[servant]
		write, err := handler.Write(string(body))
		if err != nil {
			context.JSON(500, gin.H{"error": "Failed to write request body"})
			return
		}
		context.JSON(200, gin.H{"message": "Failed to read request body", "data": write})
	})
}
