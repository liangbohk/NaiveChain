package core

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "wallet_%s.dat"

//define wallets struct and the wallets are unordered
type Wallets struct {
	WalletsMap map[string]*Wallet
}

func NewWallets(nodeID string) (*Wallets, error) {

	walletsMap := make(map[string]*Wallet)
	//load wallet data from local file
	wallets := &Wallets{walletsMap}
	err := wallets.LoadFromFile(nodeID)

	return wallets, err
}

//create a new wallet for a wallets obj
func (w *Wallets) CreateNewWallet(nodeID string) {

	wallet := NewWallet()
	fmt.Println(string(wallet.GetAddress()))
	w.WalletsMap[string(wallet.GetAddress())] = wallet
	w.SaveWallets(nodeID)
}

//save wallet
func (w *Wallets) SaveWallets(nodeID string) {
	walletFile := fmt.Sprintf(walletFile, nodeID)

	var content bytes.Buffer
	//to serialize some types
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

//load wallets data from local file
func (w *Wallets) LoadFromFile(nodeID string) error {
	walletFile := fmt.Sprintf(walletFile, nodeID)
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	content, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	w.WalletsMap = wallets.WalletsMap
	return nil
}
