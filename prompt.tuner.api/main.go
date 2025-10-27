package main

import "github.com/gin-gonic/gin"

func main() {
	ginInstance := gin.Default()
	ginInstance.POST("/webhook", Webhook)
	_ = ginInstance.Run(":4888")
}
