package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
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

func GetPrompt(c *gin.Context) {
	gitlabProjectId, exists := c.GetQuery("gitlabProjectId")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No gitlab project id found"})
		return
	}
	var prompt Prompt
	err := GetPromptsCollection().FindOne(context.TODO(), bson.M{"gitlabProjectId": gitlabProjectId}).Decode(&prompt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Prompt decode error"})
		return
	}
	c.JSON(http.StatusOK, prompt)
}

func SetBasePrompt(c *gin.Context) {
	var prompt Prompt
	if err := c.ShouldBindJSON(&prompt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
		return
	}
	_, err := GetPromptsCollection().DeleteOne(context.TODO(), bson.M{"gitlabProjectId": prompt.GitlabProjectId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error", "error": err.Error()})
		return
	}
	_, err = GetPromptsCollection().InsertOne(context.TODO(), prompt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Insert prompt error", "error": err.Error()})
	}
}
