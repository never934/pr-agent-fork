package main

import (
	"context"
	"fmt"
	"github.com/cohesion-org/deepseek-go"
	"prompt.tuner.api/entity"
)

func TuneBasePrompt(basePrompt string, reactions []entity.Reaction) (string, error) {
	// DEEPSEEK_API_KEY env
	client := deepseek.NewClient("")
	type ReactionRepresentation struct {
		typeRepresentation       string
		aiCommentsRepresentation string
	}
	var reactionsRepresentation []ReactionRepresentation
	for _, reaction := range reactions {
		var aiCommentsString string
		for _, aiCommentItem := range reaction.AiComments {
			aiCommentsString += aiCommentItem + "\n"
		}
		typeRepresentation := ""
		if reaction.Type == entity.NegativeReaction {
			typeRepresentation = "Негативная реакция"
		}
		if reaction.Type == entity.PositiveReaction {
			typeRepresentation = "Позитивная реакция"
		}
		var reactionRepresentation = ReactionRepresentation{
			typeRepresentation:       typeRepresentation,
			aiCommentsRepresentation: aiCommentsString,
		}
		reactionsRepresentation = append(reactionsRepresentation, reactionRepresentation)
	}
	var reactionsString string
	for _, reactionRepresentation := range reactionsRepresentation {
		reactionsString += reactionRepresentation.typeRepresentation + "\n"
	}
	message := fmt.Sprintf(
		"Есть базовый промт %s и есть реакции на него: %s . Дай улучшенный промпт на основе реакций",
		basePrompt,
		reactionsString,
	)
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
