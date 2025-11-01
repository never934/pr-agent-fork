package main

import (
	"context"
	"fmt"
	"github.com/cohesion-org/deepseek-go"
	"log"
	"os"
	"prompt.tuner.api/entity"
)

func TuneBasePrompt(basePrompt string, aiMrComments []entity.AiMrComment) (string, error) {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	type AiCommentRepresentation struct {
		likesRepresentation           string
		dislikesRepresentation        string
		aiCommentPointsRepresentation string
	}
	var aiCommentsRepresentation []AiCommentRepresentation
	for _, aiMrComment := range aiMrComments {
		var aiCommentPointsRepresentation string
		for _, aiCommentPoint := range aiMrComment.CommentPoints {
			aiCommentPointsRepresentation += "\n" + aiCommentPoint
		}
		var aiCommentRepresentation = AiCommentRepresentation{
			likesRepresentation:           fmt.Sprintf("Лайков %d", aiMrComment.LikesCount),
			dislikesRepresentation:        fmt.Sprintf("Дизлайков %d", aiMrComment.DislikesCount),
			aiCommentPointsRepresentation: aiCommentPointsRepresentation,
		}
		aiCommentsRepresentation = append(aiCommentsRepresentation, aiCommentRepresentation)
	}
	var aiCommentsString string
	for _, aiCommentRepresentation := range aiCommentsRepresentation {
		aiCommentsString +=
			"\n\n" + aiCommentRepresentation.aiCommentPointsRepresentation + "\n" +
				aiCommentRepresentation.likesRepresentation + "\n" +
				aiCommentRepresentation.dislikesRepresentation + "\n\n"
	}
	message := fmt.Sprintf(
		"Есть базовый промт %s и есть реакции на него: %s Дай улучшенный промпт на основе реакций",
		basePrompt,
		aiCommentsString,
	)
	log.Println(fmt.Sprintf("Process AI with message %s", message))
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekReasoner,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleSystem, Content: "Выдавай в ответ только модифицированный промпт и ничего более"},
			{Role: deepseek.ChatMessageRoleUser, Content: message},
		},
	}
	response, err := client.CreateChatCompletion(context.TODO(), request)
	if err != nil {
		return "", err
	}
	return response.Choices[0].Message.Content, nil
}
