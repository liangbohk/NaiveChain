package core

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
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
				fmt.Printf("%s\n", input.TxHash)
				fmt.Printf("%d\n", input.TxOutIndex)
				fmt.Printf("%s\n", input.ScriptSig)
			}

			for _, output := range tx.TxOuts {
				fmt.Println("TXOutput:")
				fmt.Println(output.Value)
				fmt.Println(output.ScriptPubkey)
				fmt.Println("\n")

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

func (blc *Blockchain) MineNewBlock(from []string, to []string, amount []string) *Blockchain {
	fmt.Println(from)
	fmt.Println(to)
	fmt.Println(amount)

	//set up the transactions
	var txs []*Transaction

	blc.AddBlockToBlockchain(txs)

	////get newest block info
	//var block *Block
	//blc.DB.View(func(tx *bolt.Tx) error {
	//	b:=tx.Bucket([]byte(tableName))
	//	if b!=nil{
	//		hash:=b.Get([]byte("tail"))
	//		block=Deserialize(b.Get(hash))
	//	}
	//	return nil
	//})
	//
	////build a new block
	//newBlock:=NewBlock(txs,block.Height+1,block.Hash)
	//
	////save new block to db
	//blc.DB.Update(func(tx *bolt.Tx) error {
	//	b:=tx.Bucket([]byte(tableName))
	//	if b!=nil{
	//		b.Put(newBlock.Hash,newBlock.Serialize())
	//		b.Put([]byte("tail"),newBlock.Hash)
	//		blc.Tail=newBlock.Hash
	//	}
	//	return nil
	//})
	return blc

}
