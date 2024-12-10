package repository

import (
	"github.com/emreisler/ethereum_parser/domain"
	"sync"
)

type txRepo struct {
	txs map[string]*domain.Transaction
	mu  sync.RWMutex
}

func NewInMemoryTxRepo() TransactionRepository {
	return &txRepo{
		txs: make(map[string]*domain.Transaction),
	}
}

func (tr *txRepo) AddTx(tx *domain.Transaction) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.txs[tx.Hash] = tx
	return nil
}

func (tr *txRepo) GetTx(hash string) (*domain.Transaction, error) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	return tr.txs[hash], nil
}

func (tr *txRepo) TxExist(hash string) bool {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	_, ok := tr.txs[hash]
	return ok
}
