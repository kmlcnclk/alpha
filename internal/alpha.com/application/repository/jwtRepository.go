package repository

import (
	"context"
	"fmt"
	"time"

	"alpha.com/configuration"
	"alpha.com/internal/alpha.com/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IJwtRepository interface {
	Get(ctx context.Context) ([]*domain.Jwt, error)
	GetById(ctx context.Context, userId string) (*domain.Jwt, error)
	Upsert(ctx context.Context, user *domain.Jwt) error
	Update(ctx context.Context, userID, accessToken, refreshToken string) error
}

type jwtRepository struct {
	mongoClient *mongo.Client
}

func NewJwtRepository(mongoClient *mongo.Client) IJwtRepository {
	return &jwtRepository{
		mongoClient: mongoClient,
	}
}

func (r *jwtRepository) Get(ctx context.Context) ([]*domain.Jwt, error) {

	collection := r.mongoClient.Database(configuration.MONGO_DB_NAME).Collection(configuration.MONGO_JWT_DB_NAME)

	var jwts []*domain.Jwt
	cursor, err := collection.Find(context.TODO(), bson.D{})

	if err != nil {
		fmt.Printf("jwtRepository.Get ERROR : %s\n", err.Error())
		return make([]*domain.Jwt, 0), err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var jwt *domain.Jwt
		err := cursor.Decode(&jwt)
		if err != nil {
			fmt.Printf("jwtRepository.Get ERROR : %s\n", err.Error())
			return make([]*domain.Jwt, 0), err
		}

		jwts = append(jwts, jwt)
	}

	if err := cursor.Err(); err != nil {
		fmt.Printf("jwtRepository.Get ERROR : %s\n", err.Error())
	}

	if jwts == nil {
		fmt.Println("jwtRepository.Get INFO not found users on datasource")
		return make([]*domain.Jwt, 0), nil
	}

	return jwts, nil
}

func (r *jwtRepository) GetById(ctx context.Context, userId string) (*domain.Jwt, error) {
	return nil, nil
}

func (r *jwtRepository) Upsert(ctx context.Context, jwt *domain.Jwt) error {
	collection := r.mongoClient.Database(configuration.MONGO_DB_NAME).Collection(configuration.MONGO_JWT_DB_NAME)

	insertResult, err := collection.InsertOne(context.TODO(), jwt)

	if err != nil {
		return err
	}

	objectID := insertResult.InsertedID.(primitive.ObjectID)

	fmt.Printf("jwtRepository.Upsert INFO user saved with id: %s\n", objectID.Hex())

	return nil
}

func (r *jwtRepository) Update(ctx context.Context, userID, accessToken, refreshToken string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		fmt.Printf("jwtRepository.Update ERROR :  %s\n", err.Error())
	}

	filter := bson.M{"refreshToken": refreshToken, "userId": objectID}
	update := bson.M{
		"$set": bson.M{
			"accessToken": accessToken,
			"updatedAt":   time.Now(),
		},
	}

	collection := r.mongoClient.Database(configuration.MONGO_DB_NAME).Collection(configuration.MONGO_JWT_DB_NAME)

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
