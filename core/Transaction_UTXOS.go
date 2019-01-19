package core

import (
	"bytes"
	"encoding/gob"
	"log"
)

type UTXOS struct {
	UTXOs []*UTXO
}

//serialize the txoutputs
func (txOutputs *UTXOS) Serialize() []byte {
	var res bytes.Buffer
	//initialize an encoder
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(txOutputs)
	if err != nil {
		log.Panic(err)
	}

	return res.Bytes()
}

//deserialize the txoutputs bytes
func DeserializeUTXOS(txOutputsBytes []byte) *UTXOS {
	var txOutputs UTXOS

	decoder := gob.NewDecoder(bytes.NewReader(txOutputsBytes))
	err := decoder.Decode(&txOutputs)
	if err != nil {
		log.Panic(err)
	}

	return &txOutputs
}
