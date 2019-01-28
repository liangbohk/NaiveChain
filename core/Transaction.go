package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

type Transaction struct {
	//the height of the block where the transaction packaged in
	BlockHeight int64

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
func NewCoinbaseTransaction(blockHeight int64, address string) *Transaction {
	txInput := &TXInput{[]byte{}, -1, nil, []byte{}}
	txOutput := NewTXOutput(10, address)
	//txOutput := &TXOutput{10, address}
	txCoinbase := &Transaction{blockHeight, []byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	txCoinbase.AttachHash()
	fmt.Printf("coinbase %s\n", hex.EncodeToString(txCoinbase.TxHash))
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

func NewSimpleTransaction(blockHeight int64, from string, to string, amount int64, utxoSet *UTXOSet, txs []*Transaction, nodeID string) *Transaction {

	wallets, _ := NewWallets(nodeID)
	wallet := wallets.WalletsMap[from]

	//find usable UTXOs
	restValue, dic := utxoSet.Blc.FindSpendableUTXOs(from, amount, txs)
	//restValue, dic := utxoSet.FindSpendableUTXOS(from, amount, txs)
	fmt.Println("-------------------------------")
	fmt.Println(restValue)
	fmt.Println(dic)

	//build a tx input array
	var txIns []*TXInput
	for hash, indexArray := range dic {
		for _, i := range indexArray {
			hashBytes, err := hex.DecodeString(hash)
			if err != nil {
				log.Println(err)
			}
			txInput := &TXInput{hashBytes, i, nil, wallet.PublicKey}
			txIns = append(txIns, txInput)
		}
	}

	var txOuts []*TXOutput
	//send
	txOutput := NewTXOutput(int64(amount), to)
	//txOutput := &TXOutput{int64(amount), to}
	txOuts = append(txOuts, txOutput)
	//change
	txOutput = NewTXOutput(int64(restValue)-int64(amount), from)
	//txOutput = &TXOutput{int64(restValue) - int64(amount), from}
	txOuts = append(txOuts, txOutput)

	tx := &Transaction{blockHeight, []byte{}, txIns, txOuts}
	tx.AttachHash()
	fmt.Println(amount)
	fmt.Println(hex.EncodeToString(tx.TxHash))
	fmt.Println("-------------------------------")

	//signature
	utxoSet.Blc.SignTransaction(tx, wallet.PrivateKey, txs)

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

	txCopy := Transaction{tx.BlockHeight, tx.TxHash, txInputs, txOutputs}
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
	for _, txInput := range tx.TxIns {
		if prevTXs[hex.EncodeToString(txInput.TxHash)].TxHash == nil {

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

//verify transaction signature
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbaseTransaction() {
		return true
	}

	for _, txInput := range tx.TxIns {
		if prevTXs[hex.EncodeToString(txInput.TxHash)].TxHash == nil {
			log.Panic("not valid previous transaction")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for txInputIndex, txInput := range tx.TxIns {
		prevTx := prevTXs[hex.EncodeToString(txInput.TxHash)]
		txCopy.TxIns[txInputIndex].Signature = nil
		txCopy.TxIns[txInputIndex].Pubkey = prevTx.TxOuts[txInput.TxOutIndex].Sha256Ripemd160HashPubkey
		txCopy.TxHash = txCopy.Hash()
		txCopy.TxIns[txInputIndex].Pubkey = nil

		//private key
		r, s := big.Int{}, big.Int{}
		signatureLength := len(txInput.Signature)
		r.SetBytes(txInput.Signature[:(signatureLength / 2)])
		s.SetBytes(txInput.Signature[(signatureLength / 2):])

		x, y := big.Int{}, big.Int{}
		keyLength := len(txInput.Pubkey)
		x.SetBytes(txInput.Pubkey[:(keyLength / 2)])
		y.SetBytes(txInput.Pubkey[(keyLength / 2):])

		rawPubkey := ecdsa.PublicKey{curve, &x, &y}
		if !ecdsa.Verify(&rawPubkey, txCopy.TxHash, &r, &s) {
			return false
		}

	}

	return true

}
