package core

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
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
			utxosMap := utxoSet.Blc.FindUTXOMap()

			for hash, utxos := range utxosMap {
				txHash, _ := hex.DecodeString(hash)
				table.Put(txHash, utxos.Serialize())
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

//get returned values and tx hash and index
func (utxoSet *UTXOSet) FindPackageSpendableUTXOS(from string, txs []*Transaction) []*UTXO {
	//store unspent transactions
	var UTXOs []*UTXO

	//store spent txoutput
	spentTXOutputs := make(map[string][]int)

	//traverse txs
	for _, tx := range txs {
		//txinput
		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.TxIns {
				versonSha256Ripemd160HashChecksum := Base58Decode([]byte(from))
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
			if out.UnLockScriptPubkeyWithAddress(from) {
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

	return UTXOs
}

func (utxoSet *UTXOSet) FindSpendableUTXOS(from string, amount int64, txs []*Transaction) (int64, map[string][]int) {
	unPackagedUTXOS := utxoSet.FindPackageSpendableUTXOS(from, txs)

	spendableUTXO := make(map[string][]int)

	money := int64(0)
	for _, utxo := range unPackagedUTXOS {
		money += utxo.Output.Value
		txHash := hex.EncodeToString(utxo.TXHash)
		spendableUTXO[txHash] = append(spendableUTXO[txHash], utxo.Index)

		if money >= amount {
			return money, spendableUTXO
		}
	}

	//if unpackaged txs has no enough money

	utxoSet.Blc.DB.View(func(tx *bolt.Tx) error {
		table := tx.Bucket([]byte(utxoTableName))
		if table != nil {
			c := table.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {

				utxos := DeserializeUTXOS(v)
				for _, utxo := range utxos.UTXOs {
					if utxo.Output.UnLockScriptPubkeyWithAddress(from) {
						money += utxo.Output.Value
						txHash := hex.EncodeToString(utxo.TXHash)
						spendableUTXO[txHash] = append(spendableUTXO[txHash], utxo.Index)
						if money >= amount {
							return nil
						}
					}

				}
			}
		}
		return nil
	})

	if money < amount {
		fmt.Printf("no enough balance %d, need %d ", money, amount)
		os.Exit(1)
	}

	return money, spendableUTXO
}

//update utxoset
func (utxoSet *UTXOSet) Update() {
	//newest block
	tailBlock := utxoSet.Blc.Iterator().Next()

	//traverse the transaction
	txInputs := []*TXInput{}
	utxoMap := make(map[string]*UTXOS)
	for _, tx := range tailBlock.Txs {
		for _, txInput := range tx.TxIns {
			if len(txInput.TxHash) > 0 {
				txInputs = append(txInputs, txInput)
			}

		}
	}
	// check if txOutput in current tailBlock is spend(that is, its txInput also in this tx)
	for _, tx := range tailBlock.Txs {
		tmpUtxos := &UTXOS{}
	work:
		for index, txOutput := range tx.TxOuts {
			for _, txInput := range txInputs {
				//if the txOutput has been spent
				if txInput.TxOutIndex == index && bytes.Compare(tx.TxHash, txInput.TxHash) == 0 {
					continue work
				}
			}
			tmpUtxos.UTXOs = append(tmpUtxos.UTXOs, &UTXO{tx.TxHash, index, txOutput})

		}
		txHash := hex.EncodeToString(tx.TxHash)
		utxoMap[txHash] = tmpUtxos
	}

	err := utxoSet.Blc.DB.Update(func(tx *bolt.Tx) error {
		table := tx.Bucket([]byte(utxoTableName))
		if table != nil {
			//add utxos
			for txHash, utxos := range utxoMap {
				hash, _ := hex.DecodeString(txHash)
				table.Put(hash, utxos.Serialize())

			}

			for _, txInput := range txInputs {
				utxosBytes := table.Get(txInput.TxHash)
				utxos := DeserializeUTXOS(utxosBytes)
				tmpUtxos := &UTXOS{}
				deleteFlag := false
				for _, utxo := range utxos.UTXOs {
					if txInput.TxOutIndex == utxo.Index && bytes.Compare(utxo.Output.Sha256Ripemd160HashPubkey, Sha256Ripemd160Hash(txInput.Pubkey)) == 0 {
						deleteFlag = true
					} else {
						tmpUtxos.UTXOs = append(tmpUtxos.UTXOs, utxo)

					}
				}
				if deleteFlag == true {
					table.Delete(txInput.TxHash)
					table.Put(txInput.TxHash, tmpUtxos.Serialize())

				}
			}

		} else {
			log.Panic("no utxo table!")
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}
