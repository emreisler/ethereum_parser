package handlers

import (
	"encoding/json"
	"github.com/emreisler/ethereum_parser/usecases"
	"log"
	"net/http"
)

type parserHandler struct {
	parser usecases.Parser
}

func NewParserHandler(parser usecases.Parser) *parserHandler {
	return &parserHandler{
		parser: parser,
	}
}

func (ph *parserHandler) HandleGetCurrentBlock(w http.ResponseWriter, r *http.Request, parser usecases.Parser) {

	block := ph.parser.GetCurrentBlock()
	response := map[string]interface{}{
		"current_block": block,
	}
	writeJSONResponse(w, http.StatusOK, response)
}

func (ph *parserHandler) HandleSubscribe(w http.ResponseWriter, r *http.Request, parser usecases.Parser) {

	var req struct {
		Address string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Address == "" {
		http.Error(w, "Invalid request, address is required", http.StatusBadRequest)
		return
	}

	isSubscribed := ph.parser.Subscribe(req.Address)
	response := map[string]interface{}{
		"address":      req.Address,
		"isSubscribed": isSubscribed,
	}
	writeJSONResponse(w, http.StatusOK, response)
}

func (ph *parserHandler) HandleGetTransactions(w http.ResponseWriter, r *http.Request, parser usecases.Parser) {

	// Extract the address from the URL path
	address := r.URL.Path[len("/transactions/"):]
	if address == "" {
		http.Error(w, "Address is required", http.StatusBadRequest)
		return
	}

	transactions := ph.parser.GetTransactions(address)
	if transactions == nil {
		http.Error(w, "No transactions found for the provided address", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"address":      address,
		"transactions": transactions,
	}
	writeJSONResponse(w, http.StatusOK, response)
}

// Helper function to write JSON responses
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
