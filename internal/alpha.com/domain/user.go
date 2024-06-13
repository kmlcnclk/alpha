package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `json:"firstName" validate:"required"`
	LastName  string             `json:"lastName" validate:"required"`
	Email     string             `bson:"email" validate:"required,email"`
	Password  string             `bson:"password" validate:"required,min=6"`
	Age       int32              `bson:"age" validate:"gte=0,lte=130"`
}
