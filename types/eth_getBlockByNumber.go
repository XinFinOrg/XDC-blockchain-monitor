package types

import "encoding/json"

type EthBlockRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type EthBlockResponse struct {
	Jsonrpc string         `json:"jsonrpc"`
	ID      int            `json:"id"`
	Result  EthBlockResult `json:"result"`
}

type EthBlockResult struct {
	Difficulty       string                `json:"difficulty"`
	ExtraData        string                `json:"extraData"`
	GasLimit         string                `json:"gasLimit"`
	GasUsed          string                `json:"gasUsed"`
	Hash             string                `json:"hash"`
	LogsBloom        string                `json:"logsBloom"`
	Miner            string                `json:"miner"`
	MixHash          string                `json:"mixHash"`
	Nonce            string                `json:"nonce"`
	Number           string                `json:"number"`
	ParentHash       string                `json:"parentHash"`
	Penalties        interface{}           `json:"penalties"` // mainchain and subnet has different format
	ReceiptsRoot     string                `json:"receiptsRoot"`
	Sha3Uncles       string                `json:"sha3Uncles"`
	Size             string                `json:"size"`
	StateRoot        string                `json:"stateRoot"`
	Timestamp        string                `json:"timestamp"`
	TotalDifficulty  string                `json:"totalDifficulty"`
	Transactions     []EthBlockTransaction `json:"transactions"`
	TransactionsRoot string                `json:"transactionsRoot"`
	Uncles           []interface{}         `json:"uncles"` // Empty in provided example, so using generic interface
	Validator        string                `json:"validator"`
	Validators       interface{}           `json:"validators"` // mainchain and subnet has different format
}

type EthBlockTransaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

func (r *EthBlockResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}
