package repository

import (
	"github.com/emreisler/ethereum_parser/domain"
	"sync"
)

type updatedSubscriberRepo struct {
	txMap map[string]map[string]*domain.Transaction
	mu    sync.RWMutex
}

func NewUpdatedSubscriberRepo() UpdatedSubscriberRepo {
	return &updatedSubscriberRepo{
		txMap: make(map[string]map[string]*domain.Transaction),
	}
}

func (u *updatedSubscriberRepo) AddSubscriber(address string) bool {
	u.mu.Lock()
	defer u.mu.Unlock()

	if _, ok := u.txMap[address]; !ok {
		u.txMap[address] = make(map[string]*domain.Transaction)
	}
	return u.txMap[address][address] != nil
}

func (u *updatedSubscriberRepo) SubscriberExists(address string) bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	_, ok := u.txMap[address]
	return ok
}

func (u *updatedSubscriberRepo) AddTx(address string, tx *domain.Transaction) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if _, ok := u.txMap[address]; !ok {
		u.txMap[address] = make(map[string]*domain.Transaction)
	}

	u.txMap[address][tx.Hash] = tx
}

func (u *updatedSubscriberRepo) GetTxHashes(address string) []domain.Transaction {
	u.mu.Lock()
	defer u.mu.Unlock()
	var txHashes []domain.Transaction
	for _, tx := range u.txMap[address] {
		txHashes = append(txHashes, *tx)
	}
	return txHashes
}
