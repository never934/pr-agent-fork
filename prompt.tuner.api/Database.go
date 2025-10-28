package main

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
)

var Database *mongo.Database

func GetDatabase() *mongo.Database {
	if Database == nil {
		var host = "10.30.23.8:27020"
		var login = "admin"
		var password = "adminpass"
		clientOpts := options.Client().SetHosts(
			[]string{host},
		).SetAuth(
			options.Credential{
				AuthSource:    "admin",
				AuthMechanism: "SCRAM-SHA-256",
				Username:      login,
				Password:      password,
			},
		)
		client, err := mongo.Connect(clientOpts)
		if err != nil {
			log.Println(err)
		}
		Database = client.Database("prompt_tuner_api")
		return Database
	} else {
		return Database
	}
}

func GetPromptsCollection() *mongo.Collection {
	return GetDatabase().Collection("prompts")
}

func GetReactionsCollection() *mongo.Collection {
	return GetDatabase().Collection("reactions")
}
