package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
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
	txInput := &TXInput{[]byte{}, -1, nil, []byte{}}
	txOutput := NewTXOutput(10, address)
	//txOutput := &TXOutput{10, address}
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

	wallets, _ := NewWallets()
	wallet := wallets.WalletsMap[from]

	//find usable UTXOs
	restValue, dic := blc.FindSpendableUTXOs(from, amount, txs)
	fmt.Printf("restValue:%d, dic:%x\n", restValue, dic)

	//build a tx input array
	var txIns []*TXInput
	for hash, indexArray := range dic {
		for _, i := range indexArray {
			txInput := &TXInput{[]byte(hash), i, nil, wallet.PublicKey}
			txIns = append(txIns, txInput)
		}
	}

	//txInput := &TXInput{[]byte("A Simple Transaction Hash"), 0, from}
	//txIns = append(txIns, txInput)

	var txOuts []*TXOutput
	//send
	txOutput := NewTXOutput(int64(amount), to)
	//txOutput := &TXOutput{int64(amount), to}
	txOuts = append(txOuts, txOutput)
	//change
	txOutput = NewTXOutput(int64(restValue)-int64(amount), from)
	//txOutput = &TXOutput{int64(restValue) - int64(amount), from}
	txOuts = append(txOuts, txOutput)

	tx := &Transaction{[]byte{}, txIns, txOuts}
	tx.AttachHash()

	//signature
	blc.SignTransaction(tx, wallet.PrivateKey, txs)

	return tx
}

//trimmed copy
func (tx *Transaction) TrimmedCopy() Transaction {
	var txInputs []*TXInput
	var txOutputs []*TXOutput

	for _, txInput := range tx.TxIns {
		txInputs = append(txInputs, &TXInput{txInput.TxHash, txInput.TxOutIndex, nil, nil})
	}
	for _, txOutput := range tx.TxOuts {
		txOutputs = append(txOutputs, &TXOutput{txOutput.Value, txOutput.Sha256Ripemd160HashPubkey})
	}

	txCopy := Transaction{tx.TxHash, txInputs, txOutputs}
	return txCopy
}

//return a hash
func (tx *Transaction) Hash() []byte {
	txCopy := tx
	txCopy.TxHash = []byte{}
	txCopy.AttachHash()
	return txCopy.TxHash[:]
}

//sign a transaction with a privateKey
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbaseTransaction() {
		return
	}
	//fmt.Println(prevTXs)
	for _, txInput := range tx.TxIns {
		if prevTXs[hex.EncodeToString(txInput.TxHash)].TxHash == nil {
			fmt.Println(txInput)
			fmt.Println(prevTXs[hex.EncodeToString(txInput.TxHash)].TxHash)
			log.Panic("previous transaction fault!")
		}
	}

	txCopy := tx.TrimmedCopy()

	for txInputIndex, txInput := range txCopy.TxIns {
		prevTx := prevTXs[hex.EncodeToString(txInput.TxHash)]
		txCopy.TxIns[txInputIndex].Signature = nil
		txCopy.TxIns[txInputIndex].Pubkey = prevTx.TxOuts[txInput.TxOutIndex].Sha256Ripemd160HashPubkey
		txCopy.TxHash = txCopy.Hash()
		txCopy.TxIns[txInputIndex].Pubkey = nil

		//signature
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.TxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.TxIns[txInputIndex].Signature = signature
	}

}
