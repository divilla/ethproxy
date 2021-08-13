package interfaces

type EthereumHttpClient interface {
	LatestBlockNumber() uint64
	GetLatestBlock() ([]byte, error)
	GetBlockByNumber(nr uint64) ([]byte, error)
}
