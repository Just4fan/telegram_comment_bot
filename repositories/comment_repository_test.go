package repositories

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strconv"
	"telegram_comment_bot/models"
	"testing"
	"time"
)

var commentRepo = NewCommentRepository()

func TestCommentRepository_InsertComment(t *testing.T) {
	for i := 0; i < 13; i++ {
		ok := commentRepo.InsertComment(&models.Comment{
			ID:      primitive.NewObjectID(),
			Post:    &models.Post{ChatID: 2012, MessageID: 1},
			Date:    time.Now().Unix(),
			Content: "Content " + strconv.Itoa(i),
		})
		if ok {
			log.Printf("Inserted %d", i)
		}
	}
}

func TestCommentRepository_SelectCommentsByMessage(t *testing.T) {
	comments := commentRepo.FindCommentsByMessage(2012, 1, 4, 5)
	for _, comment := range comments {
		log.Printf("%+v\n", *comment)
	}
}
