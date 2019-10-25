package types

// TokenInfo TokenInfo
type TokenInfo struct {
	Address string
	Symbol  string
	Status  string
}

// TokenMetaInfo TokenMetaInfo
type TokenMetaInfo struct {
	Address     string `json:"address"`
	Name        string `json:"name"`
	TotalSupply string `json:"totalSupply"`
	Decimals    string `json:"decimals"`
}

// TokenTxnInfo TokenTxnInfo
type TokenTxnInfo struct {
	BlockNumber  int64  `json:"blockNumber"`
	BlockHash    string `json:"blockHash"`
	Timestamp    string `json:"timestamp"`
	TransferHash string `json:"transferHash"`
	Sender       string `json:"sender"`
	Receiver     string `json:"receiver"`
	Value        string `json:"value"`
}

// TokenHolderInfo TokenHolderInfo
type TokenHolderInfo struct {
	Balance uint64 `json:"balance"`
	Address string `json:"address"`
}

//Block Block
type Block struct {
	Difficulty      string        `json:"difficulty"`
	GasLimit        string        `json:"gasLimit"`
	GasUsed         string        `json:"gasUsed"`
	Hash            string        `json:"hash"`
	Miner           string        `json:"miner"`
	Nonce           string        `json:"nonce"`
	Number          string        `json:"number"`
	ParentHash      string        `json:"parentHash"`
	Size            string        `json:"size"`
	Timestamp       string        `json:"timestamp"`
	TotalDifficulty string        `json:"totalDifficulty"`
	Transactions    []Transaction `json:"transactions"`
}

//Transaction Transaction
type Transaction struct {
	TransactionHash  string `json:"hash"`
	Input            string `json:"input"`
	From             string `json:"from"`
	To               string `json:"to"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasprice"`
	Nonce            string `json:"nonce"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
}

//Receipt Receipt
type Receipt struct {
	TransactionHash  string `json:"transactionHash"`
	TransactionIndex string `json:"transactionIndex"`
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	LogsBloom        string `json:"logsBloom"`
	Logs             []Log  `json:"logs"`
	FromAddress      string `json:"from"`
	ToAddress        string `json:"to"`
	GasUsed          string `json:"gasUsed"`
}

//Log Log
type Log struct {
	Address          string   `json:"address"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
	Data             string   `json:"data"`
	Topics           []string `json:"topics"`
	BlockNumber      string   `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}

//BlockNumberResult BlockNumberResult
type BlockNumberResult struct {
	ID      int    `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
}

//BalanceResult BalanceResult
type BalanceResult struct {
	ID      int    `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
}

//ReceiptResult ReceiptResult
type ReceiptResult struct {
	ID      int     `json:"id"`
	JSONRPC string  `json:"jsonrpc"`
	Result  Receipt `json:"result"`
}

//BlcokResult BlcokResult
type BlcokResult struct {
	ID      int    `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Result  Block  `json:"result"`
}

//LogsResult LogsResult
type LogsResult struct {
	ID      int    `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Result  []Log  `json:"result"`
}
