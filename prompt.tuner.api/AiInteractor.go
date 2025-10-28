package main

import (
	"context"
	"fmt"
	"github.com/cohesion-org/deepseek-go"
	"log"
	"os"
	"prompt.tuner.api/entity"
)

func TuneBasePrompt(basePrompt string, reactions []entity.Reaction) (string, error) {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	type ReactionRepresentation struct {
		typeRepresentation       string
		aiCommentsRepresentation string
	}
	var reactionsRepresentation []ReactionRepresentation
	for _, reaction := range reactions {
		typeRepresentation := ""
		if reaction.Type == entity.NegativeReaction {
			typeRepresentation = "Негативная реакция"
		}
		if reaction.Type == entity.PositiveReaction {
			typeRepresentation = "Позитивная реакция"
		}
		var reactionRepresentation = ReactionRepresentation{
			typeRepresentation:       typeRepresentation,
			aiCommentsRepresentation: reaction.AiComment,
		}
		reactionsRepresentation = append(reactionsRepresentation, reactionRepresentation)
	}
	var reactionsString string
	for _, reactionRepresentation := range reactionsRepresentation {
		reactionsString +=
			reactionRepresentation.typeRepresentation + " " + reactionRepresentation.aiCommentsRepresentation + "\n"
	}
	message := fmt.Sprintf(
		"Есть базовый промт %s и есть реакции на него: %s . Дай улучшенный промпт на основе реакций",
		basePrompt,
		reactionsString,
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
