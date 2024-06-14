package repository

import (
	"context"
	"fmt"

	"alpha.com/configuration"
	"alpha.com/internal/alpha.com/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserRepository interface {
	Get(ctx context.Context) ([]*domain.User, error)
	GetById(ctx context.Context, userId string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Upsert(ctx context.Context, user *domain.User) (string, error)
}

type userRepository struct {
	mongoClient *mongo.Client
}

func NewUserRepository(mongoClient *mongo.Client) IUserRepository {
	return &userRepository{

		mongoClient: mongoClient,
	}
}

func (r *userRepository) Get(ctx context.Context) ([]*domain.User, error) {
	collection := r.mongoClient.Database(configuration.MONGO_DB_NAME).Collection(configuration.MONGO_USERS_DB_NAME)

	var users []*domain.User
	cursor, err := collection.Find(context.TODO(), bson.D{})

	if err != nil {
		fmt.Printf("userRepository.Get ERROR : %s\n", err.Error())
		return make([]*domain.User, 0), err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var user *domain.User
		err := cursor.Decode(&user)
		if err != nil {
			fmt.Printf("userRepository.Get ERROR : %s\n", err.Error())
			return make([]*domain.User, 0), err
		}

		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		fmt.Printf("userRepository.Get ERROR : %s\n", err.Error())
	}

	if users == nil {
		fmt.Printf("userRepository.Get INFO not found users on datasource\n")
		return make([]*domain.User, 0), nil
	}

	return users, nil
}

func (r *userRepository) GetById(ctx context.Context, userId string) (*domain.User, error) {
	collection := r.mongoClient.Database(configuration.MONGO_DB_NAME).Collection(configuration.MONGO_USERS_DB_NAME)

	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		fmt.Printf("userRepository.GetById ERROR :  %s\n", err.Error())
	}

	var user *domain.User
	err = collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objectID}}).Decode(&user)

	if err != nil {
		fmt.Printf("userRepository.GetById ERROR :  %s\n", err.Error())
		return nil, err
	}

	fmt.Printf("userRepository.GetById INFO found user by given id: %s\n", userId)

	return user, nil

}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	collection := r.mongoClient.Database(configuration.MONGO_DB_NAME).Collection(configuration.MONGO_USERS_DB_NAME)

	var user *domain.User
	err := collection.FindOne(context.TODO(), bson.D{{Key: "email", Value: email}}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("userRepository.GetByEmail INFO : No documents found with the given email.")
			return nil, nil
		}

		fmt.Printf("userRepository.GetByEmail ERROR :  %s\n", err.Error())
		return nil, err
	}

	fmt.Printf("userRepository.GetByEmail INFO Not found user by given email: %s\n", email)

	return user, nil
}

func (r *userRepository) Upsert(ctx context.Context, user *domain.User) (string, error) {
	collection := r.mongoClient.Database(configuration.MONGO_DB_NAME).Collection(configuration.MONGO_USERS_DB_NAME)

	insertResult, err := collection.InsertOne(context.TODO(), user)

	if err != nil {
		return "", err
	}

	objectID := insertResult.InsertedID.(primitive.ObjectID)

	fmt.Printf("userRepository.Upsert INFO user saved with id: %s\n", objectID.Hex())

	return objectID.Hex(), nil
}
