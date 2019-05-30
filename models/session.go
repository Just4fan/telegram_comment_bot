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

var TypeCommentSession = reflect.TypeOf(&AddCommentSession{})
var TypeReloadSession = reflect.TypeOf(&ReloadCommentSession{})
var TypeEnableSession = reflect.TypeOf(&EnableCommentSession{})

type BaseSession struct {
	ChatID int64 `json:"chat_id"`
	Status int8  `json:"status"`
}

type AddCommentSession struct {
	BaseSession
	Post    *Post              `json:"post" bson:"post"`
	ReplyTo primitive.ObjectID `json:"reply_to"`
}

type ReloadCommentSession struct {
	BaseSession
}

type EnableCommentSession struct {
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
