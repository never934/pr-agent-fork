package main

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"prompt.tuner.api/entity"
)

func GetReactionsForGitlabProject(gitlabProjectId string) ([]entity.Reaction, error) {
	var collection = GetReactionsCollection()
	var filter = bson.M{"gitlabprojectid": gitlabProjectId}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createddate", Value: -1}})
	findOptions.SetLimit(100)
	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var reactions []entity.Reaction
	for cursor.Next(context.Background()) {
		var reaction entity.Reaction
		err := cursor.Decode(&reaction)
		if err != nil {
			return nil, err
		}
		reactions = append(reactions, reaction)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return reactions, nil
}

func GetReactionTypeFromGitlabReaction(gitlabReaction string) string {
	switch gitlabReaction {
	case "thumbsup":
		return entity.PositiveReaction
	case "thumbsdown":
		return entity.NegativeReaction
	default:
		return ""
	}
}
