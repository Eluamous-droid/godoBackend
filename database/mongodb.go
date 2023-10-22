package database

import (
	"context"
	"godoBackend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GetAllItems(c *gin.Context) {
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
	c.IndentedJSON(http.StatusOK, items)

}

func GetAllItemsInGroup(c *gin.Context) {
	initClient()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("todoApp").Collection("todoItems")
	defer cleanUp()

	var request models.TodoItem
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	filter := bson.D{primitive.E{Key: "group", Value: request.Group}}
	cursor, err := collection.Find(ctx, filter)
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
	c.IndentedJSON(http.StatusOK, items)

}
func AddItem(c *gin.Context) {

	initClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("todoApp").Collection("todoItems")
	defer cleanUp()

	var request models.TodoItem
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	collection.InsertOne(ctx, request)
	c.Status(200)
}

func DeleteItem(c *gin.Context) {

	initClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("todoApp").Collection("todoItems")
	defer cleanUp()

	var request models.TodoItem
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	filter := bson.D{primitive.E{Key: "id", Value: request.Id}}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		panic(err)
	}
	c.Status(200)
}

func MarkItemComplete(c *gin.Context) {

	initClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("todoApp").Collection("todoItems")
	defer cleanUp()

	var request models.TodoItem
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	filter := bson.D{primitive.E{Key: "id", Value: request.Id}}
	update := bson.D{primitive.E{Key: "status", Value: models.Done}}
	collection.UpdateOne(ctx, filter, update)
	c.Status(200)
}
