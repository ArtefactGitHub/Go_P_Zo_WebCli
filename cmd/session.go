package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

type SessionManager struct {
	sessions     map[string]*SessionData
	sessionKey   string
	lifetimeDate int
}

func NewSessionManager(sessionKey string, lifetimeDate int) *SessionManager {
	return &SessionManager{sessions: map[string]*SessionData{}, sessionKey: sessionKey, lifetimeDate: lifetimeDate}
}

func (m *SessionManager) StartSession(w http.ResponseWriter, r *http.Request, data *SessionData) {
	// 既に存在する場合は終了
	if sessionId := m.getSessionIdByCookie(r); sessionId != "" {
		m.EndSession(w, r, sessionId)
	}

	sessionId := m.NewSessionId()
	data.SessionId = sessionId
	m.sessions[sessionId] = data

	m.setCookie(w, sessionId)
}

func (m *SessionManager) GetSession(w http.ResponseWriter, r *http.Request) (*SessionData, error) {
	sessionId := m.getSessionIdByCookie(r)

	// セッションが存在しない
	if session, ok := m.sessions[sessionId]; !ok {
		m.deleteCookie(w, r, sessionId)
		return nil, errors.New("session not found")
	} else {
		return session, nil
	}
}

func (m *SessionManager) EndSession(w http.ResponseWriter, r *http.Request, sessionId string) {
	delete(m.sessions, sessionId)
	m.deleteCookie(w, r, sessionId)
}

func (m *SessionManager) NewSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (m *SessionManager) NewSession() *SessionData {
	return &SessionData{}
}

func (m *SessionManager) setCookie(w http.ResponseWriter, sessionId string) {
	cookie := http.Cookie{Name: m.sessionKey, Value: sessionId, Expires: time.Now().AddDate(0, 0, m.lifetimeDate)}
	http.SetCookie(w, &cookie)
}

func (m *SessionManager) deleteCookie(w http.ResponseWriter, r *http.Request, sessionId string) {
	if c, err := r.Cookie(m.sessionKey); err == nil {
		c.MaxAge = -1
		http.SetCookie(w, c)
	}
}

func (m *SessionManager) getSessionIdByCookie(r *http.Request) string {
	c, err := r.Cookie(m.sessionKey)
	// クッキーが存在しない
	if err != nil {
		log.Printf("cookie not found")
		return ""
	}

	return c.Value
}
