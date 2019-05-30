package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"reflect"
	"strings"
	"telegram_comment_bot/configs"
	"telegram_comment_bot/models"
	"telegram_comment_bot/repositories"
	"telegram_comment_bot/utils"
	"time"
)

const (
	CommandRegister = "register"
	CommandStart    = "start"
	CommandReload   = "reload"
	CommandCancel   = "cancel"
	CommandEnable   = "enable"

	EnableCommentSection  = int8(0x0)
	ReloadCommentSection  = int8(0x1)
	RecoverCommentSection = int8(0x2)

	PageSize = 5
)

type Handler struct {
	bot         *tgbotapi.BotAPI
	repo        *repositories.Repository
	channelRepo *repositories.ChannelRepository
	sessionRepo *repositories.SessionRepository
	commentRepo *repositories.CommentRepository
}

func (h *Handler) Pause() {
	updates := h.bot.ListenForWebhook("/" + h.bot.Token)
	go http.ListenAndServe("0.0.0.0:8080", nil)

	for update := range updates {
		if update.Message != nil || update.CallbackQuery != nil {
			var chatID int64
			if update.Message != nil {
				chatID = update.Message.Chat.ID
			} else {
				chatID = int64(update.CallbackQuery.From.ID)
			}
			h.bot.Send(tgbotapi.NewMessage(chatID, "机器人正在搬家"))
		}
	}
}

func (h *Handler) Start() {
	updates := h.bot.ListenForWebhook("/" + h.bot.Token)
	go http.ListenAndServe("0.0.0.0:8080", nil)

	info, err := h.bot.GetWebhookInfo()

	if err != nil {
		log.Printf("%+v\n", info)
	}

	for update := range updates {
		/*		config := tgbotapi.NewUpdate(offset)
				config.Timeout = 60
				updates, err := h.bot.GetUpdatesChan(config)
				if err != nil {
					continue
				}
				for update := range updates {*/

		//log.Printf("%+v\n", update.Message.ForwardFrom)
		if update.ChannelPost != nil {
			h.handleChannelMessage(update.ChannelPost)
		} else if update.EditedChannelPost != nil {
			log.Println("Edited Channel Post")
		} else if update.EditedMessage != nil {
			log.Println("Edited Post")
		} else if update.CallbackQuery != nil {
			h.handleCallbackQuery(update.CallbackQuery)
		} else if update.Message != nil {
			handled := false
			if update.Message.IsCommand() {
				log.Println(update.Message.Command())
				switch update.Message.Command() {
				case CommandRegister:
					h.handleRegisterRequest(update.Message)
				case CommandStart:
					h.handleStartRequest(update.Message)
				case CommandReload:
					h.handleReloadRequest(update.Message)
				case CommandCancel:
					h.handleCancelRequest(update.Message)
				case CommandEnable:
					h.handleEnableRequest(update.Message)
				}
				handled = true
			}
			if !handled {
				h.handleMessage(update.Message)
			}
		}
		//offset = update.UpdateID
		//}
	}
}

func (h *Handler) updateCommentPanel(postChat int64, postMessage int, page int, chatID int64, panelID int, updateOnEmpty bool) bool {

	comments := h.commentRepo.FindCommentsByMessage(postChat, postMessage, page, PageSize)

	if !updateOnEmpty && len(comments) == 0 {
		return false
	}

	keyboard := utils.NewInlineKeyboardMarkup(
		utils.CommentPostKeyBoardRow(postChat, postMessage),
		utils.CommentIndexKeyBoardRow(comments),
		utils.PageSwitchKeyboardRow(page, postChat, postMessage, panelID))

	editConfig := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      chatID,
			MessageID:   panelID,
			ReplyMarkup: &keyboard,
		},
		Text:      utils.CommentsText(comments),
		ParseMode: tgbotapi.ModeHTML,
	}

	_, err := h.bot.Send(editConfig)
	if err != nil {
		log.Println(err)
		return false
	}

	return true

}

func (h *Handler) updateCommentSection(section *models.CommentSection) bool {

	comments := h.commentRepo.FindCommentsByMessage(section.ChatID, section.MessageID, 1, PageSize)
	commentsText := utils.CommentsText(comments)
	params, err := utils.EncodeParam(section.ChatID, section.MessageID)
	if err != nil {
		log.Println(err)
		return false
	}

	url := "http://t.me/" + h.bot.Self.UserName + "?start=" + params

	keyboard := utils.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{Text: "查看详情️", URL: &url}))

	edit := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      section.ChatID,
			MessageID:   section.SectionID,
			ReplyMarkup: &keyboard,
		},
		Text:      commentsText,
		ParseMode: tgbotapi.ModeHTML,
	}

	log.Printf("%+v\n", edit)
	_, err = h.bot.Send(edit)

	if err != nil {
		log.Println(err)
		return false
	}

	return true

}

func (h *Handler) resolveCommentSection(chatID int64, messageID int, mode int8) bool {

	//var comments []*models.Comment
	section, found := h.commentRepo.FindCommentSectionByMessage(chatID, messageID)

	if (!found && mode == EnableCommentSection) || (found && mode == RecoverCommentSection) {
		if !found {
			section = &models.CommentSection{
				ChatID:    chatID,
				MessageID: messageID,
			}
		}

		replyConfig := &tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:              chatID,
				ReplyToMessageID:    messageID,
				DisableNotification: true,
			},
			Text: "加载中",
		}

		msg, err := h.bot.Send(replyConfig)
		if err != nil {
			log.Println(err)
			return false
		}

		section.SectionID = msg.MessageID

		if !found {

			ok := h.commentRepo.InsertCommentSection(section)

			if !ok {
				return false
			}

		} else {

			ok := h.commentRepo.UpdateCommentSection(section)

			if !ok {
				return false
			}

		}

	} else if (!found && mode == ReloadCommentSection) || (!found && mode == RecoverCommentSection) || (found && mode == EnableCommentSection) {
		return false
	}

	return h.updateCommentSection(section)

}

func (h *Handler) handleEnableRequest(message *tgbotapi.Message) {
	h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "转发开启评论的频道消息到此对话"))
	h.sessionRepo.InsertSession(message.Chat.ID, &models.BaseSession{ChatID: message.Chat.ID, Status: models.SessionWaitingPost})
}

func (h *Handler) handleCancelRequest(message *tgbotapi.Message) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	h.sessionRepo.RemoveSession(message.Chat.ID)
	h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "会话状态已清空"))
}

func (h *Handler) handleReloadRequest(message *tgbotapi.Message) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: message.Chat.ID,
		},
		Text: `⚠️此功能将尝试修复评论功能异常的帖子⚠️
				转发评论功能异常的频道消息到此或者输入指令 /cancel 取消操作`,
	}
	h.bot.Send(msg)
	session := &models.ReloadCommentSession{BaseSession: models.BaseSession{
		ChatID: message.Chat.ID,
		Status: models.SessionWaitingPost,
	}}
	h.sessionRepo.InsertSession(session.ChatID, session)
}

func (h *Handler) handleStartRequest(message *tgbotapi.Message) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	log.Println(message.CommandArguments())
	if len(message.CommandArguments()) != 0 {
		h.handleCommentRequest(message)
	}
}

//Extract Edit Reply Method
func (h *Handler) handleCommentRequest(message *tgbotapi.Message) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	chatID, messageID, err := utils.DecodeParam(message.CommandArguments())
	if err != nil {
		log.Println(err)
		return
	}

	_, found := h.channelRepo.FindChannelForCreatorByChatID(chatID)
	if !found {
		h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "频道未注册！"))
	}

	//var msg tgbotapi.Message
	forwardConfig := tgbotapi.NewForward(message.Chat.ID, chatID, messageID)
	forwardMsg, err := h.bot.Send(forwardConfig)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("channel:%+v\n", *forwardMsg.ForwardFromChat)
	msgConfig := &tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           message.Chat.ID,
			ReplyToMessageID: forwardMsg.MessageID,
		},
		Text: "加载中",
	}

	msg, err := h.bot.Send(msgConfig)
	if err != nil {
		log.Println(err)
		return
	}

	h.updateCommentPanel(chatID, messageID, 1, message.Chat.ID, msg.MessageID, true)

}

func (h *Handler) handleCallbackQuery(query *tgbotapi.CallbackQuery) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	arr := strings.Split(query.Data, " ")
	log.Println(arr)
	/*	if query.Post == nil {
			return
		}
		session, found := h.sessionRepo.SelectCommentSessionByChatID(query.Post.Chat.ID)
		if !found || !(session.Panel == query.Post.MessageID) {
			return
		}*/
	switch arr[0] {

	case "cp", "cr":

		var chatID int64
		var messageID int
		var comment *models.Comment

		session := &models.AddCommentSession{
			BaseSession: models.BaseSession{
				ChatID: int64(query.From.ID),
				Status: models.SessionWaitingInputComment,
			},
		}

		if arr[0] == "cp" {

			var err error
			chatID, messageID, err = utils.DecodeData(arr[1])
			if err != nil {
				log.Println(err)
				return
			}

			session.Post = &models.Post{
				ChatID:    chatID,
				MessageID: messageID,
			}

		} else {

			var id primitive.ObjectID
			var err error
			id, err = utils.DecodeDataR(arr[1])
			if err != nil {
				log.Println(err)
				return
			}

			session.ReplyTo = id

			var found bool
			comment, found = h.commentRepo.FindCommentByID(id)
			if !found {
				h.bot.Send(tgbotapi.NewMessage(int64(query.From.ID), "找不到该回复"))
				return
			}

		}

		h.sessionRepo.InsertSession(session.ChatID, session)
		text := ""

		if comment != nil {
			text += "评论用户 <b>" + comment.UserName + "</b>:\n<pre>\n</pre>" + "<pre>      </pre><i>" + comment.Content + "</i>\n<pre>\n</pre>"
		}

		text += "输入评论并发送，发送 /cancel 指令退出评论:"

		config := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: session.ChatID,
			},
			Text:      text,
			ParseMode: tgbotapi.ModeHTML,
		}

		_, err := h.bot.Send(config)

		if err != nil {
			log.Println(err)
		}

		h.bot.AnswerCallbackQuery(tgbotapi.NewCallback(query.ID, ""))

	case "pg":

		page, chatID, messageID, panelID, err := utils.DecodeDataI(arr[1])

		if err != nil {
			log.Println(err)
			return
		}

		if page > 0 {

			//comments := h.commentRepo.FindCommentsByMessage(chatID, messageID, page, 5)

			if h.updateCommentPanel(chatID, messageID, page, int64(query.From.ID), panelID, false) {

				h.bot.AnswerCallbackQuery(tgbotapi.NewCallback(query.ID, ""))
				return

			}
		}

		h.bot.AnswerCallbackQuery(tgbotapi.NewCallback(query.ID, "所请求页面没有更多数据"))
	}
}

func (h *Handler) handleMessage(message *tgbotapi.Message) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	v, found := h.sessionRepo.SelectSessionByChatID(message.Chat.ID)
	if found {
		typ := reflect.TypeOf(v)
		log.Println(typ)
		switch typ {

		case models.TypeEnableSession:

			session, _ := v.(*models.EnableCommentSession)
			switch session.Status {
			case models.SessionWaitingPost:

			}

		case models.TypeReloadSession:

			session, _ := v.(*models.ReloadCommentSession)
			switch session.Status {
			case models.SessionWaitingPost:
				if message.ForwardFromChat == nil {
					h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "转发频道消息！！"))
					return
				}
				channel, found := h.channelRepo.FindChannelForCreatorByChatID(message.ForwardFromChat.ID)
				if !found {
					h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "该频道未注册"))
					return
				}
				if channel.CreatorID != message.From.ID {
					h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "没有权限"))
					return
				}

				ok := h.resolveCommentSection(message.ForwardFromChat.ID, message.ForwardFromMessageID, ReloadCommentSection)

				if ok {

					h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "评论功能已重新加载"))
					h.sessionRepo.RemoveSession(message.Chat.ID)
				}
			}

		case models.TypeCommentSession:

			session, _ := v.(*models.AddCommentSession)
			log.Printf("%+v\n", session)
			switch session.Status {
			case models.SessionWaitingInputComment:

				//var replyTo bool
				var replyTo *models.Comment

				comment := &models.Comment{
					ID:       primitive.NewObjectID(),
					Date:     time.Now().Unix(),
					Content:  message.Text,
					UserID:   message.From.ID,
					UserName: message.From.FirstName + " " + message.From.LastName,
				}

				if session.ReplyTo.Hex() != "000000000000000000000000" {
					//replyTo = true
					var found bool
					replyTo, found = h.commentRepo.FindCommentByID(session.ReplyTo)
					if !found {
						log.Println("ReplyTo Not Found")
						return
					}
					comment.Post = replyTo.Post
					comment.Reply = replyTo.ID
					comment.ReplyTo = replyTo.UserName
					comment.ReplyID = replyTo.UserID
				} else {
					comment.Post = session.Post
				}

				section, found := h.commentRepo.FindCommentSectionByMessage(comment.Post.ChatID, comment.Post.MessageID)
				if !found {
					log.Println("Post Not Found")
					return
				}

				log.Printf("%+v", comment)
				ok := h.commentRepo.InsertComment(comment)
				if !ok {
					return
				}

				h.updateCommentSection(section)

				if replyTo != nil && comment.ReplyID != comment.UserID {

					markup := utils.NewInlineKeyboardMarkup(utils.QuickReplyKeyBoardRow(comment))
					notification := tgbotapi.MessageConfig{
						BaseChat: tgbotapi.BaseChat{
							ChatID:      int64(comment.ReplyID),
							ReplyMarkup: markup,
						},
						Text:      utils.NotificationText(comment.UserName, replyTo.Content, comment.Content),
						ParseMode: tgbotapi.ModeHTML,
					}
					_, err := h.bot.Send(notification)
					if err != nil {
						log.Println(err)
						h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "评论成功，但是对方可能关闭了跟我的对话，通知无法送达"))
						return
					}

				}
				h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "评论成功"))
				session.Status = models.SessionFinished
			}
		}
	}
	/*	session, found := h.repo.SelectSession(message.Chat.ID)
		if found {
			switch session.Status {
			case models.SESSION_WAITING_CHANNEL:
				var reply tgbotapi.MessageConfig
				if message.ForwardFromChat == nil {
					reply = tgbotapi.NewMessage(message.Chat.ID, "从频道转发一条消息到此对话")
				} else {
					//log.Println(message.ForwardFromChat.ID)
					admins, err := h.bot.GetChatAdministrators(tgbotapi.ChatConfig{ChatID: message.ForwardFromChat.ID})
					isCreator := false
					isAdmin := false
					if err == nil {
						for _, admin := range admins {
							log.Println(admin.User.UserName)
							if admin.User.ID == message.From.ID && admin.IsCreator() {
								isCreator = true
							}
							if admin.User.ID == h.bot.Self.ID {
								isAdmin = true
							}
							if isAdmin && isCreator {
								break;
							}
						}
						if !isCreator {
							reply = tgbotapi.NewMessage(message.Chat.ID, "只有频道创建者才能注册频道")
						} else if !isAdmin {
							reply = tgbotapi.NewMessage(message.Chat.ID, "请把bot设置为频道管理员")
						} else {
							session.Status = models.SESSION_FINISHED
							channel := &models.RegisteredChannelForCreator{ChatID: message.ForwardFromChat.ID, CreatorID: message.From.ID, Settings: &models.ChannelSettingsForCreator{Mode: models.MODE_AUTO}}
							if h.repo.UpdateSession(session) && h.channelRepo.InsertRegisterChannelForCreator(channel) {
								reply = tgbotapi.NewMessage(message.Chat.ID, "注册成功")
							} else {
								reply = tgbotapi.NewMessage(message.Chat.ID, "未知错误,请重试")
							}
						}
					} else {
						reply = tgbotapi.NewMessage(message.Chat.ID, "请检查是否把bot添加到频道")
					}
				}
				h.bot.Send(reply)
			}
		}*/
}

func (h *Handler) handleRegisterRequest(message *tgbotapi.Message) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	reply := tgbotapi.NewMessage(message.Chat.ID, "开发中，暂不支持注册")
	h.bot.Send(reply)
	/*	session, found := h.repo.SelectSession(message.Chat.ID)
		if !found {
			session = models.NewSession(message.Chat)
			session.Status = models.SessionWaitingChannel
			ok := h.repo.InsertSession(session)
			if ok {
				reply := tgbotapi.NewMessage(message.Chat.ID, "开发中，暂不支持注册")
				h.bot.Send(reply)
			}
		}*/
}

func (h *Handler) handleChannelMessage(message *tgbotapi.Message) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	channel, found := h.channelRepo.FindChannelForCreatorByChatID(message.Chat.ID)
	if found {
		if channel.Settings.Mode == models.ModeAuto {

			h.resolveCommentSection(message.Chat.ID, message.MessageID, EnableCommentSection)
			//data := string(message.MessageID)
			/*reply := &tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:              channel.ChatID,
					ReplyToMessageID:    message.MessageID,
					DisableNotification: true,
				},
				Text:      "暂无评论",
				ParseMode: tgbotapi.ModeMarkdown,
			}
			msg, err := h.bot.Send(reply)
			if err == nil {
				post := &models.Post{MessageID: message.MessageID,
					ChatID:    message.Chat.ID,
					SectionID: msg.MessageID,
				}
				param, err := utils.EncodeParam(post)
				if err != nil {
					log.Println(err)
					return
				}
				url := "http://t.me/" + h.bot.Self.UserName + "?start=" + param
				keyboard := utils.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.InlineKeyboardButton{Text: "查看详情️", URL: &url}))
				_, err = h.bot.Send(tgbotapi.NewEditMessageReplyMarkup(msg.Chat.ID, msg.MessageID, keyboard))
				h.commentRepo.InsertCommentSection(post)
			}*/
		}
	}
}

func NewHandler(config configs.BotConfig) *Handler {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	ret, err := bot.SetWebhook(tgbotapi.NewWebhook(config.Server + "/" + bot.Token))
	if err != nil {
		panic(err)
	}
	log.Println(ret)
	return &Handler{
		bot: bot,
		//repo:        repositories.NewRepository(),
		channelRepo: repositories.NewChannelRepository(),
		sessionRepo: repositories.NewSessionRepository(),
		commentRepo: repositories.NewCommentRepository()}
}
