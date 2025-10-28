package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}
	ginInstance := gin.Default()
	ginInstance.POST("/webhook", Webhook)
	ginInstance.POST("/setBasePrompt", SetBasePrompt)
	ginInstance.GET("/getPrompt", GetPrompt)
	_ = ginInstance.Run(":4888")
}
