package mongodb

import (
	"context"
	"fmt"

	"alpha.com/configuration"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB() *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI(configuration.MONGO_URI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err.Error())
	}

	PingMongoDB(client)

	return client
}

func PingMongoDB(client *mongo.Client) {

	// Check the connection
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Connected to MongoDB!")

	// Close the connection once no longer needed
	// defer func() {
	// 	if err := client.Disconnect(context.TODO()); err != nil {
	// 		fmt.Println(err.Error())
	// 	}
	// }()
}
