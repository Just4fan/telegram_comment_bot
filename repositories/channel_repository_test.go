package repositories

import (
	"log"
	"telegram_comment_bot/models"
	"testing"
)

var channelRepo = NewChannelRepository()

func TestChannelRepository_FindChannelForCreatorByChatID(t *testing.T) {
	channel, found := channelRepo.FindChannelForCreatorByChatID(-1001467892869)
	if found {
		log.Printf("%+v\n", channel)
	}
}

func TestChannelRepository_InsertChannelForCreator(t *testing.T) {
	ok := channelRepo.InsertChannelForCreator(&models.ChannelForCreator{
		CreatorID: 510546065,
		ChatID:    -1001467892869,
		Settings: &models.ChannelSettingsForCreator{
			Mode: models.ModeAuto,
		},
	})
	if ok {
		log.Println("ok")
	}
}

func TestChannelRepository_FindChannelForUser(t *testing.T) {
	channel, found := channelRepo.FindChannelForUser(-1001467892869, 888197665)
	if found {
		log.Printf("%+v\n", channel)
	}
}

func TestChannelRepository_InsertChannelForUser(t *testing.T) {
	ok := channelRepo.InsertChannelForUser(&models.ChannelForUser{
		UserID: 888197665,
		ChatID: -1001467892869,
		Settings: &models.ChannelSettingsForUser{
			ReceiveNotification: true,
		},
	})
	if ok {
		log.Println("ok")
	}
}
