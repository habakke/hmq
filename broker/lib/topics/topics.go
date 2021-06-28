package topics

import (
	"fmt"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

const (
	// MWC is the multi-level wildcard
	MWC = "#"

	// SWC is the single level wildcard
	SWC = "+"

	// SEP is the topic level separator
	SEP = "/"

	// SYS is the starting character of the system level topics
	SYS = "$"

	// Both wildcards
	//lint:ignore U1000 This may be used later
	_WC = "#+"
)

var (
	providers = make(map[string]TopicsProvider)
)

type TopicsProvider interface {
	Subscribe(topic []byte, qos byte, subscriber interface{}) (byte, error)
	Unsubscribe(topic []byte, subscriber interface{}) error
	Subscribers(topic []byte, qos byte, subs *[]interface{}, qoss *[]byte) error
	Retain(msg *packets.PublishPacket) error
	Retained(topic []byte, msgs *[]*packets.PublishPacket) error
	Close() error
}

func Register(name string, provider TopicsProvider) {
	if provider == nil {
		panic("topics: Register provide is nil")
	}

	if _, dup := providers[name]; dup {
		panic("topics: Register called twice for provider " + name)
	}

	providers[name] = provider
}

func Unregister(name string) {
	delete(providers, name)
}

type Manager struct {
	p TopicsProvider
}

func NewManager(providerName string) (*Manager, error) {
	p, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provider %q", providerName)
	}

	return &Manager{p: p}, nil
}

func (m *Manager) Subscribe(topic []byte, qos byte, subscriber interface{}) (byte, error) {
	return m.p.Subscribe(topic, qos, subscriber)
}

func (m *Manager) Unsubscribe(topic []byte, subscriber interface{}) error {
	return m.p.Unsubscribe(topic, subscriber)
}

func (m *Manager) Subscribers(topic []byte, qos byte, subs *[]interface{}, qoss *[]byte) error {
	return m.p.Subscribers(topic, qos, subs, qoss)
}

func (m *Manager) Retain(msg *packets.PublishPacket) error {
	return m.p.Retain(msg)
}

func (m *Manager) Retained(topic []byte, msgs *[]*packets.PublishPacket) error {
	return m.p.Retained(topic, msgs)
}

func (m *Manager) Close() error {
	return m.p.Close()
}
