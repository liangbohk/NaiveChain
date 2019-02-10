package core

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Version struct {
	Version     int
	BaseHeight  int
	AddressFrom string
}

//serialize version structure
func (v *Version) Serialize() []byte {
	var res bytes.Buffer
	//initialize an encoder
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(v)
	if err != nil {
		log.Panic(err)
	}

	return res.Bytes()
}

//deserialize version structure
func DeserializeVersion(versionBytes []byte) *Version {
	var version Version

	decoder := gob.NewDecoder(bytes.NewReader(versionBytes))
	err := decoder.Decode(&version)
	if err != nil {
		log.Panic(err)
	}

	return &version
}
