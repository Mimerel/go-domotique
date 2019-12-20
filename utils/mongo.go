package utils

import (
	"context"
	"fmt"
	"go-domotique/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func InsertMany(c *models.Configuration, collectionName string, data []interface{}) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:admin@192.168.222.240:27017"))
	if err != nil {
		fmt.Printf( "error connecting %v", err)
		return
	}
	database := client.Database("domotique")
	collection := database.Collection(collectionName)

	_, err = collection.InsertMany(ctx, data)
	if err != nil {
		fmt.Printf("error inserting %v", err)
	}
}
