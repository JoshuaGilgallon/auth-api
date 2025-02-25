package services

import (
	"auth-api/internal/models"
	"sync"
	"time"
)

type SessionStore struct {
	sessions map[string]models.AdminSession
	mutex    sync.RWMutex
}

var (
	store     *SessionStore
	onceStore sync.Once
)

func GetSessionStore() *SessionStore {
	onceStore.Do(func() {
		store = &SessionStore{
			sessions: make(map[string]models.AdminSession),
		}
		// Start cleanup routine
		go store.cleanupRoutine()
	})
	return store
}

func (s *SessionStore) Set(token string, session models.AdminSession) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sessions[token] = session
}

func (s *SessionStore) Get(token string) (models.AdminSession, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	session, exists := s.sessions[token]
	return session, exists
}

func (s *SessionStore) Delete(token string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.sessions, token)
}

func (s *SessionStore) cleanupRoutine() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		s.mutex.Lock()
		now := time.Now()
		for token, session := range s.sessions {
			if now.After(session.AccessExpiresAt) {
				delete(s.sessions, token)
			}
		}
		s.mutex.Unlock()
	}
}
