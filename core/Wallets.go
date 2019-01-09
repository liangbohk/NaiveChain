package core

//define wallets struct and the wallets are unordered
type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallets() *Wallets {
	wallets := make(map[string]*Wallet)
	return &Wallets{wallets}
}

//create a new wallet for a wallets obj
func (w *Wallets) CreateNewWallet() {
	wallet := NewWallet()
	w.Wallets[string(wallet.GetAddress())] = wallet
}
