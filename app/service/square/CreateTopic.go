package scoredatabase

import (
	"context"
	"hr/configs/models/square"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateTopic(userId string, title string, content string) (interface{}, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())
	database := client.Database("yourDatabase")
	collectionName := "yourCollection"

	topic := square.Topic{
		Title:    title,
		Content:  content,
		AutherId: userId,
	}
	collection := database.Collection(collectionName)
	result, err := collection.InsertOne(context.TODO(), topic)
	if err != nil {
		return "", err
	}

	return result.InsertedID, nil
}
