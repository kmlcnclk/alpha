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

type IBusinessAccountRepository interface {
	Get(ctx context.Context) ([]*domain.BusinessAccount, error)
	GetByID(ctx context.Context, businessAccountId string) (*domain.BusinessAccount, error)
	Upsert(ctx context.Context, businessAccount *domain.BusinessAccount) error
	GetByIDAndUserID(ctx context.Context, businessAccountId string, userID string) (*domain.BusinessAccount, error)
}

type businessAccountRepository struct {
	mongoClient *mongo.Client
}

func NewBusinessAccountRepository(mongoClient *mongo.Client) IBusinessAccountRepository {
	return &businessAccountRepository{
		mongoClient: mongoClient,
	}
}

func (r *businessAccountRepository) Get(ctx context.Context) ([]*domain.BusinessAccount, error) {

	collection := r.mongoClient.Database(configuration.MONGO_DB_NAME).Collection(configuration.MONGO_BUSINESS_ACCOUNT_DB_NAME)

	var businessAccounts []*domain.BusinessAccount
	cursor, err := collection.Find(context.TODO(), bson.D{})

	if err != nil {
		fmt.Printf("businessAccountRepository.Get ERROR : %s\n", err.Error())
		return make([]*domain.BusinessAccount, 0), err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var businessAccount *domain.BusinessAccount
		err := cursor.Decode(&businessAccount)
		if err != nil {
			fmt.Printf("businessAccountRepository.Get ERROR : %s\n", err.Error())
			return make([]*domain.BusinessAccount, 0), err
		}

		businessAccounts = append(businessAccounts, businessAccount)
	}

	if err := cursor.Err(); err != nil {
		fmt.Printf("businessAccountRepository.Get ERROR : %s\n", err.Error())
	}

	if businessAccounts == nil {
		fmt.Println("businessAccountRepository.Get INFO not found users on datasource")
		return make([]*domain.BusinessAccount, 0), nil
	}

	return businessAccounts, nil
}

func (r *businessAccountRepository) Upsert(ctx context.Context, businessAccount *domain.BusinessAccount) error {
	collection := r.mongoClient.Database(configuration.MONGO_DB_NAME).Collection(configuration.MONGO_BUSINESS_ACCOUNT_DB_NAME)

	insertResult, err := collection.InsertOne(context.TODO(), businessAccount)

	if err != nil {
		return err
	}

	objectID := insertResult.InsertedID.(primitive.ObjectID)

	fmt.Printf("businessAccountRepository.Upsert INFO user saved with id: %s\n", objectID.Hex())

	return nil
}

func (r *businessAccountRepository) GetByID(ctx context.Context, businessAccountId string) (*domain.BusinessAccount, error) {
	collection := r.mongoClient.Database(configuration.MONGO_DB_NAME).Collection(configuration.MONGO_BUSINESS_ACCOUNT_DB_NAME)

	objectID, err := primitive.ObjectIDFromHex(businessAccountId)
	if err != nil {
		fmt.Printf("businessAccountRepository.GetByID ERROR :  %s\n", err.Error())
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var businessAccount *domain.BusinessAccount
	err = collection.FindOne(context.Background(), filter).Decode(&businessAccount)

	if err != nil {
		return nil, err
	}

	return businessAccount, nil
}

func (r *businessAccountRepository) GetByIDAndUserID(ctx context.Context, businessAccountId string, userID string) (*domain.BusinessAccount, error) {
	collection := r.mongoClient.Database(configuration.MONGO_DB_NAME).Collection(configuration.MONGO_BUSINESS_ACCOUNT_DB_NAME)

	objectID, err := primitive.ObjectIDFromHex(businessAccountId)
	if err != nil {
		fmt.Printf("businessAccountRepository.GetByIDAndUserID ERROR :  %s\n", err.Error())
		return nil, err
	}

	objectIDForUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		fmt.Printf("businessAccountRepository.GetByIDAndUserID ERROR :  %s\n", err.Error())
		return nil, err
	}

	fmt.Println(objectID, objectIDForUser)

	filter := bson.M{"_id": objectID, "userId": objectIDForUser}

	var businessAccount *domain.BusinessAccount
	err = collection.FindOne(context.Background(), filter).Decode(&businessAccount)

	if err != nil {
		return nil, err
	}

	return businessAccount, nil
}
