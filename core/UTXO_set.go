package core

import (
	"encoding/hex"
	"github.com/boltdb/bolt"
	"log"
)

//save all UTXO to database

const utxoTableName = "utxoTable"

type UTXOSet struct {
	Blc *Blockchain
}

//reset the UTXOSet database
func (utxoSet *UTXOSet) Reset() {
	err := utxoSet.Blc.DB.Update(func(tx *bolt.Tx) error {
		table := tx.Bucket([]byte(utxoTableName))
		if table != nil {
			err := tx.DeleteBucket([]byte(utxoTableName))
			if err != nil {
				log.Panic(err)
			}
		}
		table, _ = tx.CreateBucket([]byte(utxoTableName))
		if table != nil {
			txOutputsMap := utxoSet.Blc.FindUTXOMap()

			for hash, txOutputs := range txOutputsMap {
				txHash, _ := hex.DecodeString(hash)
				table.Put(txHash, txOutputs.Serialize())
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func (utxoSet *UTXOSet) findUTXOWithAddress(address string) []*UTXO {
	var utxos []*UTXO

	utxoSet.Blc.DB.View(func(tx *bolt.Tx) error {
		table := tx.Bucket([]byte(utxoTableName))
		if table == nil {
			log.Panic("no utxo table")
		}
		c := table.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {

			utxoss := DeserializeUTXOS(v)
			for _, utxo := range utxoss.UTXOs {
				if utxo.Output.UnLockScriptPubkeyWithAddress(address) {
					utxos = append(utxos, utxo)
				}
			}
		}
		return nil
	})

	return utxos
}

//get balance of an address
func (utxoSet *UTXOSet) GetBalance(address string) int64 {
	utxos := utxoSet.findUTXOWithAddress(address)

	var amount int64
	for _, utxo := range utxos {
		amount += utxo.Output.Value
	}

	return amount
}
