package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

type Transaction struct {
	//transaction hash
	TxHash []byte

	//input and output
	TxIns  []*TXInput
	TxOuts []*TXOutput
}

//chech curTx is coinbase or not
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return len(tx.TxIns[0].TxHash) == 0 && tx.TxIns[0].TxOutIndex == -1
}

//transaction from coinbase
func NewCoinbaseTransaction(address string) *Transaction {
	txInput := &TXInput{[]byte{}, -1, "Genesis Data"}
	txOutput := &TXOutput{10, address}
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

func NewSimpleTransaction(from string, to string, amount int, blc *Blockchain, txs []*Transaction) *Transaction {

	//find usable UTXOs
	restValue, dic := blc.FindSpendableUTXOs(from, amount, txs)
	fmt.Printf("restValue:%d, dic:%x", restValue, dic)

	//build a tx input array
	var txIns []*TXInput
	for hash, indexArray := range dic {
		for _, i := range indexArray {
			txInput := &TXInput{[]byte(hash), i, from}
			txIns = append(txIns, txInput)
		}
	}

	//txInput := &TXInput{[]byte("A Simple Transaction Hash"), 0, from}
	//txIns = append(txIns, txInput)

	var txOuts []*TXOutput
	//send
	txOutput := &TXOutput{int64(amount), to}
	txOuts = append(txOuts, txOutput)
	//change
	txOutput = &TXOutput{int64(restValue) - int64(amount), from}
	txOuts = append(txOuts, txOutput)

	tx := &Transaction{[]byte{}, txIns, txOuts}
	tx.AttachHash()
	return tx
}
