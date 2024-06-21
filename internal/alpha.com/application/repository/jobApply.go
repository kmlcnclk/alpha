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

type IJobApplyRepository interface {
	Get(ctx context.Context) ([]*domain.JobApply, error)
	Upsert(ctx context.Context, jobApply *domain.JobApply) error
}

type jobApplyRepository struct {
	mongoClient *mongo.Client
}

func NewJobApplyRepository(mongoClient *mongo.Client) IJobApplyRepository {
	return &jobApplyRepository{
		mongoClient: mongoClient,
	}
}

func (r *jobApplyRepository) Get(ctx context.Context) ([]*domain.JobApply, error) {
	collection := r.mongoClient.Database(configuration.MONGO_DB_NAME).Collection(configuration.MONGO_JOB_APPLIES_DB_NAME)

	var jobApplys []*domain.JobApply
	cursor, err := collection.Find(context.TODO(), bson.D{})

	if err != nil {
		fmt.Printf("jobApplyRepository.Get ERROR : %s\n", err.Error())
		return make([]*domain.JobApply, 0), err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var jobApply *domain.JobApply
		err := cursor.Decode(&jobApply)
		if err != nil {
			fmt.Printf("jobApplyRepository.Get ERROR : %s\n", err.Error())
			return make([]*domain.JobApply, 0), err
		}

		jobApplys = append(jobApplys, jobApply)
	}

	if err := cursor.Err(); err != nil {
		fmt.Printf("jobApplyRepository.Get ERROR : %s\n", err.Error())
	}

	if jobApplys == nil {
		fmt.Println("jobApplyRepository.Get INFO not found users on datasource")
		return make([]*domain.JobApply, 0), nil
	}

	return jobApplys, nil
}

func (r *jobApplyRepository) Upsert(ctx context.Context, jobApply *domain.JobApply) error {
	collection := r.mongoClient.Database(configuration.MONGO_DB_NAME).Collection(configuration.MONGO_JOB_APPLIES_DB_NAME)

	insertResult, err := collection.InsertOne(context.TODO(), jobApply)

	if err != nil {
		return err
	}

	objectID := insertResult.InsertedID.(primitive.ObjectID)

	fmt.Printf("jobApplyRepository.Upsert INFO user saved with id: %s\n", objectID.Hex())

	return nil
}
