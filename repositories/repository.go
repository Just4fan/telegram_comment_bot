package repositories

import "comment_bot/models"

type Repository struct {
	sessions map[int64]*models.Session
}

func (r *Repository) SelectSession(id int64) (session *models.Session, found bool) {
	session = r.sessions[id]
	if session != nil {
		found = true
	} else {
		found = false
	}
	return
}

func (r *Repository) InsertSession(session *models.Session) (ok bool) {
	r.sessions[session.Chat.ID] = session
	ok = true
	return
}

func (r *Repository) UpdateSession(session *models.Session) (ok bool) {
	r.sessions[session.Chat.ID] = session
	ok = true
	return
}

func NewRepository() *Repository {
	return &Repository{sessions: make(map[int64]*models.Session)}
}
