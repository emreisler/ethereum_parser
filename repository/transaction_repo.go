package repository

import (
	"github.com/emreisler/ethereum_parser/domain"
	"sync"
)

type InMemoryTxRepo struct {
	txs map[string]*domain.Transaction
	mu  sync.RWMutex
}

func NewInMemoryTxRepo() TransactionRepository {
	return &InMemoryTxRepo{
		txs: make(map[string]*domain.Transaction),
	}
}

func (tr *InMemoryTxRepo) AddTx(tx *domain.Transaction) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.txs[tx.Hash] = tx
	return nil
}

func (tr *InMemoryTxRepo) GetTx(hash string) (*domain.Transaction, error) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	return tr.txs[hash], nil
}

func (tr *InMemoryTxRepo) TxExist(hash string) bool {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	_, ok := tr.txs[hash]
	return ok
}
