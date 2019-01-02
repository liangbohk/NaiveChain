package core

type UTXO struct {
	//hash
	TXHash []byte
	//index of the output in the tx
	Index int
	//output
	Output *TXOutput
}
