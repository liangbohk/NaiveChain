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
const dbName = "naivechain.db"
const tableName = "blocks"

//blockchain structure
type Blockchain struct {
	//Blocks []*Block //store ordered blocks, no need if there is a db

	Tail []byte //newest block hash
	DB   *bolt.DB
}

//check if the db exist
func DBExist() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

//generate a blockchain with a genesis block
func CreateBlockchainWithAGenesisBlock(address string) *Blockchain {
	//check if the db exist
	if DBExist() {
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
			txCoinbase := NewCoinbaseTransaction(address)

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
		lastBlock := Deserialize(table.Get(blc.Tail))
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
				fmt.Printf("%s\n", hex.EncodeToString(input.TxHash))
				fmt.Printf("%d\n", input.TxOutIndex)
				fmt.Printf("%s\n", input.Pubkey)
			}

			for _, output := range tx.TxOuts {
				fmt.Println("TXOutput:")
				fmt.Println(output.Value)
				fmt.Println(output.Sha256Ripemd160HashPubkey)
			}
		}

	}

}

//get a blockchain object
func BlockchainObject() *Blockchain {
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
			log.Fatal("no blockchain")
		}
		return nil

	})

	return &Blockchain{tailHash, db}
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
func (blc *Blockchain) FindSpendableUTXOs(from string, amount int, txs []*Transaction) (int64, map[string][]int) {
	//get utxo
	utxos := blc.UnspentTxOuts(from, txs)

	//traverse utxos
	var value int64 = 0
	selectedUTXOs := make(map[string][]int)
	for _, utxo := range utxos {
		value = value + utxo.Output.Value
		selectedUTXOs[string(utxo.TXHash)] = append(selectedUTXOs[string(utxo.TXHash)], utxo.Index)
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

func (blc *Blockchain) MineNewBlock(from []string, to []string, amount []string) *Blockchain {

	//build new transactions
	//set up the transactions.Note: the tx order in txs is specially setup
	var txs []*Transaction
	for index, _ := range from {
		value, err := strconv.Atoi(amount[index])
		if err != nil {
			log.Panic(err)
		}
		tx := NewSimpleTransaction(from[index], to[index], value, blc, txs)
		txs = append(txs, tx)
	}

	blc.AddBlockToBlockchain(txs)

	return blc

}
