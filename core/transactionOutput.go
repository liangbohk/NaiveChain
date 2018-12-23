package core

type TXOutput struct {
	//value
	Value int64

	//public key
	ScriptPubkey string
}

//check if the sig equals address
func (txOutput *TXOutput) UnLockScriptPubkeyWithAddress(address string) bool {
	return txOutput.ScriptPubkey == address
}
