package main

import (
	"fmt"
	"github.com/emreisler/ethereum_parser/client"
	"github.com/emreisler/ethereum_parser/handlers"
	"github.com/emreisler/ethereum_parser/repository"
	"github.com/emreisler/ethereum_parser/usecases/parser"
	"log"
	"net/http"
)

func main() {
	// Initialize dependencies
	ethClient := client.NewEthereumClient("https://ethereum-rpc.publicnode.com")
	txRepo := repository.NewInMemoryTxRepo()
	subscriberRepo := repository.NewInMemorySubscriberRepo()
	parser := parser.NewEthereumParser(ethClient, txRepo, subscriberRepo)
	handler := handlers.NewParserHandler(parser)

	// Define routes
	http.HandleFunc("GET /current-block", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleGetCurrentBlock(w, r, parser)
	})

	http.HandleFunc("POST /subscribe", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleSubscribe(w, r, parser)
	})

	http.HandleFunc("GET /transactions/", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleGetTransactions(w, r, parser)
	})

	// Start the server
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
