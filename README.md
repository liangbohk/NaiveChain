# NaiveChain
A blockchain demo refering to online courses and ethereum/go-ethereum

### build
In project directory and execute
```
# go build .
# ls
cmd  core  LICENSE  main.go  NaiveChain  README.md
```
NaiveChain is the executable file

### Execute `NaiveChain` without options to get help
```
# ./NaiveChain
Usage:
        getaddresslist -- output all address
        createwallet -- create a wallet
        createblockchain -address ADDRESS -- create a blockchain
        send -from FROM -to TO -amount AMOUNT -- send value by transaction
        getbalance -address ADDRESS -- get balance of an address
        printchain -- print the blockchain
        test -- test

```

### Some options
#### Create a wallet (with an address)
```
# ./NaiveChain createwallet
1NrQM8LXD5CRworrmKZPs4Pc6ZSbEzCPTZ
```
#### Create a blockchain and send a transaction (watch out the form)
```
# ./NaiveChain createblockchain -address 1NrQM8LXD5CRworrmKZPs4Pc6ZSbEzCPTZ
./NaiveChain send -from '["1NrQM8LXD5CRworrmKZPs4Pc6ZSbEzCPTZ","1NrQM8LXD5CRworrmKZPs4Pc6ZSbEzCPTZ"]' -to '["1Jf76f7tHBpH2W4o8gGo4brAxWRe25vi86","1FXTcPQB3VT5pwA3tUX2D4EM6an4CnGn7Y"]' -amount '["5","3"]'
```
#### Try other options and test it!

