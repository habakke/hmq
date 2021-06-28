package sessions

import (
	"fmt"
	"sync"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

const (
	// Queue size for the ack queue
	//lint:ignore U1000 This may be used later
	defaultQueueSize = 16
)

type Session struct {

	// cmsg is the CONNECT message
	cmsg *packets.ConnectPacket

	// Will message to publish if connect is closed unexpectedly
	Will *packets.PublishPacket

	// Retained publish message
	Retained *packets.PublishPacket

	// topics stores all the topis for this session/client
	topics map[string]byte

	// Initialized?
	initted bool

	// Serialize access to this session
	mu sync.Mutex

	id string
}

func (s *Session) Init(msg *packets.ConnectPacket) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.initted {
		return fmt.Errorf("Session already initialized")
	}

	s.cmsg = msg

	if s.cmsg.WillFlag {
		s.Will = packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
		s.Will.Qos = s.cmsg.Qos
		s.Will.TopicName = s.cmsg.WillTopic
		s.Will.Payload = s.cmsg.WillMessage
		s.Will.Retain = s.cmsg.WillRetain
	}

	s.topics = make(map[string]byte, 1)

	s.id = string(msg.ClientIdentifier)

	s.initted = true

	return nil
}

func (s *Session) Update(msg *packets.ConnectPacket) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cmsg = msg
	return nil
}

func (s *Session) RetainMessage(msg *packets.PublishPacket) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Retained = msg

	return nil
}

func (s *Session) AddTopic(topic string, qos byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.initted {
		return fmt.Errorf("Session not yet initialized")
	}

	s.topics[topic] = qos

	return nil
}

func (s *Session) RemoveTopic(topic string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.initted {
		return fmt.Errorf("Session not yet initialized")
	}

	delete(s.topics, topic)

	return nil
}

func (s *Session) Topics() ([]string, []byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.initted {
		return nil, nil, fmt.Errorf("Session not yet initialized")
	}

	var (
		topics []string
		qoss   []byte
	)

	for k, v := range s.topics {
		topics = append(topics, k)
		qoss = append(qoss, v)
	}

	return topics, qoss, nil
}

func (s *Session) ID() string {
	return s.cmsg.ClientIdentifier
}

func (s *Session) WillFlag() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cmsg.WillFlag
}

func (s *Session) SetWillFlag(v bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cmsg.WillFlag = v
}

func (s *Session) CleanSession() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cmsg.CleanSession
}
