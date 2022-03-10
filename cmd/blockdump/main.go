package main

import (
	"bytes"
	"context"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/parlia"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

func main() {
	eth, err := ethclient.Dial("https://rpc.ankr.com/bsc")
	if err != nil {
		panic(err)
	}
	block, err := eth.BlockByNumber(context.Background(), big.NewInt(13082000))
	if err != nil {
		panic(err)
	}
	b := new(bytes.Buffer)
	rlp.Encode(b, []interface{}{
		big.NewInt(65),
	})
	println(hexutil.Encode(b.Bytes()))
	println("---------------------------------------------")
	json, _ := block.Header().MarshalJSON()
	println(string(json))
	println("---------------------------------------------")
	payload, err := rlp.EncodeToBytes(block.Header())
	println(hexutil.Encode(payload))
	println("---------------------------------------------")
	println(hexutil.Encode(block.Header().Extra[:len(block.Header().Extra)-65]))
	println("---------------------------------------------")
	signingData := parlia.ParliaRLP(block.Header(), big.NewInt(65))
	println(hexutil.Encode(signingData))
}
