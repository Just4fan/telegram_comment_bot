package repositories

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"telegram_comment_bot/models"
	"testing"
)

var database = NewDatabase("bot")

func TestDatabase_InsertOne(t *testing.T) {
	channel := &models.ChannelForCreator{
		CreatorID: 510546065,
		ChatID:    -1001467892869,
		Settings: &models.ChannelSettingsForCreator{
			Mode: models.ModeAuto,
		},
	}
	_, err := database.InsertOne("creator_channels", channel)
	if err != nil {
		panic(err)
	}
}

func TestDatabase_InsertMany(t *testing.T) {
	var comments []*models.Comment
	err := database.FindMany("comments", bson.M{}, &options.FindOptions{}, &comments)
	if err != nil {
		panic(err)
	} else {
		var m []interface{}
		for _, v := range comments {
			m = append(m, bson.M{
				"_id": v.ID,
				"message": bson.M{
					"message_id": v.Post.MessageID,
					"chat_id":    v.Post.ChatID,
				},
				"reply_id":  v.ReplyID,
				"reply_to":  v.ReplyTo,
				"reply":     v.Reply,
				"user_id":   v.UserID,
				"user_name": v.UserName,
				"date":      v.Date,
				"content":   v.Content,
			})
		}
		ret, err := database.InsertMany("comments_backup", m, &options.InsertManyOptions{})
		if err != nil {
			panic(err)
		}
		log.Println(ret.InsertedIDs)
	}
}

func TestDatabase_FindMany(t *testing.T) {

}

func TestDatabase_DeleteMany(t *testing.T) {
	query := bson.M{"message.chat_id": 2012, "message.message_id": 1}
	ret, err := database.DeleteMany("comments", query)
	if err != nil {
		panic(err)
	}
	log.Println(ret.DeletedCount)
}

func TestDatabase_FindOne(t *testing.T) {
	query := bson.M{"message.chat_id": 2011, "message.message_id": 1}
	var comments []*models.Comment
	err := database.FindMany("comments", query, &options.FindOptions{}, &comments)
	if err != nil {
		panic(err)
	} else {
		for _, v := range comments {
			log.Printf("%+v\n", *v)
		}
	}
}
