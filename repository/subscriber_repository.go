package repository

import (
	"errors"
	"github.com/emreisler/ethereum_parser/domain"
	"sync"
)

type subscriberRepo struct {
	subscribers map[string]*domain.Subscriber
	mu          sync.RWMutex
}

func NewInMemorySubscriberRepo() SubscriberRepository {
	return &subscriberRepo{
		subscribers: make(map[string]*domain.Subscriber),
	}
}

func (s *subscriberRepo) GetSubscribers() []string {
	var subscribersAddresses []string
	for _, subscriber := range s.subscribers {
		subscribersAddresses = append(subscribersAddresses, subscriber.Address)
	}
	return subscribersAddresses
}

func (s *subscriberRepo) SubscriberExists(address string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.subscribers[address]
	return ok
}

func (s *subscriberRepo) AddTxHash(address, hash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	subscriber, ok := s.subscribers[address]
	if !ok {
		return errors.New("subscriber not found")
	}
	if _, ok := subscriber.TxHashes[hash]; !ok {
		subscriber.TxHashes[hash] = struct{}{}
	}
	return nil
}

func (s *subscriberRepo) AddSubscriber(address string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subscribers[address] = &domain.Subscriber{
		TxHashes: make(map[string]struct{}),
		Address:  address,
	}
	return true
}

func (s *subscriberRepo) GetTxHashes(address string) map[string]struct{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.subscribers[address].TxHashes
}
