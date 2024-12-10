package repository

import "github.com/emreisler/ethereum_parser/domain"

type TransactionRepository interface {
	AddTx(tx *domain.Transaction) error
	GetTx(hash string) (*domain.Transaction, error)
	TxExist(hash string) bool
}

type SubscriberRepository interface {
	GetSubscribers() []string
	SubscriberExists(address string) bool
	AddSubscriber(address string) bool
	AddTxHash(address, hash string) error
	GetTxHashes(address string) map[string]struct{}
}
