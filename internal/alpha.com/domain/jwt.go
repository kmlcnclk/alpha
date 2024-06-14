package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Jwt struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"userId" validate:"required"`
	AccessToken  string             `bson:"accessToken" validate:"required"`
	RefreshToken string             `bson:"refreshToken" validate:"required"`
	CreatedAt    time.Time          `bson:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"`
}
