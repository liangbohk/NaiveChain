package core

import (
	"github.com/boltdb/bolt"
	"log"
	"math/big"
)

//the iterator structure of the blockchain
type BlockchainIterator struct {
	CurHash []byte
	DB      *bolt.DB
}

func (blcIter *BlockchainIterator) Next() *Block {
	//check if current block is genesis block
	var curHashInt big.Int
	curHashInt.SetBytes(blcIter.CurHash)
	if (big.NewInt(0).Cmp(&curHashInt)) == 0 {
		return nil
	}

	var curBlock *Block
	err := blcIter.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tableName))
		if bucket != nil {
			curBlock = Deserialize(bucket.Get(blcIter.CurHash))
			blcIter.CurHash = curBlock.PrevHash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return curBlock
}
