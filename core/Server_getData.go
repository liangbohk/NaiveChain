package core

import (
	"bytes"
	"encoding/gob"
	"log"
)

type GetData struct {
	AddrFrom string
	Type     string
	Hash     []byte
}

//serialize getdata structure
func (gd *GetData) Serialize() []byte {
	var res bytes.Buffer
	//initialize an encoder
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(gd)
	if err != nil {
		log.Panic(err)
	}

	return res.Bytes()
}

//deserialize getdata structure
func DeserializeGetData(getDataBytes []byte) *GetData {
	var getData GetData

	decoder := gob.NewDecoder(bytes.NewReader(getDataBytes))
	err := decoder.Decode(&getData)
	if err != nil {
		log.Panic(err)
	}

	return &getData
}
