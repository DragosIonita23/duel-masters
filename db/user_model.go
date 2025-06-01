package db

import "go.mongodb.org/mongo-driver/mongo"

type UserSession struct {
	Token   string `json:"token"`
	IP      string `json:"ip"`
	Expires int    `json:"expires"`
}

type User struct {
	UID         string        `json:"uid"`
	Permissions []string      `json:"permissions"`
	Username    string        `json:"username"`
	Password    string        `json:"-"`
	Email       string        `json:"email"`
	Color       string        `json:"color"`
	Playmat     string        `json:"playmat"`
	Sessions    []UserSession `json:"-"`
	MutedUsers  []string      `json:"muted_users"`
	Chatblocked bool          `json:"-" bson:"chat_blocked"`
}

func Users() *mongo.Collection {
	return conn().Collection("users")
}
