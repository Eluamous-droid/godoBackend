package main

import (
	"context"
	"time"
	"todoBackend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@localhost:27017"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	collection := client.Database("todoApp").Collection("todoItems")
	item := models.NewTodoItem("work", "sut min dut2")
	collection.InsertOne(ctx, item)
	cursor, err := collection.Find(ctx, bson.D{})
	defer cursor.Close(ctx)
	var items []models.TodoItem
	for cursor.Next(ctx) {
		var item models.TodoItem
		if err := cursor.Decode(&item); err != nil {
			panic(err)
		}
		items = append(items, item)
	}
	for _, item := range items {
		println(item.Group + ": " + item.Task)
	}
}
