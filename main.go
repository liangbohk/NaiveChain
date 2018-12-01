package main

import (
	"NaiveChain/core"
	"fmt"
)

func main() {
	block := core.NewBlock("Genesis Block", 0, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	fmt.Println(block)
	fmt.Println(string(block.Hash))

	//fmt.Println(strconv.FormatInt(54,2))
	//fmt.Println([]byte(strconv.FormatInt(54,2)))
	//fmt.Println(core.Int2ByteArray(54))
	//fmt.Println(core.Int2ByteArray(1))
	//fmt.Println('1')
}
