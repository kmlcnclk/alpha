package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Jwt struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `json:"userId" validate:"required"`
	AccessToken  string             `json:"accessToken" validate:"required"`
	RefreshToken string             `json:"refreshToken" validate:"required"`
}
