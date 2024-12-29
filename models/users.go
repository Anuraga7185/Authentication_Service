package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Email         string             `bson:"email"`
	Password      string             `bson:"password"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
	UserType      string             `json:"user_type" bson:"user_type"`
	Address       string             `bson:"address"`
	Gender        string             `bson:"gender"`
	Biography     string             `bson:"bio_graphy"`
	DateOfBirth   time.Time          `bson:"date_of_birth"`
	FCMToken      string             `bson:"fcm_token"`
	IsPremiumUser bool               `bson:"is_premium_user"`
	Mobile        string             `bson:"mobile"`
	Name          string             `bson:"name"`
	UniqueId      string             `bson:"unique_id"`
}
