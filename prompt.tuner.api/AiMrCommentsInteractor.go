package main

import (
	"os"
	"prompt.tuner.api/entity"
	"strconv"
	"time"
)

func GetNewAiCommentRecord(webhookRequest entity.GitlabWebhookRequest) (entity.AiMrComment, error) {
	likesCount := 0
	dislikesCount := 0
	incValue := 0
	if IsWebhookRequestAddReaction(webhookRequest) {
		incValue = 1
	} else {
		incValue = -1
	}
	isWebhookRequestLike, err := IsWebhookRequestLike(webhookRequest)
	if err != nil {
		return entity.AiMrComment{}, err
	}
	if isWebhookRequestLike {
		likesCount += incValue
	} else {
		dislikesCount += incValue
	}
	return entity.AiMrComment{
		CommentPoints:   ParseCommentPoints(webhookRequest.Note.Description),
		CreateDate:      time.Now(),
		GitlabProjectId: strconv.Itoa(webhookRequest.ProjectId),
		Url:             webhookRequest.ObjectAttributes.AwardedOnUrl,
		LikesCount:      likesCount,
		DislikesCount:   dislikesCount,
	}, nil
}

func IsWebhookRequestLike(webhookRequest entity.GitlabWebhookRequest) (bool, error) {
	if webhookRequest.ObjectAttributes.Name == "thumbsup" {
		return true, nil
	}
	if webhookRequest.ObjectAttributes.Name == "thumbsdown" {
		return false, nil
	}
	return false, os.ErrInvalid
}

func IsWebhookRequestAddReaction(webhookRequest entity.GitlabWebhookRequest) bool {
	if webhookRequest.EventType == "revoke" {
		return false
	}
	return true
}
