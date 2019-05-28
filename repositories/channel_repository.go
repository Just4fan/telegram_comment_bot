package repositories

import (
	"comment_bot/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	creatorChannels = "creator_channels"
	userChannels    = "user_channels"
)

type ChannelRepository struct {
	//channels map[int64]*models.ChannelForCreator
	database *Database
}

func NewChannelRepository() *ChannelRepository {
	return &ChannelRepository{database: GetDatabase("bot")}
}

func (r *ChannelRepository) FindChannelForCreatorByChatID(chatID int64) (channel *models.ChannelForCreator, found bool) {
	channel = &models.ChannelForCreator{}
	err := r.database.FindOne(
		creatorChannels,
		bson.M{"chat_id": chatID},
		&options.FindOneOptions{}, channel)
	if err == nil {
		found = true
	}
	return
}

func (r *ChannelRepository) InsertChannelForCreator(channel *models.ChannelForCreator) (ok bool) {
	_, err := r.database.InsertOne(creatorChannels, channel)
	if err == nil {
		ok = true
	}
	return
}

func (r *ChannelRepository) FindChannelForUser(chatID int64, userID int) (channel *models.ChannelForUser, found bool) {
	channel = &models.ChannelForUser{}
	err := r.database.FindOne(
		userChannels,
		bson.M{"chat_id": chatID, "user_id": userID},
		&options.FindOneOptions{}, channel)
	if err == nil {
		found = true
	}
	return
}

func (r *ChannelRepository) InsertChannelForUser(channel *models.ChannelForUser) (ok bool) {
	_, err := r.database.InsertOne(userChannels, channel)
	if err == nil {
		ok = true
	}
	return
}
