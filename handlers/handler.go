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
)

type Handler struct {
	bot         *tgbotapi.BotAPI
	repo        *repositories.Repository
	channelRepo *repositories.ChannelRepository
	sessionRepo *repositories.SessionRepository
	commentRepo *repositories.CommentRepository
}

func (h *Handler) Start(offset int) {
	updates := h.bot.ListenForWebhook("/" + h.bot.Token)
	go http.ListenAndServe("0.0.0.0:80", nil)

	for update := range updates {
		/*		config := tgbotapi.NewUpdate(offset)
				config.Timeout = 60
				updates, err := h.bot.GetUpdatesChan(config)
				if err != nil {
					continue
				}
				for update := range updates {*/

		log.Printf("%+v\n", update)
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

func (h *Handler) handleCancelRequest(message *tgbotapi.Message) {
	h.sessionRepo.RemoveSession(message.Chat.ID)
	h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "会话状态已清空"))
}

func (h *Handler) handleReloadRequest(message *tgbotapi.Message) {
	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: message.Chat.ID,
		},
		Text: `⚠️此功能将尝试修复评论功能异常的帖子⚠️
				转发评论功能异常的频道消息到此或者输入指令 /cancel 取消操作`,
	}
	h.bot.Send(msg)
	session := &models.ReloadSession{BaseSession: models.BaseSession{
		ChatID: message.Chat.ID,
		Status: models.SessionWaitingPost,
	}}
	h.sessionRepo.InsertSession(session.ChatID, session)
}

func (h *Handler) handleStartRequest(message *tgbotapi.Message) {
	log.Println(message.CommandArguments())
	if len(message.CommandArguments()) != 0 {
		h.handleCommentRequest(message)
	}
}

func (h *Handler) handleCommentRequest(message *tgbotapi.Message) {
	target, err := utils.DecodeParam(message.CommandArguments())
	if err != nil {
		log.Println(err)
		return
	}
	_, found := h.channelRepo.FindChannelForCreatorByChatID(target.ChatID)
	if !found {
		h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "频道未注册！"))
	}
	comments := h.commentRepo.FindCommentsByMessage(target.ChatID, target.MessageID, 1, 5)
	//var msg tgbotapi.Message
	forward := tgbotapi.NewForward(message.Chat.ID, target.ChatID, target.MessageID)
	forwardMsg, err := h.bot.Send(forward)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("channel:%+v\n", *forwardMsg.ForwardFromChat)
	reply := &tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:              message.Chat.ID,
			ReplyToMessageID:    forwardMsg.MessageID,
			DisableNotification: false,
		},
		Text:      utils.CommentsText(comments),
		ParseMode: tgbotapi.ModeHTML,
	}
	msg, err := h.bot.Send(reply)
	if err != nil {
		log.Println(err)
		return
	}

	keyboard := utils.NewInlineKeyboardMarkup(
		utils.CommentPostKeyBoardRow(target.ChatID, target.MessageID, msg.MessageID),
		utils.CommentIndexKeyBoardRow(comments, msg.MessageID),
		utils.PageSwitchKeyboardRow(1, target.ChatID, target.MessageID, msg.MessageID))
	edit := tgbotapi.NewEditMessageReplyMarkup(msg.Chat.ID, msg.MessageID, keyboard)
	_, err = h.bot.Send(edit)
	if err != nil {
		log.Println(err)
		return
	}
	/*	session := &models.CommentSession{
			BaseSession: models.BaseSession{
				ChatID: message.Chat.ID,
				Status: models.SESSION_WAITING_ADD_COMMENT,
			},
			Post: &models.Post{
				ChatID: target.ChatID,
				MessageID: target.MessageID,
			},
			Comments: comments,
			Area: msgID,
			Panel: msg.MessageID,
			Page: 1,
			Total: total,
			Params: message.CommandArguments(),
		}
		h.sessionRepo.InsertCommentSession(session)*/
}

func (h *Handler) handleCallbackQuery(query *tgbotapi.CallbackQuery) {
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

		session := &models.CommentSession{
			BaseSession: models.BaseSession{
				ChatID: int64(query.From.ID),
				Status: models.SessionWaitingInputComment,
			},
		}

		if arr[0] == "cp" {
			var err error
			chatID, messageID, panelID, err := utils.DecodeData(arr[1])
			if err != nil {
				log.Println(err)
				return
			}
			post, found := h.commentRepo.FindPostByMessage(chatID, messageID)
			if !found {
				log.Println("Post Not Found")
				return
			}
			session.Post = post
			session.Panel = panelID
		} else {
			id, panelID, err := utils.DecodeDataR(arr[1])
			if err != nil {
				log.Println(err)
				return
			}
			comment, found := h.commentRepo.FindCommentByID(id)
			if !found {
				log.Println("Comment Not Found")
				return
			}
			session.Comment = comment
			session.Panel = panelID
		}

		h.sessionRepo.InsertSession(session.ChatID, session)
		text := ""

		if session.Comment != nil {
			text += "评论用户 <b>" + session.Comment.UserName + "</b>:\n<pre>\n</pre>" + "<pre>      </pre><i>" + session.Comment.Content + "</i>\n<pre>\n</pre>"
		}

		text += "输入评论并发送，发送 /cancel 指令退出评论:"

		msg := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: session.ChatID,
			},
			Text:      text,
			ParseMode: tgbotapi.ModeHTML,
		}

		_, err := h.bot.Send(msg)

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

			comments := h.commentRepo.FindCommentsByMessage(chatID, messageID, page, 5)
			if len(comments) > 0 {

				text := utils.CommentsText(comments)

				markup := utils.NewInlineKeyboardMarkup(
					utils.CommentPostKeyBoardRow(chatID, messageID, panelID),
					utils.CommentIndexKeyBoardRow(comments, panelID),
					utils.PageSwitchKeyboardRow(page, chatID, messageID, panelID))

				edit := tgbotapi.EditMessageTextConfig{
					BaseEdit: tgbotapi.BaseEdit{
						ChatID:      int64(query.From.ID),
						MessageID:   panelID,
						ReplyMarkup: &markup,
					},
					Text:      text,
					ParseMode: tgbotapi.ModeHTML,
				}

				log.Printf("%+v\n", edit)

				h.bot.Send(edit)

				h.bot.AnswerCallbackQuery(tgbotapi.NewCallback(query.ID, ""))

				return

			}
		}

		h.bot.AnswerCallbackQuery(tgbotapi.NewCallback(query.ID, "所请求页面没有更多数据"))
	}
}

func (h *Handler) handleMessage(message *tgbotapi.Message) {
	v, found := h.sessionRepo.SelectSessionByChatID(message.Chat.ID)
	if found {
		typ := reflect.TypeOf(v)
		log.Println(typ)
		switch typ {

		case models.TypeReloadSession:

			session, _ := v.(*models.ReloadSession)
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
				post, found := h.commentRepo.FindPostByMessage(message.ForwardFromChat.ID, message.ForwardFromMessageID)
				if !found {
					h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "该频道消息没有开启评论功能"))
					return
				}
				comments := h.commentRepo.FindCommentsByMessage(post.ChatID, post.MessageID, 1, 5)
				commentsText := utils.CommentsText(comments)
				params, err := utils.EncodeParam(post)
				if err != nil {
					log.Println(err)
					return
				}
				url := "http://t.me/comment_it_bot?start=" + params
				keyboard := utils.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.InlineKeyboardButton{Text: "查看详情️", URL: &url}))
				edit := tgbotapi.EditMessageTextConfig{
					BaseEdit: tgbotapi.BaseEdit{
						ChatID:      post.ChatID,
						MessageID:   post.AreaID,
						ReplyMarkup: &keyboard,
					},
					Text:      commentsText,
					ParseMode: tgbotapi.ModeHTML,
				}
				log.Printf("%+v\n", edit)
				_, err = h.bot.Send(edit)
				if err != nil {
					log.Println(err)
					h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "评论功能已重新加载"))
					h.sessionRepo.RemoveSession(message.Chat.ID)
				}
			}

		case models.TypeCommentSession:

			session, _ := v.(*models.CommentSession)
			log.Printf("%+v\n", session)
			switch session.Status {
			case models.SessionWaitingInputComment:
				reply := &models.Comment{
					ID:       primitive.NewObjectID(),
					Post:     session.Post,
					UserID:   message.From.ID,
					UserName: message.From.FirstName + " " + message.From.LastName,
					Date:     time.Now().Unix(),
					Content:  message.Text,
				}
				if session.Comment != nil {
					log.Printf("%+v\n", session.Comment)
					//content := utils.CommentText(message.From.UserName, comment.UserName, message.Text)
					reply.Post = session.Comment.Post
					reply.ReplyTo = session.Comment.UserName
					reply.Reply = session.Comment.ID
					reply.ReplyID = session.Comment.UserID
				}
				log.Printf("%+v", reply)
				ok := h.commentRepo.InsertComment(reply)
				if !ok {
					return
				}
				comments := h.commentRepo.FindCommentsByMessage(reply.Post.ChatID, reply.Post.MessageID, 1, 5)
				commentsText := utils.CommentsText(comments)
				params, err := utils.EncodeParam(reply.Post)
				if err != nil {
					log.Println(err)
					return
				}
				url := "http://t.me/comment_it_bot?start=" + params
				keyboard := utils.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.InlineKeyboardButton{Text: "查看详情️", URL: &url}))
				edit := tgbotapi.EditMessageTextConfig{
					BaseEdit: tgbotapi.BaseEdit{
						ChatID:      reply.Post.ChatID,
						MessageID:   reply.Post.AreaID,
						ReplyMarkup: &keyboard,
					},
					Text:      commentsText,
					ParseMode: tgbotapi.ModeHTML,
				}
				log.Printf("%+v\n", edit)
				h.bot.Send(edit)
				h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "评论成功"))
				if session.Comment != nil {
					markup := utils.NewInlineKeyboardMarkup(utils.QuickReplyKeyBoardRow(reply, 0))
					notification := tgbotapi.MessageConfig{
						BaseChat: tgbotapi.BaseChat{
							ChatID:      int64(reply.ReplyID),
							ReplyMarkup: markup,
						},
						Text:      utils.NotificationText(reply.UserName, session.Comment.Content, reply.Content),
						ParseMode: tgbotapi.ModeHTML,
					}
					_, err := h.bot.Send(notification)
					if err != nil {
						log.Println(err)
						h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "评论成功，但是对方可能关闭了跟我的对话，通知无法送达"))
					}
				}
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

	channel, found := h.channelRepo.FindChannelForCreatorByChatID(message.Chat.ID)
	if found {
		if channel.Settings.Mode == models.ModeAuto {
			//data := string(message.MessageID)
			reply := &tgbotapi.MessageConfig{
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
					ChatID: message.Chat.ID,
					AreaID: msg.MessageID,
				}
				param, err := utils.EncodeParam(post)
				if err != nil {
					log.Println(err)
					return
				}
				url := "http://t.me/comment_it_bot?start=" + param
				keyboard := utils.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.InlineKeyboardButton{Text: "查看详情️", URL: &url}))
				_, err = h.bot.Send(tgbotapi.NewEditMessageReplyMarkup(msg.Chat.ID, msg.MessageID, keyboard))
				h.commentRepo.InsertPost(post)
			}
		}
	}
}

func NewHandler(config configs.BotConfig) *Handler {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	ret, err := bot.SetWebhook(tgbotapi.NewWebhook(config.Server + bot.Token))
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
