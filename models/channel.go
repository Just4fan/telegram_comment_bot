package models

const (
	MODE_MANUAL = 0x0
	ModeAuto    = 0x1
)

type ChannelSettingsForCreator struct {
	Mode int8 `json:"mode" bson:"mode"`
}

type ChannelForCreator struct {
	ChatID    int64                      `json:"channel_id" bson:"chat_id"`
	CreatorID int                        `json:"creator_id" bson:"creator_id"`
	Link      string                     `json:"link" bson:"link"`
	Settings  *ChannelSettingsForCreator `json:"settings" bson:"settings"`
}

type ChannelSettingsForUser struct {
	ReceiveNotification bool `json:"receive_notification" bson:"receive_notification"`
}

type ChannelForUser struct {
	ChatID   int64                   `json:"chat_id" bson:"chat_id"`
	UserID   int                     `json:"user_id" bson:"user_id"`
	Settings *ChannelSettingsForUser `json:"settings" bson:"settings"`
}
