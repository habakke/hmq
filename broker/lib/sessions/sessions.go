package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var (
	ErrSessionsProviderNotFound = errors.New("Session: Session provider not found")
	ErrKeyNotAvailable          = errors.New("Session: not item found for key")

	providers = make(map[string]SessionsProvider)
)

type SessionsProvider interface {
	New(id string) (*Session, error)
	Get(id string) (*Session, error)
	Del(id string)
	Save(id string) error
	Count() int
	Close() error
}

// Register makes a session provider available by the provided name.
// If a Register is called twice with the same name or if the driver is nil,
// it panics.
func Register(name string, provider SessionsProvider) {
	if provider == nil {
		panic("session: Register provide is nil")
	}

	if _, dup := providers[name]; dup {
		panic("session: Register called twice for provider " + name)
	}

	providers[name] = provider
}

func Unregister(name string) {
	delete(providers, name)
}

type Manager struct {
	p SessionsProvider
}

func NewManager(providerName string) (*Manager, error) {
	p, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provider %q", providerName)
	}

	return &Manager{p: p}, nil
}

func (m *Manager) New(id string) (*Session, error) {
	if id == "" {
		id = m.sessionId()
	}
	return m.p.New(id)
}

func (m *Manager) Get(id string) (*Session, error) {
	return m.p.Get(id)
}

func (m *Manager) Del(id string) {
	m.p.Del(id)
}

func (m *Manager) Save(id string) error {
	return m.p.Save(id)
}

func (m *Manager) Count() int {
	return m.p.Count()
}

func (m *Manager) Close() error {
	return m.p.Close()
}

func (m *Manager) sessionId() string {
	b := make([]byte, 15)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
