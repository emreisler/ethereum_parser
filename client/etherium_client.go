package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/emreisler/ethereum_parser/domain"
	"io/ioutil"
	"net/http"
)

type EthereumClient interface {
	GetCurrentBlockNumber() (int, error)
	GetTxHashes(blockNumber int) ([]string, error)
	GetTxObjects(blockNumber int) ([]*domain.Transaction, error)
	GetTxByHash(hash string) (*domain.Transaction, error)
}

type ethereumClient struct {
	url string
}

func NewEthereumClient(url string) EthereumClient {
	return &ethereumClient{
		url: url,
	}
}

func (e *ethereumClient) GetCurrentBlockNumber() (int, error) {
	response := e.call("eth_blockNumber", nil)
	var result string
	if err := json.Unmarshal(response["result"], &result); err != nil {
		return 0, err
	}
	blockNumber := hexToInt(result)
	return blockNumber, nil
}

func (e *ethereumClient) GetTxHashes(blockNumber int) ([]string, error) { //21372053

	params := []interface{}{fmt.Sprintf("0x%x", blockNumber), false}
	response := e.call("eth_getBlockByNumber", params)
	var block struct {
		Transactions []string `json:"transactions"`
	}

	if err := json.Unmarshal(response["result"], &block); err != nil {
		return nil, err
	}

	return block.Transactions, nil
}

func (e *ethereumClient) GetTxObjects(blockNumber int) ([]*domain.Transaction, error) {
	params := []interface{}{fmt.Sprintf("0x%x", blockNumber), true}
	response := e.call("eth_getBlockByNumber", params)

	var block struct {
		Transactions []*domain.Transaction `json:"transactions"`
	}
	if err := json.Unmarshal(response["result"], &block); err != nil {
		return nil, err
	}
	return block.Transactions, nil
}

func (e *ethereumClient) GetTxByHash(hash string) (*domain.Transaction, error) {
	params := []interface{}{hash} // JSON-RPC requires the transaction hash as a parameter
	response := e.call("eth_getTransactionByHash", params)

	if response == nil {
		return nil, fmt.Errorf("failed to fetch transaction details")
	}

	var tx domain.Transaction
	if err := json.Unmarshal(response["result"], &tx); err != nil {
		return nil, fmt.Errorf("error parsing transaction: %v", err)
	}

	// If the transaction result is empty, it may not exist or be pending
	if tx.Hash == "" {
		return nil, fmt.Errorf("transaction not found: %s", hash)
	}

	return &tx, nil
}

// Helper method to make JSON-RPC calls
func (e *ethereumClient) call(method string, params []interface{}) map[string]json.RawMessage {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	}
	payloadBytes, _ := json.Marshal(payload)

	resp, err := http.Post(e.url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error making JSON-RPC call:", err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var response map[string]json.RawMessage
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error parsing JSON-RPC response:", err)
		return nil
	}

	return response
}

func hexToInt(hexStr string) int {
	var i int
	fmt.Sscanf(hexStr, "0x%x", &i)
	return i
}

func intToHex(i int) string {
	return fmt.Sprintf("%x", i)
}
