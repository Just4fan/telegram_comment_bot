package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	MessageID int   `json:"message_id" bson:"message_id"`
	ChatID    int64 `json:"chat_id" bson:"chat_id"`
	AreaID    int   `json:"area_id" bson:"area_id"`
}

type Comment struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Post     *Post              `json:"post" bson:"post"`
	ReplyID  int                `json:"reply_id" bson:"reply_id"`
	ReplyTo  string             `json:"reply_to" bson:"reply_to"`
	Reply    primitive.ObjectID `json:"reply" bson:"reply"`
	UserID   int                `json:"user_id" bson:"user_id"`
	UserName string             `json:"user_name" bson:"user_name"`
	Date     int64              `json:"date" bson:"date"`
	Content  string             `json:"content" bson:"content"`
}
