package storage

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"

	_ "github.com/lib/pq"

	"github.com/google/go-safeweb/safesql"
	"golang.org/x/crypto/scrypt"

	"github.com/aiit2022-pbl-okuhara/incident-training/config"
)

type Note struct {
	Title, Text string
}

type DB struct {
	conn safesql.DB
	mu   sync.Mutex
	// user -> note title -> notes
	notes map[string]map[string]Note
	// user -> token
	sessionTokens map[string]string
	// token -> user
	userSessions map[string]string
	// user -> pw hash
	credentials map[string]string
}

func NewDB() (*DB, error) {
	c := config.Config

	conn, err := safesql.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBPort, c.DBUsername, c.DBPassword, c.DBName),
	)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &DB{
		conn:          conn,
		notes:         map[string]map[string]Note{},
		sessionTokens: map[string]string{},
		userSessions:  map[string]string{},
		credentials:   map[string]string{},
	}, nil
}

// Notes

func (s *DB) AddOrEditNote(user string, n Note) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.notes[user] == nil {
		s.notes[user] = map[string]Note{}
	}
	s.notes[user][n.Title] = n
}

func (s *DB) GetNotes(user string) []Note {
	s.mu.Lock()
	defer s.mu.Unlock()
	var ns []Note
	for _, n := range s.notes[user] {
		ns = append(ns, n)
	}
	return ns
}

// Sessions

func (s *DB) GetUser(token string) (user string, valid bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, valid = s.sessionTokens[token]
	return user, valid
}

func (s *DB) GetToken(user string) (token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	token, has := s.userSessions[user]
	if has {
		return token
	}
	token = genToken()
	s.userSessions[user] = token
	s.sessionTokens[token] = user
	return token
}

func (s *DB) DelSession(user string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	token, has := s.userSessions[user]
	if !has {
		return
	}
	delete(s.userSessions, user)
	delete(s.sessionTokens, token)
}

func genToken() string {
	b := make([]byte, 20)
	rand.Read(b)
	tok := base64.RawStdEncoding.EncodeToString(b)
	return tok
}

// Credentials

// HasUser checks if the user exists.
func (s *DB) HasUser(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, has := s.credentials[name]
	return has
}

// AddUser adds a user to the storage if it is not already there.
func (s *DB) AddOrAuthUser(name, password string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if password == "" {
		return errors.New("password cannot be empty")
	}
	if storedHash, has := s.credentials[name]; has {
		if storedHash != hash(password) {
			return errors.New("wrong password")
		}
		return nil
	}
	s.credentials[name] = hash(password)
	return nil
}

func hash(pw string) string {
	salt := []byte("please use a proper salt in production")
	hash, err := scrypt.Key([]byte(pw), salt, 32768, 8, 1, 32)
	if err != nil {
		// TODO: 適切に error 処理を行う
		panic("this should not happen")
	}
	return string(hash)
}
