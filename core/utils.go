package core

import (
	"bytes"
	"encoding/binary"
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
