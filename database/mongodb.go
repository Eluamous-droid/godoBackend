package database

import (
	"context"
	"godoBackend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var ctx = context.TODO()

func getCollection() {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		panic(err)
	}
	collection = client.Database("todoApp").Collection("todoItems")
}

func cleanUp() {
	if err := collection.Database().Client().Disconnect(ctx); err != nil {
		panic(err)
	}
}

func GetAllItems() []models.TodoItem {
	getCollection()

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@localhost:27017"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	collection := client.Database("todoApp").Collection("todoItems")
	collection.InsertOne(ctx, item)
}
