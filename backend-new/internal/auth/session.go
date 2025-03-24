package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type Session struct {
	Username string
	Expires  time.Time
}

type Store struct {
	sync.RWMutex
	data map[string]Session
}

func NewStore() *Store {
	return &Store{data: make(map[string]Session)}
}

func (s *Store) Create(user string, ttl time.Duration) string {
	b := make([]byte, 32)
	rand.Read(b)
	token := hex.EncodeToString(b)
	s.Lock()
	s.data[token] = Session{user, time.Now().Add(ttl)}
	s.Unlock()
	return token
}

func (s *Store) Get(token string) (Session, bool) {
	s.RLock()
	sess, ok := s.data[token]
	s.RUnlock()
	if !ok || time.Now().After(sess.Expires) {
		return Session{}, false
	}
	return sess, true
}

func (s *Store) Delete(token string) {
	s.Lock()
	delete(s.data, token)
	s.Unlock()
}
