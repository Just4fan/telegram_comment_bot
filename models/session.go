package models

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

const (
	SessionAlive               = 0x0
	SESSION_EXPIRED            = 0x1
	SESSION_ABORTED            = 0x2
	SESSION_ERROR              = 0x3
	SessionFinished            = 0x4
	SessionWaitingChannel      = 0x5
	SessionWaitingPost         = 0x6
	SessionWaitingEnablePost   = 0x7
	SessionWaitingInputComment = 0x8
)

var TypeCommentSession = reflect.TypeOf(&CommentSession{})
var TypeReloadSession = reflect.TypeOf(&ReloadSession{})

type BaseSession struct {
	ChatID int64 `json:"chat_id"`
	Status int8  `json:"status"`
}

type CommentSession struct {
	BaseSession
	Post *Post `json:"message"`
	//Area     int        `json:"area"`
	Panel int `json:"panel"`
	//Page     int        `json:"page"`
	//Total    int        `json:"total"`
	//Index    int        `json:"index"`
	//Params   string     `json:"params"`
	Comment *Comment `json:"comment"`
	//Comments []*Comment `json:"comments"`
}

type ReloadSession struct {
	BaseSession
}

type Session struct {
	Id     primitive.ObjectID `json:"id"`
	Status int8               `json:"status"`
	Chat   *tgbotapi.Chat     `json:"chat"`
}

func NewSession(chat *tgbotapi.Chat) *Session {
	return &Session{Id: primitive.NewObjectID(), Status: SessionAlive, Chat: chat}
}
