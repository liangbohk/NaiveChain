package core

type TXInput struct {
	//transaction id
	TxHash []byte

	//txoutput index
	TxOutIndex int

	//public key
	ScriptSig string
}
