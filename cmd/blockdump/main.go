package main

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

func main() {
	eth, err := ethclient.Dial("https://rpc.ankr.com/bsc")
	if err != nil {
		panic(err)
	}
	confirmations := 12
	for i := 0; i < confirmations; i++ {
		block, err := eth.BlockByNumber(context.Background(), big.NewInt(int64(15946200+i)))
		if err != nil {
			panic(err)
		}
		payload, _ := rlp.EncodeToBytes(block.Header())
		println(hexutil.Encode(payload))
	}
	//println("-----------BLOCK HEADER JSON-----------------")
	//json, _ := block.Header().MarshalJSON()
	//println(string(json))
	//println("------------BLOCK HEADER---------------------")
	//payload, err := rlp.EncodeToBytes(block.Header())
	//println(hexutil.Encode(payload)[2:])
	//println("----------EXTRA DATA SHOULD BE---------------")
	//println(hexutil.Encode(block.Header().Extra[:len(block.Header().Extra)-65])[2:])
	//println("----------SIGNING DATA-----------------------")
	//signingData := parlia.ParliaRLP(block.Header(), big.NewInt(56))
	//println(hexutil.Encode(signingData)[2:])
}
