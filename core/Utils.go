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

//convert json to array
func Json2Array(jsonStr string) []string {
	var s []string
	if err := json.Unmarshal([]byte(jsonStr), &s); err != nil {
		log.Panic(err)
	}
	return s
}

//reverse the byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
