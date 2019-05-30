package utils

import (
	"strconv"
	"telegram_comment_bot/models"
)

func commentText(user string, reply string, content string) string {
	/*	if len(reply) == 0 {
			return "<pre>" + user + ": " + content + "</pre>"
		}
		return "<pre>" + user + " 回复 " + reply + ": " + content + "</pre>"*/
	if len(reply) == 0 {
		return "<b>" + user + "</b>: " + content
	}
	return "<b>" + user + "</b> 回复 <b>" + reply + "</b>: " + content
}

func CommentsText(comments []*models.Comment) (text string) {
	if len(comments) == 0 {
		return "暂无评论"
	}
	for i, comment := range comments {
		text += strconv.Itoa(i+1) + ". " + commentText(comment.UserName, comment.ReplyTo, comment.Content) + "\n"
	}
	return
}

func NotificationText(user, comment, reply string) (text string) {
	text += "🔔 用户 <b>" + user + "</b>回复了评论：\n<pre>\n</pre>"
	text += "<pre>      </pre><i>" + comment + "</i>\n<pre>\n</pre>"
	text += reply
	return
}
