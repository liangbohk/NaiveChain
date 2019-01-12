package core

import "bytes"

type TXInput struct {
	//transaction id
	TxHash []byte

	//txoutput index
	TxOutIndex int

	//signature and public key
	Signature []byte
	Pubkey    []byte
}

//check if the sig equals address
func (txInput *TXInput) UnLockRipemd160Hash(sha256Ripemd160HashPubkey []byte) bool {
	publicKey := Sha256Ripemd160Hash(txInput.Pubkey)
	return bytes.Compare(publicKey, sha256Ripemd160HashPubkey) == 0
}
