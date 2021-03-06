package core

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

//database name
const dbName = "naivechain_%s.db"
const tableName = "blocks"

//blockchain structure
type Blockchain struct {
	//Blocks []*Block //store ordered blocks, no need if there is a db

	Tail []byte //newest block hash
	DB   *bolt.DB
}

//check if the db exist
func DBExist(dbName string) bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

//generate a blockchain with a genesis block
func CreateBlockchainWithAGenesisBlock(address string, nodeID string) *Blockchain {
	//specify db file name
	dbName := fmt.Sprintf(dbName, nodeID)

	//check if the db exist
	if DBExist(dbName) {
		fmt.Println("genesis block already exists!")
		os.Exit(1)
	}

	fmt.Println("creating the genesis block")
	//initialize the database
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var genesisHash []byte
	err = db.Update(func(tx *bolt.Tx) error {

		bucket, err := tx.CreateBucket([]byte(tableName))
		if err != nil {
			log.Panic(err)
		}

		if bucket != nil {
			//create transaction
			txCoinbase := NewCoinbaseTransaction(0, address)

			//generate a genesis block
			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
			//save genesisBlock to the table
			err = bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			//save the newest block hash
			err = bucket.Put([]byte("tail"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			genesisHash = genesisBlock.Hash
		}
		return nil
	})

	return &Blockchain{genesisHash, db}
}

//add a block to the blockchain
func (blc *Blockchain) AddBlockToBlockchain(txs []*Transaction) {
	//save the block to the database
	err := blc.DB.Update(func(tx *bolt.Tx) error {

		//get the table from the database
		table := tx.Bucket([]byte(tableName))
		//get the newest block
		lastBlock := DeserializeBlock(table.Get(blc.Tail))

		//create a new block
		block := NewBlock(txs, lastBlock.Height+1, lastBlock.Hash)
		//save the block to the database
		err := table.Put(block.Hash, block.Serialize())
		if err != nil {
			log.Panic(err)
		}
		//update the tail in the database
		err = table.Put([]byte("tail"), block.Hash)
		if err != nil {
			log.Panic(err)
		}
		//update the tail in the blockchain
		blc.Tail = block.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

}

func (blc *Blockchain) VerifyTransactions(tx *Transaction, txs []*Transaction) bool {
	prevTXs := make(map[string]Transaction)
	for _, txInput := range tx.TxIns {
		prevTX, err := blc.FindTransaction(txInput.TxHash, txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}

	return tx.Verify(prevTXs)
}

//return a point to the iterator
func (blc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blc.Tail, blc.DB}
}

//print the info of a blockchain by iterating visit
func (blc *Blockchain) PrintChain() {
	//initialize a blockchain iterator
	blcIter := blc.Iterator()
	for {
		block := blcIter.Next()
		if block == nil {
			return
		}
		fmt.Println("-----------------------------------------------------------")
		fmt.Printf("Height: %d\n PredHash:%x\n TimeStamp:%s\n Hash:%x\n Nonce:%d\n",
			block.Height, block.PrevHash, time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"), block.Hash, block.Nonce)
		fmt.Println("Transactions:")
		for _, tx := range block.Txs {
			fmt.Printf("%x\n", tx.TxHash)
			for _, input := range tx.TxIns {
				fmt.Println("TXInput:")
				fmt.Printf("\t%s\n", hex.EncodeToString(input.TxHash))
				fmt.Printf("\t%d\n", input.TxOutIndex)
				fmt.Printf("\t%x\n", input.Pubkey)
			}

			for _, output := range tx.TxOuts {
				fmt.Println("TXOutput:")
				fmt.Printf("\tvalue:%d\n", output.Value)
				fmt.Printf("\tsha256Ripemd160HashPubkey:%x\n", output.Sha256Ripemd160HashPubkey)
			}
		}

	}

}

//get a blockchain object
func BlockchainObject(nodeID string) *Blockchain {
	//specify the db file name
	dbName := fmt.Sprintf(dbName, nodeID)

	//check if the database file exist

	if !DBExist(dbName) {
		log.Fatal("no blockchain")
	}

	//open the database
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var tailHash []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))

		if b != nil {
			tailHash = b.Get([]byte("tail"))

		} else {
			err := os.Remove(dbName)
			if err != nil {
				log.Panic(err)
			}
			//log.Fatal("no blockchain")
		}
		return nil

	})

	return &Blockchain{tailHash, db}
}

func (blc *Blockchain) GetBlockchainHeight() (int64, error) {
	var height int64 = 0
	err := blc.DB.View(func(tx *bolt.Tx) error {
		//get the table from the database
		table := tx.Bucket([]byte(tableName))
		if table != nil {
			//get the newest block
			lastBlock := DeserializeBlock(table.Get(blc.Tail))
			height = lastBlock.Height
		} else {
			log.Fatal("no blockchain")
		}
		return nil

	})
	return height, err
}

//return transactions with unspent TXOutput
func (blc *Blockchain) UnspentTxOuts(address string, txs []*Transaction) []*UTXO {

	//store unspent transactions
	var UTXOs []*UTXO
	//store spent txoutput
	spentTXOutputs := make(map[string][]int)

	//traverse txs
	for _, tx := range txs {
		//txinput
		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.TxIns {
				versonSha256Ripemd160HashChecksum := Base58Decode([]byte(address))
				if in.UnLockRipemd160Hash(versonSha256Ripemd160HashChecksum[1 : len(versonSha256Ripemd160HashChecksum)-addressChecksumLen]) {
					key := hex.EncodeToString(in.TxHash)
					spentTXOutputs[key] = append(spentTXOutputs[key], in.TxOutIndex)
				}
			}
		}

	}
	for _, tx := range txs {
		//txoutput
	loop0:
		for index, out := range tx.TxOuts {
			if out.UnLockScriptPubkeyWithAddress(address) {
				if spentTXOutputs != nil {
					if len(spentTXOutputs) != 0 {
						for txHash, indexArr := range spentTXOutputs {
							for _, i := range indexArr {
								if txHash == hex.EncodeToString(tx.TxHash) && index == i {
									continue loop0
								}
							}

						}
						//if not be spend, add it as a utxo
						utxo := &UTXO{tx.TxHash, index, out}
						UTXOs = append(UTXOs, utxo)
					} else {
						utxo := &UTXO{tx.TxHash, index, out}
						UTXOs = append(UTXOs, utxo)
					}

				}

			}
		}
	}

	//get iterator
	iter := blc.Iterator()

	for {
		block := iter.Next()

		//deal with transaction
		for _, tx := range block.Txs {
			//txinput
			if tx.IsCoinbaseTransaction() == false {
				for _, in := range tx.TxIns {
					versonSha256Ripemd160HashChecksum := Base58Decode([]byte(address))
					if in.UnLockRipemd160Hash(versonSha256Ripemd160HashChecksum[1 : len(versonSha256Ripemd160HashChecksum)-addressChecksumLen]) {
						//if in.UnLockScriptSigWithAddress(address) {
						key := hex.EncodeToString(in.TxHash)
						spentTXOutputs[key] = append(spentTXOutputs[key], in.TxOutIndex)
					}
				}
			}
		}

		for _, tx := range block.Txs {
			//txoutput
		loop1:
			for index, out := range tx.TxOuts {
				if out.UnLockScriptPubkeyWithAddress(address) {
					if spentTXOutputs != nil {
						if len(spentTXOutputs) != 0 {
							for txHash, indexArr := range spentTXOutputs {
								for _, i := range indexArr {
									if txHash == hex.EncodeToString(tx.TxHash) && index == i {
										continue loop1

									}
								}

							}
							//if not be spend, add it as a utxo
							utxo := &UTXO{tx.TxHash, index, out}
							UTXOs = append(UTXOs, utxo)
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							UTXOs = append(UTXOs, utxo)
						}

					}
				}
			}

		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevHash)

		//stop condition
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return UTXOs
}

//get balance
func (blc *Blockchain) GetBanlance(address string) int64 {
	utxos := blc.UnspentTxOuts(address, []*Transaction{})
	var amount int64 = 0
	for _, utxo := range utxos {
		amount += utxo.Output.Value
	}
	return amount
}

//find a transaction by a id(txhash)
func (blc *Blockchain) FindTransaction(txHash []byte, txs_packaged []*Transaction) (Transaction, error) {

	//before discussing tx in previous blocks, also you need to discuss txs already packaged in this block
	for _, tx := range txs_packaged {
		if bytes.Compare(txHash, tx.TxHash) == 0 {
			return *tx, nil
		}
	}

	//traverse the blockchain
	iter := blc.Iterator()

	for {
		block := iter.Next()
		for _, tx := range block.Txs {
			//fmt.Printf("%s<=>%s\n",hex.EncodeToString(txHash),hex.EncodeToString(tx.TxHash))
			if bytes.Compare(txHash, tx.TxHash) == 0 {
				return *tx, nil
			}
		}

		//stop condition
		var hashInt big.Int
		hashInt.SetBytes(block.PrevHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return Transaction{}, nil
}

//sign transaction
func (blc *Blockchain) SignTransaction(tx *Transaction, privateKey ecdsa.PrivateKey, txs_packaged []*Transaction) {
	if tx.IsCoinbaseTransaction() {
		return
	}

	prevTXs := make(map[string]Transaction)
	for _, txInput := range tx.TxIns {
		prevTX, err := blc.FindTransaction(txInput.TxHash, txs_packaged)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}

	tx.Sign(privateKey, prevTXs)

}

//sign transaction with private key

//find utxo that can be used
func (blc *Blockchain) FindSpendableUTXOs(from string, amount int64, txs []*Transaction) (int64, map[string][]int) {
	//get utxo
	utxos := blc.UnspentTxOuts(from, txs)

	//traverse utxos
	var value int64 = 0
	selectedUTXOs := make(map[string][]int)
	for _, utxo := range utxos {
		value = value + utxo.Output.Value
		txHash := hex.EncodeToString(utxo.TXHash)
		selectedUTXOs[txHash] = append(selectedUTXOs[txHash], utxo.Index)
		//selectedUTXOs[string(utxo.TXHash)] = append(selectedUTXOs[string(utxo.TXHash)], utxo.Index)
		if value >= int64(amount) {
			break
		}
	}

	if value < int64(amount) {

		fmt.Printf("no enough balance %d, need %d ", value, amount)
		os.Exit(1)
	}

	return value, selectedUTXOs
}

func (blc *Blockchain) MineNewBlock(from []string, to []string, amount []string, nodeID string) *Blockchain {

	//build new transactions
	//set up the transactions.Note: the tx order in txs is specially setup
	blockHeight, err := blc.GetBlockchainHeight()
	if err != nil {
		log.Panic(err)
	}

	//utxoset
	utxoSet := &UTXOSet{blc}

	//build transactions
	var txs []*Transaction
	for index, _ := range from {
		value, err := strconv.Atoi(amount[index])
		if err != nil {
			log.Panic(err)
		}
		tx := NewSimpleTransaction(blockHeight+1, from[index], to[index], int64(value), utxoSet, txs, nodeID)
		txs = append(txs, tx)
	}

	//coinbase transaction,let from[0] be the miner node
	tx := NewCoinbaseTransaction(blockHeight+1, from[0])
	txs = append(txs, tx)

	//verify the signature of txs
	tmpTxs := []*Transaction{}
	for _, tx := range txs {
		if !blc.VerifyTransactions(tx, tmpTxs) {
			log.Panic("verifying signature failed")
		}
		tmpTxs = append(tmpTxs, tx)
	}

	blc.AddBlockToBlockchain(tmpTxs)

	return blc

}

//return UTXOMap
func (blc *Blockchain) FindUTXOMap() map[string]UTXOS {
	//traverse the blockchain
	iter := blc.Iterator()
	//spent
	spentUTXOMap := make(map[string][]*TXInput)
	utxoMap := make(map[string]UTXOS)
	for {
		block := iter.Next()
		for i := len(block.Txs) - 1; i >= 0; i-- {
			utxos := UTXOS{[]*UTXO{}}

			tx := block.Txs[i]
			if !tx.IsCoinbaseTransaction() {
				for _, txInput := range tx.TxIns {
					txHash := hex.EncodeToString(txInput.TxHash)
					spentUTXOMap[txHash] = append(spentUTXOMap[txHash], txInput)
				}
			}

			txHash := hex.EncodeToString(tx.TxHash)
		work:
			for index, txOutput := range tx.TxOuts {
				txInputs := spentUTXOMap[txHash]
				if len(txInputs) > 0 {
					for _, txInput := range txInputs {
						outPubkey := txOutput.Sha256Ripemd160HashPubkey
						inPubkey := txInput.Pubkey
						if bytes.Compare(outPubkey, Sha256Ripemd160Hash(inPubkey)) == 0 {
							if index == txInput.TxOutIndex {
								continue work
							}
						}

					}
					utxo := &UTXO{tx.TxHash, index, txOutput}
					utxos.UTXOs = append(utxos.UTXOs, utxo)
				} else {
					utxo := &UTXO{tx.TxHash, index, txOutput}
					utxos.UTXOs = append(utxos.UTXOs, utxo)
				}
			}

			utxoMap[txHash] = utxos

		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return utxoMap
}

func (blc *Blockchain) GetBlockHashes() [][]byte {

	iter := blc.Iterator()
	var hashes [][]byte

	for {
		block := iter.Next()
		hashes = append(hashes, block.Hash)

		var hashInt big.Int
		hashInt.SetBytes(block.PrevHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return hashes
}
