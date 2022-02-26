package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

var Sm *SessionManager

type SessionManager struct {
	sessions     map[string]*SessionData
	sessionKey   string
	lifetimeDate int
}

func NewSessionManager(ctx context.Context, sessionKey string, lifetimeDate int) *SessionManager {
	if Sm != nil {
		return Sm
	}
	Sm := &SessionManager{sessions: map[string]*SessionData{}, sessionKey: sessionKey, lifetimeDate: lifetimeDate}
	Sm.StartDeleteSessionAsync(ctx)
	return Sm
}

func (m *SessionManager) StartDeleteSessionAsync(ctx context.Context) {
	timer := time.NewTicker(time.Hour * 24 * time.Duration(m.lifetimeDate))

	go func() {
		defer timer.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Printf("delete session process is canceled")
				return
			case <-timer.C:
				m.endExpiredSession()
			}
		}
	}()
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
	cookie := http.Cookie{
		Name:     m.sessionKey,
		Value:    sessionId,
		Expires:  time.Now().AddDate(0, 0, m.lifetimeDate),
		HttpOnly: true}
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

// 期限切れのセッションを削除
func (m *SessionManager) endExpiredSession() {
	now := time.Now().Unix()
	for k, v := range m.sessions {
		if v.ExpiredAt.Unix() < now {
			log.Printf("delete session key: %s", k)
			delete(m.sessions, k)
		}
	}
}
