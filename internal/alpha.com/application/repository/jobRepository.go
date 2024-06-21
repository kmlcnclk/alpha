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

type IJobRepository interface {
	Get(ctx context.Context) ([]*domain.Job, error)
	Upsert(ctx context.Context, job *domain.Job) error
	GetByIDAndBusinessAccountID(ctx context.Context, id, businessAccountID string) (*domain.Job, error)
}

type jobRepository struct {
	mongoClient *mongo.Client
}

func NewJobRepository(mongoClient *mongo.Client) IJobRepository {
	return &jobRepository{
		mongoClient: mongoClient,
	}
}

func (r *jobRepository) Get(ctx context.Context) ([]*domain.Job, error) {
	collection := r.mongoClient.Database(configuration.MONGO_DB_NAME).Collection(configuration.MONGO_JOBS_DB_NAME)

	var jobs []*domain.Job
	cursor, err := collection.Find(context.TODO(), bson.D{})

	if err != nil {
		fmt.Printf("jobRepository.Get ERROR : %s\n", err.Error())
		return make([]*domain.Job, 0), err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var job *domain.Job
		err := cursor.Decode(&job)
		if err != nil {
			fmt.Printf("jobRepository.Get ERROR : %s\n", err.Error())
			return make([]*domain.Job, 0), err
		}

		jobs = append(jobs, job)
	}

	if err := cursor.Err(); err != nil {
		fmt.Printf("jobRepository.Get ERROR : %s\n", err.Error())
	}

	if jobs == nil {
		fmt.Println("jobRepository.Get INFO not found users on datasource")
		return make([]*domain.Job, 0), nil
	}

	return jobs, nil
}

func (r *jobRepository) Upsert(ctx context.Context, job *domain.Job) error {
	collection := r.mongoClient.Database(configuration.MONGO_DB_NAME).Collection(configuration.MONGO_JOBS_DB_NAME)

	insertResult, err := collection.InsertOne(context.TODO(), job)

	if err != nil {
		return err
	}

	objectID := insertResult.InsertedID.(primitive.ObjectID)

	fmt.Printf("jobRepository.Upsert INFO user saved with id: %s\n", objectID.Hex())

	return nil
}

func (r *jobRepository) GetByIDAndBusinessAccountID(ctx context.Context, id, businessAccountID string) (*domain.Job, error) {
	collection := r.mongoClient.Database(configuration.MONGO_DB_NAME).Collection(configuration.MONGO_JOBS_DB_NAME)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Printf("jobRepository.GetByIDAndUserID ERROR :  %s\n", err.Error())
		return nil, err
	}

	objectIDForBusinessAccount, err := primitive.ObjectIDFromHex(businessAccountID)
	if err != nil {
		fmt.Printf("jobRepository.GetByIDAndUserID ERROR :  %s\n", err.Error())
		return nil, err
	}

	filter := bson.M{"_id": objectID, "businessAccountId": objectIDForBusinessAccount}

	var job *domain.Job
	err = collection.FindOne(context.Background(), filter).Decode(&job)

	if err != nil {
		return nil, err
	}

	return job, nil
}
