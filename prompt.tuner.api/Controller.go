package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"net/http"
	"prompt.tuner.api/entity"
	"strings"
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
	var collection = GetAiMergeRequestCommentsCollection()
	var aiMergeRequestCommentRecord entity.AiMrComment
	mongoFindErr := collection.FindOne(
		context.TODO(),
		bson.M{
			"url": gitlabWebhookRequest.ObjectAttributes.AwardedOnUrl,
		},
	).Decode(&aiMergeRequestCommentRecord)
	if mongoFindErr != nil {
		if errors.Is(mongoFindErr, mongo.ErrNoDocuments) {
			newAiCommentRecord, newAiCommentRecordErr := GetNewAiCommentRecord(gitlabWebhookRequest)
			if newAiCommentRecordErr != nil {
				log.Println("[Webhook] create new ai comment mr record error " + newAiCommentRecordErr.Error())
				c.JSON(http.StatusOK, gin.H{})
				return
			}
			var result, err = collection.InsertOne(context.TODO(), newAiCommentRecord)
			if err != nil {
				log.Println(fmt.Sprintf("[Webhook] Database insert ai comment mr record error %s", err.Error()))
				c.JSON(http.StatusOK, gin.H{})
				return
			}
			log.Println("[Webhook] completed, insterted new ai comment mr record")
			c.JSON(http.StatusOK, result)
			return
		}
		log.Println("[Webhook] mongo find error " + mongoFindErr.Error())
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	var intLikeOrDislikeRepresentation int
	if IsWebhookRequestAddReaction(gitlabWebhookRequest) {
		intLikeOrDislikeRepresentation = 1
	} else {
		intLikeOrDislikeRepresentation = -1
	}
	var updateTransaction any
	isWebhookRequestLike, isWebhookRequestLikeErr := IsWebhookRequestLike(gitlabWebhookRequest)
	if isWebhookRequestLikeErr != nil {
		log.Println("[Webhook] not like or dislike reaction " + mongoFindErr.Error())
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	if isWebhookRequestLike {
		updateTransaction = bson.M{
			"$inc": bson.M{
				"likescount": intLikeOrDislikeRepresentation,
			},
		}
	} else {
		updateTransaction = bson.M{
			"$inc": bson.M{
				"dislikescount": intLikeOrDislikeRepresentation,
			},
		}
	}
	var updateResult, updateErr = collection.UpdateOne(
		context.TODO(),
		bson.M{
			"url": gitlabWebhookRequest.ObjectAttributes.AwardedOnUrl,
		},
		updateTransaction,
	)
	if updateErr != nil {
		log.Println("[Webhook] update error " + updateErr.Error())
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	c.JSON(http.StatusOK, updateResult)
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
	projectReactions, err := GetAiMrCommentsForGitlabProject(gitlabProjectId)
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
