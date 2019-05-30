package repositories

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"telegram_comment_bot/models"
)

const (
	collectionComments        = "comments"
	collectionCommentSections = "comment_sections"
)

type CommentRepository struct {
	comments []*models.Comment
	database *Database
}

func (r *CommentRepository) SelectAllComment() []*models.Comment {
	return r.comments
}

func (r *CommentRepository) FindCommentByID(id primitive.ObjectID) (comment *models.Comment, found bool) {
	comment = &models.Comment{}
	query := bson.M{"_id": id}
	err := r.database.FindOne(collectionComments, query, &options.FindOneOptions{}, comment)
	if err == nil {
		found = true
	}
	return
}

func (r *CommentRepository) FindCommentsByMessage(chatID int64, messageID, page, size int) (comments []*models.Comment) {
	limit := int64(size)
	skip := int64(size * (page - 1))
	opts := &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
		Sort:  bson.M{"date": -1},
	}
	query := bson.M{"post.chat_id": chatID, "post.message_id": messageID}
	err := r.database.FindMany(collectionComments, query, opts, &comments)
	if err != nil {
		log.Println(err)
	}
	return
}

func (r *CommentRepository) CountCommentsByMessage(chatID int64, messageID int) (count int64) {
	query := bson.M{"post.chat_id": chatID, "post.message_id": messageID}
	count, _ = r.database.Count(collectionComments, query, &options.CountOptions{})
	return
}

func (r *CommentRepository) InsertComment(comment *models.Comment) (ok bool) {
	_, err := r.database.InsertOne(collectionComments, comment)
	if err == nil {
		ok = true
	}
	return
}

func (r *CommentRepository) FindCommentSectionByMessage(chatID int64, messageID int) (post *models.CommentSection, found bool) {
	post = &models.CommentSection{}
	query := bson.M{"chat_id": chatID, "message_id": messageID}
	err := r.database.FindOne(collectionCommentSections, query, &options.FindOneOptions{}, post)
	if err != nil {
		log.Println(err)
	} else {
		found = true
	}
	return
}

func (r *CommentRepository) UpdateCommentSection(section *models.CommentSection) (ok bool) {
	query := bson.M{"chat_id": section.ChatID, "message_id": section.MessageID}
	update := bson.M{"section_id": section.SectionID}
	_, err := r.database.UpdateOne(collectionCommentSections, query, update, &options.UpdateOptions{})
	if err == nil {
		ok = true
	}
	return
}

func (r *CommentRepository) InsertCommentSection(section *models.CommentSection) (ok bool) {
	_, err := r.database.InsertOne(collectionCommentSections, section)
	if err == nil {
		ok = true
	}
	return
}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{database: GetDatabase("bot")}
}
