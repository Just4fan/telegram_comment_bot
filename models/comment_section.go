package models

type CommentSection struct {
	MessageID int   `json:"message_id" bson:"message_id"`
	ChatID    int64 `json:"chat_id" bson:"chat_id"`
	SectionID int   `json:"section_id" bson:"section_id"`
}
