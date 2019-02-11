package core

import "log"

//handle command from other nodes

//handle version info
func handleVersion(request []byte, blc *Blockchain) {
	//get version
	v := DeserializeVersion(request[COMMANDLENGTH:])

	//compare current blockchain height and foreign blockchain height
	curHeight, err := blc.GetBlockchainHeight()
	if err != nil {
		log.Panic(err)
	}
	foreignHeight := v.BaseHeight
	if curHeight > foreignHeight {
		sendVersion(v.AddressFrom, blc)
	} else if curHeight < foreignHeight {
		//request missing blocks
		sendGetBlocks(v.AddressFrom)
	}

}

func handleAddr(request []byte, blc *Blockchain) {

}

func handleBlock(request []byte, blc *Blockchain) {
	gbs := DeserializeGetBlocks(request[COMMANDLENGTH:])
	blcHashes := blc.GetBlockHashes()
	sendInv(gbs.AddrFrom, "block", BLOCK_COMMAND, blcHashes)
}

func handleInv(request []byte, blc *Blockchain) {

}

func handleTx(request []byte, blc *Blockchain) {

}

func handleGetBlocks(request []byte, blc *Blockchain) {

}

func handleGetData(request []byte, blc *Blockchain) {

}
