package core

import (
	"bytes"
	"encoding/gob"
	"log"
)

type GetBlocks struct {
	AddrFrom string
}

//serialize getblocks structure
func (gd *GetBlocks) Serialize() []byte {
	var res bytes.Buffer
	//initialize an encoder
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(gd)
	if err != nil {
		log.Panic(err)
	}

	return res.Bytes()
}

//deserialize getblocks structure
func DeserializeGetBlocks(getBlocksBytes []byte) *GetBlocks {
	var getBlocks GetBlocks

	decoder := gob.NewDecoder(bytes.NewReader(getBlocksBytes))
	err := decoder.Decode(&getBlocks)
	if err != nil {
		log.Panic(err)
	}

	return &getBlocks
}
