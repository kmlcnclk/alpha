package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BusinessAccount struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      primitive.ObjectID `bson:"userId" validate:"required"`
	Name        string             `bson:"name" validate:"required"`
	Description string             `bson:"description" validate:"required"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}
