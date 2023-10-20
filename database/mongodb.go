package database

import (
	"context"
	"godoBackend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func initClient() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		panic(err)
	}
}

func cleanUp() {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func GetAllItems() []models.TodoItem {
	initClient()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("todoApp").Collection("todoItems")
	defer cleanUp()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)
	var items []models.TodoItem
	for cursor.Next(ctx) {
		var item models.TodoItem
		if err := cursor.Decode(&item); err != nil {
			panic(err)
		}
		items = append(items, item)
	}
	return items

}

func AddItem(item models.TodoItem) {

	initClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("todoApp").Collection("todoItems")
	defer cleanUp()

	collection.InsertOne(ctx, item)
}
