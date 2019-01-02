package core

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
)

//int64 to byte array
func Int2ByteArray(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func Json2Array(jsonStr string) []string {
	var s []string
	if err := json.Unmarshal([]byte(jsonStr), &s); err != nil {
		log.Panic(err)
	}
	return s
}
