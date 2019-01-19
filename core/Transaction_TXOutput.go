package core

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TXOutput struct {
	//value
	Value int64

	//public key
	Sha256Ripemd160HashPubkey []byte
}

//serialize the txoutput
func (txOutput *TXOutput) Serialize() []byte {
	var res bytes.Buffer
	//initialize an encoder
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(txOutput)
	if err != nil {
		log.Panic(err)
	}

	return res.Bytes()
}

//deserialize the txoutput bytes
func DeserializeTXOutput(txOutputBytes []byte) *TXOutput {
	var txOutput TXOutput

	decoder := gob.NewDecoder(bytes.NewReader(txOutputBytes))
	err := decoder.Decode(&txOutput)
	if err != nil {
		log.Panic(err)
	}

	return &txOutput
}

//check if the sig equals address
func (txOutput *TXOutput) UnLockScriptPubkeyWithAddress(address string) bool {
	publicKeyHash := Base58Decode([]byte(address))
	return bytes.Compare(txOutput.Sha256Ripemd160HashPubkey, publicKeyHash[1:len(publicKeyHash)-addressChecksumLen]) == 0
}

//create a new txoutput
func NewTXOutput(value int64, address string) *TXOutput {
	txOutput := &TXOutput{value, nil}
	//setup public key (sha256 ripemd160)
	txOutput.Lock(address)

	return txOutput
}

//lock a txoutput, in fact convert address to pubkey
func (txOutput *TXOutput) Lock(address string) {
	publicKeyHash := Base58Decode([]byte(address))
	//the first "1" is version id and last addressChecksumLen is checksum
	txOutput.Sha256Ripemd160HashPubkey = publicKeyHash[1 : len(publicKeyHash)-addressChecksumLen]
}
