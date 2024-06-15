package mongodb

import (
	"context"
	"fmt"
	"time"

	"alpha.com/configuration"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB() *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI(configuration.MONGO_URI).SetMaxPoolSize(4).
		SetMinPoolSize(2).
		SetMaxConnIdleTime(1 * time.Second)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err.Error())
	}

	PingMongoDB(client)

	return client
}

func DisconnectMongoDB(client *mongo.Client) {
	// Disconnect from MongoDB
	if err := client.Disconnect(context.TODO()); err != nil {
		fmt.Printf("MongoDB Shutdown Error: %v\n", err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func PingMongoDB(client *mongo.Client) {
	// Check the connection
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Connected to MongoDB!")
}
