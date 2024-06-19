package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	Id                primitive.ObjectID `bson:"_id,omitempty"`
	BusinessAccountID primitive.ObjectID `bson:"businessAccountId" validate:"required"`
	Name              string             `bson:"name" validate:"required"`
	Description       string             `bson:"description" validate:"required"`
	Price             float32            `bson:"price" validate:"required"`
	Category          string             `bson:"category" validate:"required"`
	CreatedAt         time.Time          `bson:"createdAt"`
	UpdatedAt         time.Time          `bson:"updatedAt"`
}
