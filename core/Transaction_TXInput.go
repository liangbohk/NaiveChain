package core

type TXInput struct {
	//transaction id
	TxHash []byte

	//txoutput index
	TxOutIndex int

	//public key
	ScriptSig string
}

//check if the sig equals address
func (txInput *TXInput) UnLockScriptSigWithAddress(address string) bool {
	return txInput.ScriptSig == address
}
