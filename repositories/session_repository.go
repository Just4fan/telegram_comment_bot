package repositories

import "telegram_comment_bot/models"

type SessionRepository struct {
	sessions map[int64]interface{}
	database *Database
}

func (r *SessionRepository) SelectCommentSessionByChatID(chatID int64) (session *models.CommentSession, found bool) {
	v := r.sessions[chatID]
	if v != nil {
		session, found = v.(*models.CommentSession)
		return
	}
	return nil, false
}

func (r *SessionRepository) InsertCommentSession(session *models.CommentSession) (ok bool) {
	r.sessions[session.ChatID] = session
	ok = true
	return
}

func (r *SessionRepository) UpdateCommentSession(session *models.CommentSession) (ok bool) {
	r.sessions[session.ChatID] = session
	ok = true
	return
}

func (r *SessionRepository) SelectSessionByChatID(chatID int64) (v interface{}, found bool) {
	v = r.sessions[chatID]
	if v != nil {
		found = true
	} else {
		found = false
	}
	return
}

func (r *SessionRepository) InsertSession(chatID int64, v interface{}) (ok bool) {
	r.sessions[chatID] = v
	ok = true
	return
}

func (r *SessionRepository) RemoveSession(chatID int64) (ok bool) {
	r.sessions[chatID] = nil
	ok = true
	return
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{database: GetDatabase("bot"), sessions: make(map[int64]interface{})}
}
