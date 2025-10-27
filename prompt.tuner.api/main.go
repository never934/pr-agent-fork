package main

import "github.com/gin-gonic/gin"

func main() {
	ginInstance := gin.Default()
	ginInstance.POST("/webhook", Webhook)
	ginInstance.POST("/setBasePrompt", SetBasePrompt)
	ginInstance.GET("/getPrompt", GetPrompt)
	_ = ginInstance.Run(":4888")
}
