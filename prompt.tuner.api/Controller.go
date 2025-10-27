package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func Webhook(c *gin.Context) {
	log.Println("=== HEADERS ===")
	for key, values := range c.Request.Header {
		for _, value := range values {
			log.Printf("%s: %s\n", key, value)
		}
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	fmt.Println("=== BODY ===")
	fmt.Println(string(body))
	c.JSON(http.StatusOK, gin.H{})
}
