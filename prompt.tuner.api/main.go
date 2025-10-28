package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	ginInstance := gin.Default()
	ginInstance.POST("/webhook", Webhook)
	ginInstance.POST("/setBasePrompt", SetBasePrompt)
	ginInstance.GET("/getPrompt", GetPrompt)
	_ = ginInstance.Run(":" + os.Getenv("APP_PORT"))
}
