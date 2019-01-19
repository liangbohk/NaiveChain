package core

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TXInput struct {
	//transaction id
	TxHash []byte

	//txoutput index
	TxOutIndex int

	//signature and public key
	Signature []byte
	Pubkey    []byte
}

//serialize the txinput
func (txInput *TXInput) Serialize() []byte {
	var res bytes.Buffer
	//initialize an encoder
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(txInput)
	if err != nil {
		log.Panic(err)
	}

	return res.Bytes()
}

//deserialize the txinput bytes
func DeserializeTXInput(txInputBytes []byte) *TXInput {
	var txInput TXInput

	decoder := gob.NewDecoder(bytes.NewReader(txInputBytes))
	err := decoder.Decode(&txInput)
	if err != nil {
		log.Panic(err)
	}

	return &txInput
}

//check if the sig equals address
func (txInput *TXInput) UnLockRipemd160Hash(sha256Ripemd160HashPubkey []byte) bool {
	publicKey := Sha256Ripemd160Hash(txInput.Pubkey)
	return bytes.Compare(publicKey, sha256Ripemd160HashPubkey) == 0
}
