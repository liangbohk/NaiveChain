package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Transaction struct {
	//transaction hash
	TxHash []byte

	//input and output
	TxIns  []*TXInput
	TxOuts []*TXOutput
}

//transaction from coinbase
func NewCoinbaseTransaction(address string) *Transaction {
	txInput := &TXInput{[]byte{}, -1, "Genesis Data"}
	txOutput := &TXOutput{100, address}
	txCoinbase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	txCoinbase.AttachHash()
	return txCoinbase
}

//attach a hash to the transaction
func (tx *Transaction) AttachHash() {
	var res bytes.Buffer
	//initialize an encoder
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(res.Bytes())
	tx.TxHash = hash[:]
}

func NewSimpleTransaction(from string, to string, amount int) *Transaction {

	//build a tx input array
	var txIns []*TXInput
	txInput := &TXInput{[]byte("A Simple Transaction Hash"), 0, from}
	txIns = append(txIns, txInput)

	var txOuts []*TXOutput
	//send
	txOutput := &TXOutput{int64(amount), to}
	txOuts = append(txOuts, txOutput)
	//change
	txOutput = &TXOutput{100 - int64(amount), to}
	txOuts = append(txOuts, txOutput)

	tx := &Transaction{[]byte{}, txIns, txOuts}
	tx.AttachHash()
	return tx
}
