package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
	Password []byte             `json:"password"`
}
