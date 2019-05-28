package utils

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strconv"
	"strings"
	"telegram_comment_bot/models"
)

func EncodeData(chatID int64, messageID int, panelID int) string {
	return strconv.FormatInt(chatID, 10) + "/" + strconv.Itoa(messageID) + "/" + strconv.Itoa(panelID)
}

func EncodeDataR(id primitive.ObjectID, panelID int) string {
	return id.Hex() + "/" + strconv.Itoa(panelID)
}

func EncodeDataI(page int, chatID int64, messageID int, panelID int) string {
	return strconv.Itoa(page) + "/" + strconv.FormatInt(chatID, 10) + "/" + strconv.Itoa(messageID) + "/" + strconv.Itoa(panelID)
}

func DecodeDataI(data string) (page int, chatID int64, messageID int, panelID int, err error) {
	params := strings.Split(data, "/")
	if len(params) != 4 {
		err = errors.New("params' count no equal to 4")
		return
	}
	page, err = strconv.Atoi(params[0])
	if err != nil {
		return
	}
	chatID, err = strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return
	}
	messageID, err = strconv.Atoi(params[2])
	if err != nil {
		return
	}
	panelID, err = strconv.Atoi(params[3])
	if err != nil {
		return
	}
	return
}

func DecodeData(data string) (chatID int64, messageID int, panelID int, err error) {
	params := strings.Split(data, "/")
	if len(params) != 3 {
		err = errors.New("params' count no equal to 3")
		return
	}
	chatID, err = strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		return
	}
	messageID, err = strconv.Atoi(params[1])
	if err != nil {
		return
	}
	panelID, err = strconv.Atoi(params[2])
	if err != nil {
		return
	}
	return
}

func DecodeDataR(data string) (id primitive.ObjectID, panelID int, err error) {
	params := strings.Split(data, "/")
	if len(params) != 2 {
		err = errors.New("params' count no equal to 2")
		return
	}
	id, err = primitive.ObjectIDFromHex(params[0])
	if err != nil {
		return
	}
	panelID, err = strconv.Atoi(params[1])
	if err != nil {
		return
	}
	return
}

func EncodeParam(post *models.Post) (param string, err error) {
	param += strconv.FormatInt(post.ChatID, 10) + "_" + strconv.Itoa(post.MessageID)
	/*data, err := json.Marshal(post)
	if err != nil {
		return
	}
	param, err = Base64Encode(data)*/
	return
}

func DecodeParam(param string) (post *models.Post, err error) {
	post = &models.Post{}
	params := strings.Split(param, "_")
	if len(params) != 2 {
		err = errors.New("params' count not equal to 2")
	}
	post.ChatID, err = strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		return
	}
	post.MessageID, err = strconv.Atoi(params[1])
	if err != nil {
		return
	}
	return
	/*post = &models.Post{}
	data, err := Base64Decode(param)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, post)*/
}

func PageSwitchKeyboardRow(page int, chatID int64, messageID int, panelID int) (row []tgbotapi.InlineKeyboardButton) {
	prePage := "pg " + EncodeDataI(page-1, chatID, messageID, panelID)
	row = append(row, tgbotapi.InlineKeyboardButton{Text: "⬅️上一页", CallbackData: &prePage})
	curPage := "pg " + EncodeDataI(page, chatID, messageID, panelID)
	row = append(row, tgbotapi.InlineKeyboardButton{Text: "🔄刷新", CallbackData: &curPage})
	nextPage := "pg " + EncodeDataI(page+1, chatID, messageID, panelID)
	row = append(row, tgbotapi.InlineKeyboardButton{Text: "下一页➡️", CallbackData: &nextPage})
	return
}

func CommentPostKeyBoardRow(chatID int64, messageID int, panelID int) (row []tgbotapi.InlineKeyboardButton) {
	data := EncodeData(chatID, messageID, panelID)
	commentPost := "cp " + data
	row = append(row, tgbotapi.InlineKeyboardButton{Text: "点击此按钮评论主贴或选择序号评论回复", CallbackData: &commentPost})
	return
}

func CommentIndexKeyBoardRow(comments []*models.Comment, panelID int) (row []tgbotapi.InlineKeyboardButton) {
	for i, comment := range comments {
		row = append(row, ReplyKeyBoardRow(strconv.Itoa(i+1), comment, panelID))
	}
	return
}

func QuickReplyKeyBoardRow(comment *models.Comment, panelID int) (row []tgbotapi.InlineKeyboardButton) {
	row = append(row, ReplyKeyBoardRow("回复", comment, panelID))
	return
}

func ReplyKeyBoardRow(text string, comment *models.Comment, panelID int) (btn tgbotapi.InlineKeyboardButton) {
	callback := "cr " + EncodeDataR(comment.ID, panelID)
	btn = tgbotapi.InlineKeyboardButton{Text: text, CallbackData: &callback}
	log.Printf("comment" + text + " callback:" + callback)
	return
}

func NewInlineKeyboardMarkup(rows ...[]tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup {
	var keyboard [][]tgbotapi.InlineKeyboardButton
	for _, row := range rows {
		if len(row) > 0 {
			keyboard = append(keyboard, row)
		}
	}
	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}
