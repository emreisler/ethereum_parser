
# Ethereum Parser HTTP Server

A simple HTTP server for interacting with an Ethereum parser. This server allows you to:
- Get the current block number.
- Subscribe to an Ethereum address.
- Retrieve transactions for a subscribed Ethereum address.

## Endpoints

### 1. Get Current Block
- **Endpoint**: `GET /current-block`
- **Description**: Retrieves the latest block number on the Ethereum blockchain.
- **Example Curl Command**:
  ```bash
  curl http://localhost:8080/current-block
  ```
- **Sample Response**:
  ```json
  {
    "current_block": 12345678
  }
  ```

---

### 2. Subscribe to an Address
- **Endpoint**: `POST /subscribe`
- **Description**: Subscribes to an Ethereum address to track its transactions.
- **Request Body**:
  ```json
  {
    "address": "0x68d3a973e7272eb388022a5c6518d9b2a2e66fbf"
  }
  ```
- **Example Curl Command**:
  ```bash
  curl -X POST http://localhost:8080/subscribe \
       -H "Content-Type: application/json" \
       -d '{"address": "0x68d3a973e7272eb388022a5c6518d9b2a2e66fbf"}'
  ```
- **Sample Response**:
  ```json
  {
    "address": "0x68d3a973e7272eb388022a5c6518d9b2a2e66fbf",
    "isSubscribed": true
  }
  ```

---

### 3. Get Transactions for an Address
- **Endpoint**: `GET /transactions/:address`
- **Description**: Retrieves transactions for a specific Ethereum address that has been subscribed.
- **Example Curl Command**:
  ```bash
  curl http://localhost:8080/transactions/0x68d3a973e7272eb388022a5c6518d9b2a2e66fbf
  ```
- **Sample Response**:
  ```json
  {
    "address": "0x68d3a973e7272eb388022a5c6518d9b2a2e66fbf",
    "transactions": [
      {
        "hash": "0x123456789abcdef...",
        "from": "0xFromAddress...",
        "to": "0xToAddress...",
        "value": "1000000000000000000",
        "nonce": 1,
        "blockNumber": "12345678",
        "gas": "21000",
        "gasPrice": "50000000000",
        "input": "0x",
        "transactionIndex": "0"
      }
    ]
  }
  ```

---

## Running the Server
1. Start the server:
   ```bash
   go run main.go
   ```
2. The server will be available at `http://localhost:8080`.
