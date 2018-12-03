package core

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//the least number of zero at begin of valid hash
const targetZeroBit = 16

type ProofOfWork struct {
	block  *Block  //block to be verified
	target big.Int //difficult-first zero number
}

func NewProofOfWork(block *Block) *ProofOfWork {
	// initialize target
	target := big.NewInt(1)
	target = target.Lsh(target, 256-targetZeroBit)

	return &ProofOfWork{block, *target}
}

//prepare []byte data
func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevHash,
			pow.block.Data,
			Int2ByteArray(pow.block.Height),
			Int2ByteArray(pow.block.Timestamp),
			Int2ByteArray(int64(targetZeroBit)),
			Int2ByteArray(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

//verify if the hash is valid
func (pow *ProofOfWork) IsValid() bool {
	//compare the block hash with the target

	var hashInt big.Int
	hashInt.SetBytes(pow.block.Hash)

	if pow.target.Cmp(&hashInt) == 1 {
		return true
	}
	return false
}

//running of POW
func (pow *ProofOfWork) Run() ([]byte, int64) {
	var nonce int64 = 0
	var hashInt *big.Int = big.NewInt(0)
	var hash [32]byte
	for {
		//transfer the block to []byte
		dataBytes := pow.prepareData(nonce)
		//generate hash
		hash = sha256.Sum256(dataBytes)

		//verify hash, return if satisfied
		hashInt.SetBytes(hash[:])
		if pow.target.Cmp(hashInt) == 1 {
			fmt.Printf("\r%x\n", hash)
			break
		}
		nonce = nonce + 1
	}

	return hash[:], nonce
}
