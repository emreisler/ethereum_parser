package domain

type Subscriber struct {
	Address  string
	TxHashes map[string]struct{}
}
