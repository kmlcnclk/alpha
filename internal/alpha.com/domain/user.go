package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"firstName" validate:"required"`
	LastName  string             `bson:"lastName" validate:"required"`
	Email     string             `bson:"email" validate:"required,email"`
	Password  string             `bson:"password" validate:"required,min=6"`
	Age       int32              `bson:"age" validate:"gte=0,lte=130"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}
