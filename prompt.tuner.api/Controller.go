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
	var collection = GetPromptsCollection()
	var filter = bson.M{"gitlabProjectId": gitlabProjectId}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Prompt not found"})
		return
	}
	var prompt Prompt
	err = collection.FindOne(context.TODO(), filter).Decode(&prompt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Prompt decode error"})
		return
	}
	c.JSON(http.StatusOK, prompt)
}

func SetBasePrompt(c *gin.Context) {
	var prompt Prompt
	if err := c.ShouldBindJSON(&prompt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}
	var gitlabProjectId = prompt.GitlabProjectId
	var collection = GetPromptsCollection()
	var filter = bson.M{"gitlabProjectId": gitlabProjectId}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	_, err = collection.InsertOne(context.TODO(), prompt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Insert prompt error"})
	}
}
