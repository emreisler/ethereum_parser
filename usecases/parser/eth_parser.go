package parser

import (
	"fmt"
	"github.com/emreisler/ethereum_parser/client"
	"github.com/emreisler/ethereum_parser/domain"
	"github.com/emreisler/ethereum_parser/repository"
	"github.com/emreisler/ethereum_parser/usecases"
	"log/slog"
	"time"
)

type ethParser struct {
	ethClient      client.EthereumClient
	txRepo         repository.TransactionRepository
	subscriberRepo repository.SubscriberRepository
	currentBlock   int
}

func NewEthereumParser(ethClient client.EthereumClient, txRepo repository.TransactionRepository, subscriberRepo repository.SubscriberRepository) usecases.Parser {
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

// GetTransactions return retrieved transactions belongs to given address
// TODO It is not returning the transactions happened after last block update handling, so some new transactions are missing
// but the purpose is notification so 5 second delay may be acceptable ?
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

// check last parsed block and populate storage with missing transactions
func (p *ethParser) handleBlockUpdates(blocksChan <-chan int) {
	for newBlock := range blocksChan {
		if newBlock >= p.currentBlock {
			p.populateTxsWithBlockRange(p.currentBlock, newBlock)
			p.currentBlock = newBlock
			slog.Info(fmt.Sprintf("Current block updated to : %d", newBlock))
		}
	}
}

// populateTxsWithBlockRange loop through from the current block to the latest retrieved block from eth
func (p *ethParser) populateTxsWithBlockRange(startBlock, endBlock int) {
	for block := startBlock; block <= endBlock; block++ {
		p.populateTxs(block)
	}
}

// populateTxs if current block
// TODO name of the method does not say too much about the logic
func (p *ethParser) populateTxs(block int) error {
	if block == p.currentBlock {
		return p.populateExistingBlock(block)
	}
	return p.populateNewBlock(block)
}

// addToObserver check the retrieved transaction belongs to any subscriber, if it is, add hash to the subscriber
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

// populateNewBlock fetch all transaction with details
// it is a new block for our parser, so we need to fetch all transaction details
// TODO(1) why I keep all the transactions in repo ? Bec, if there is no subscriber at this time keeping the tx is non-sense
// or I should suppose to give subscriber txs from the beginning block ?
// design error
// I need to check if the transaction belongs to any subscriber, otherwise I should not keep it..
// Another scenario is giving all the history from starting block no on , but is this the requirement ?
func (p *ethParser) populateNewBlock(block int) error {
	txs, err := p.ethClient.GetTxObjects(block)
	if err != nil {
		return err
	}

	for _, tx := range txs {
		//I should have been checked if tx belongs to any subscriber..
		if !p.txRepo.TxExist(tx.Hash) {
			p.txRepo.AddTx(tx)
		}
		p.addToObserver(tx)
	}
	return nil
}

// populateExistingBlock fetch just transaction hashes and retrieve complete transaction if the hash does not exist in storage
// TODO starts with fetching just hashes bec, we are assuming we have most of the txs from this already scanned block
// to be able to reduce the payload from eth node
// looping through the hashes and fetch the tx details by using tx hash
// downside is what if we couldn't retrieve most of the transactions so there will be too many requests one by one
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

// watchLatestBlock watches latest block and return a channel of block number
// TODO I may check if the retrieved current block is higher than existing block here ?
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
