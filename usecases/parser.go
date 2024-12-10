package usecases

import (
	"fmt"
	"github.com/emreisler/ethereum_parser/client"
	"github.com/emreisler/ethereum_parser/domain"
	"github.com/emreisler/ethereum_parser/repository"
	"log/slog"
	"time"
)

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []domain.Transaction
}

type ethParser struct {
	ethClient      client.EthereumClient
	txRepo         repository.TransactionRepository
	subscriberRepo repository.SubscriberRepository
	currentBlock   int
}

func NewEthereumParser(ethClient client.EthereumClient, txRepo repository.TransactionRepository, subscriberRepo repository.SubscriberRepository) Parser {
	p := &ethParser{ethClient: ethClient, txRepo: txRepo, subscriberRepo: subscriberRepo}

	p.currentBlock = p.GetCurrentBlock()
	go p.handleBlockUpdates(p.watchLatestBlock())

	return p
}

func (p *ethParser) GetCurrentBlock() int {
	currentBlock, err := p.ethClient.GetCurrentBlockNumber()
	if err != nil {
		return 0
	}
	return currentBlock
}

func (p *ethParser) Subscribe(address string) bool {
	p.subscriberRepo.AddSubscriber(address)
	return true
}

func (p *ethParser) GetTransactions(address string) []domain.Transaction {
	if !p.subscriberRepo.SubscriberExists(address) {
		return nil
	}
	txHashes := p.subscriberRepo.GetTxHashes(address)

	var transactions []domain.Transaction
	for txHash := range txHashes {
		tx, _ := p.txRepo.GetTx(txHash)
		transactions = append(transactions, *tx)
	}

	return transactions
}

// check last parsed currentBlock and populate storage with missing transactions
func (p *ethParser) handleBlockUpdates(blocksChan <-chan int) {
	for newBlock := range blocksChan {
		if newBlock >= p.currentBlock {
			p.populateTxsWithBlockRange(p.currentBlock, newBlock)
			p.currentBlock = newBlock
			slog.Info(fmt.Sprintf("Current block updated to : %d", newBlock))
		}
	}
}

func (p *ethParser) populateTxsWithBlockRange(startBlock, endBlock int) {
	for block := startBlock; block <= endBlock; block++ {
		p.populateTxs(block)
	}
}

func (p *ethParser) populateTxs(block int) error {
	if block == p.currentBlock {
		return p.populateExistingBlock(block)
	}
	return p.populateNewBlock(block)

}

func (p *ethParser) addToObserver(tx *domain.Transaction) {
	subscriberAddresses := p.subscriberRepo.GetSubscribers()
	for _, subscriber := range subscriberAddresses {
		for _, address := range subscriberAddresses {
			if tx.From == address || tx.To == address {
				p.subscriberRepo.AddTxHash(subscriber, tx.Hash)
			}
		}
	}
}

func (p *ethParser) populateNewBlock(block int) error {
	txs, err := p.ethClient.GetTxObjects(block)
	if err != nil {
		return err
	}

	for _, tx := range txs {
		if !p.txRepo.TxExist(tx.Hash) {
			p.txRepo.AddTx(tx)
		}
		p.addToObserver(tx)
	}
	return nil
}

func (p *ethParser) populateExistingBlock(block int) error {
	txHashes, err := p.ethClient.GetTxHashes(block)
	if err != nil {
		return err
	}

	for _, txHash := range txHashes {
		if !p.txRepo.TxExist(txHash) {
			tx, err := p.ethClient.GetTxByHash(txHash)
			if err != nil {
				return err
			}
			p.txRepo.AddTx(tx)
			p.addToObserver(tx)
		}
	}
	return nil
}

func (p *ethParser) watchLatestBlock() <-chan int {
	ticker := time.NewTicker(5 * time.Second)
	blockChan := make(chan int)

	go func() {
		defer ticker.Stop()
		defer close(blockChan)
		for {
			select {
			case <-ticker.C:
				blockChan <- p.GetCurrentBlock()
			}
		}
	}()
	return blockChan
}
