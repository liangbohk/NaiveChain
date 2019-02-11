package core

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

//serialize inv structure
func (inv *Inv) Serialize() []byte {
	var res bytes.Buffer
	//initialize an encoder
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(inv)
	if err != nil {
		log.Panic(err)
	}

	return res.Bytes()
}

//deserialize inv structure
func DeserializeInv(invBytes []byte) *Inv {
	var inv Inv

	decoder := gob.NewDecoder(bytes.NewReader(invBytes))
	err := decoder.Decode(&inv)
	if err != nil {
		log.Panic(err)
	}

	return &inv
}
