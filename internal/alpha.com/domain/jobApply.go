package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobApply struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	JobID     primitive.ObjectID `bson:"jobId" validate:"required"`
	UserID    primitive.ObjectID `bson:"userId" validate:"required"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}
