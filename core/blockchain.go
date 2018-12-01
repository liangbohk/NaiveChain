package core

//blockchain structure
type Blockchain struct {
	Blocks []*Block //store ordered blocks
}

//generate a blockchain with a genesis block
func CreateBlockchainWithAGenesisBlock() *Blockchain {
	//generate a genesis block
	genesisBlock := CreateGenesisBlock("Genesis block")

	return &Blockchain{[]*Block{genesisBlock}}
}

//add a block to the blockchain
func (blc *Blockchain) AddBlockToBlockchain(data string, height int64, predHash []byte) {
	block := NewBlock(data, height, predHash)
	blc.Blocks = append(blc.Blocks, block)
}
