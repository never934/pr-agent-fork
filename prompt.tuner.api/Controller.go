package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"net/http"
	"prompt.tuner.api/entity"
	"strconv"
	"strings"
	"time"
)

func Webhook(c *gin.Context) {
	gitlabEventHeader := c.Request.Header.Get("X-Gitlab-Event")
	if gitlabEventHeader == "" {
		log.Println("[Webhook] X-Gitlab-Event header not found")
		c.JSON(http.StatusBadRequest, gin.H{"message": "X-Gitlab-Event header not found"})
		return
	}
	if gitlabEventHeader != "Emoji Hook" {
		log.Println("[Webhook] X-Gitlab-Event header not emoji hook")
		c.JSON(http.StatusBadRequest, gin.H{"message": "X-Gitlab-Event header not emoji hook"})
		return
	}
	var gitlabWebhookRequest entity.GitlabWebhookRequest
	if err := c.ShouldBindJSON(&gitlabWebhookRequest); err != nil {
		log.Println(fmt.Sprintf("[Webhook] Invalid request body %s", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}
	if strings.Contains(gitlabWebhookRequest.User.Username, "bot") {
		log.Println("[Webhook] reaction from bot")
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	if !strings.Contains(gitlabWebhookRequest.Note.Description, "AI Ревью") {
		log.Println("[Webhook] reaction not for ai review")
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	var collection = GetReactionsCollection()
	if gitlabWebhookRequest.ObjectAttributes.Action == "revoke" {
		collection.DeleteOne(context.TODO(), bson.M{"reactionurl": gitlabWebhookRequest.AwardedOnUrl})
		log.Println("[Webhook] reaction deleted")
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	reactionType := ""
	if gitlabWebhookRequest.ObjectAttributes.Name == "thumbsup" {
		reactionType = entity.PositiveReaction
	}
	if gitlabWebhookRequest.ObjectAttributes.Name == "thumbsdown" {
		reactionType = entity.NegativeReaction
	}
	reaction := entity.Reaction{
		Type:            reactionType,
		AiComment:       gitlabWebhookRequest.Note.Description,
		CreateDate:      time.Now().Format(time.RFC1123),
		GitlabProjectId: strconv.Itoa(gitlabWebhookRequest.ProjectId),
		ReactionUrl:     gitlabWebhookRequest.AwardedOnUrl,
	}
	var result, err = collection.InsertOne(context.Background(), reaction)
	if err != nil {
		log.Println(fmt.Sprintf("[Webhook] Database insert reaction error %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database insert reaction error"})
		return
	}
	log.Println("[Webhook] completed")
	c.JSON(http.StatusOK, result)
}

func GetPrompt(c *gin.Context) {
	gitlabProjectId, exists := c.GetQuery("gitlabProjectId")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No gitlab project id found"})
		return
	}
	forceRegenerate, exists := c.GetQuery("forceRegenerate")
	collection := GetPromptsCollection()
	filter := bson.M{"gitlabprojectid": gitlabProjectId}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Prompt not found"})
		return
	}
	var prompt entity.Prompt
	err = collection.FindOne(context.TODO(), filter).Decode(&prompt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Prompt decode error"})
		return
	}
	var tunedPromptsCache = GetTunedPromptsCache()
	if tunedPromptsCache.Exists(gitlabProjectId) && forceRegenerate != "true" {
		c.JSON(
			http.StatusOK,
			gin.H{
				"basePrompt":  prompt.Text,
				"tunedPrompt": tunedPromptsCache.Get(gitlabProjectId),
			},
		)
		return
	}
	projectReactions, err := GetReactionsForGitlabProject(gitlabProjectId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Project reactions decode error"})
		return
	}
	tunedPrompt, err := TuneBasePrompt(prompt.Text, projectReactions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Tune prompt error"})
		return
	}
	tunedPromptsCache.Add(gitlabProjectId, tunedPrompt)
	c.JSON(
		http.StatusOK,
		gin.H{
			"basePrompt":  prompt.Text,
			"tunedPrompt": tunedPrompt,
		},
	)
}

func SetBasePrompt(c *gin.Context) {
	var prompt entity.Prompt
	if err := c.ShouldBindJSON(&prompt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}
	collection := GetPromptsCollection()
	filter := bson.M{"gitlabprojectid": prompt.GitlabProjectId}
	_, err := collection.ReplaceOne(
		context.TODO(),
		filter,
		prompt,
		options.Replace().SetUpsert(true),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{"message": "Prompt updated successfully"},
	)
}
