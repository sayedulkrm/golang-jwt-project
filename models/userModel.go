package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Go lang doesnt understand JSON, But MongoDB understands JSON. So the struct is a helper or a bridge to Mongo DB and Go lang

type User struct {
	ID primitive.ObjectID `bson:"_id"`

	// Id          string  `json:"_id,omitempty" bson:"_id,omitempty"`

	FirstName string `json:"first_name" validate:"required, min=2, max=100"`
	LastName  string `json:"last_name" validate:"required, min=2, max=100"`
	Password  string `json:"password" validate:"required, min=4"`
	Email     string `json:"email" validate:"email, required"`
	Phone     string `json:"phone" validate:"required"`
	Token     string `json:"token"`

	User_type     string    `json:"user_type" validate:"required, eq=admin|eq=user"`
	Refresh_token string    `json:"refresh_token"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
	User_status   string    `json:"user_status"`

	User_id string `json:"user_id"`
}
