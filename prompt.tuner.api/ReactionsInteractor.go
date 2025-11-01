package main

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"prompt.tuner.api/entity"
)

func GetAiMrCommentsForGitlabProject(gitlabProjectId string) ([]entity.AiMrComment, error) {
	var collection = GetAiMergeRequestCommentsCollection()
	var filter = bson.M{"gitlabprojectid": gitlabProjectId}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createddate", Value: -1}})
	findOptions.SetLimit(100)
	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var aiMrComments []entity.AiMrComment
	for cursor.Next(context.Background()) {
		var reaction entity.AiMrComment
		err := cursor.Decode(&reaction)
		if err != nil {
			return nil, err
		}
		aiMrComments = append(aiMrComments, reaction)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return aiMrComments, nil
}
