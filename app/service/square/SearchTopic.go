package scoredatabase

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SearchTopic(UserId string, collection *mongo.Collection) ([]bson.M, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{"author": UserId},
		},
	}
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	// 获取所有匹配的文档
	var results []bson.M
	if err := cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
