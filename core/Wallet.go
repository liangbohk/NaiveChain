package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Wallet struct {
	//private key
	PrivateKey ecdsa.PrivateKey

	//corresponding public key
	PublicKey []byte
}

//create a wallet
func NewWallet() *Wallet {

	privateKey, publicKey := newKeyPair()
	return &Wallet{privateKey, publicKey}
}

//get a privatekey,publickey pair
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privateKeyPtr, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		log.Panic(err)
	}

	publicKey := append(privateKeyPtr.PublicKey.X.Bytes(), privateKeyPtr.PublicKey.Y.Bytes()...)

	return *privateKeyPtr, publicKey

}

func (w *Wallet) IsValidAddress(address []byte) bool {
	versionAndRipemdHashAndChecksum := Base58Decode(address)
	versionAndRipemdHash := versionAndRipemdHashAndChecksum[:len(versionAndRipemdHashAndChecksum)-addressChecksumLen]
	checkSumBytes := versionAndRipemdHashAndChecksum[len(versionAndRipemdHashAndChecksum)-addressChecksumLen:]
	if bytes.Compare(checkSumBytes, CheckSum(versionAndRipemdHash)) == 0 {
		return true
	}

	return false
}

//get an address
func (w *Wallet) GetAddress() []byte {
	//sha256 and ripemd160
	publicKeyHash := w.Ripemd160Hash(w.PublicKey)
	versionAndRipemdHash := append([]byte{version}, publicKeyHash...)
	checkSumBytes := CheckSum(versionAndRipemdHash)

	tmpBytes := append(versionAndRipemdHash, checkSumBytes...)

	return Base58Encode(tmpBytes)

}

//ripemd160
func (w *Wallet) Ripemd160Hash(publicKey []byte) []byte {
	//256
	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash := hash256.Sum(nil)
	//160
	hash160 := ripemd160.New()
	hash160.Write(hash)
	return hash160.Sum(nil)
}

func CheckSum(bytes []byte) []byte {
	hash1 := sha256.Sum256(bytes)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:addressChecksumLen]
}
