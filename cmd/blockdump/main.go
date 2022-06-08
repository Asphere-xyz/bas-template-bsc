package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/parlia"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	extraVanity      = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal        = 65 // Fixed number of extra-data suffix bytes reserved for signer seal
	nextForkHashSize = 4  // Fixed number of extra-data suffix bytes reserved for nextForkHash.
)

func encodeSigHeader(w io.Writer, header *types.Header, chainId *big.Int) {
	err := rlp.Encode(w, []interface{}{
		chainId,
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-65], // this will panic if extra is too short, should check before calling encodeSigHeader
		header.MixDigest,
		header.Nonce,
	})
	if err != nil {
		panic("can't encode: " + err.Error())
	}
}

func SealHash(header *types.Header, chainId *big.Int) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	encodeSigHeader(hasher, header, chainId)
	hasher.Sum(hash[:0])
	return hash
}

func ecrecover(header *types.Header, chainId *big.Int) (common.Address, error) {
	// Retrieve the signature from the header extra-data
	if len(header.Extra) < extraSeal {
		return common.Address{}, fmt.Errorf("missing signature")
	}
	signature := header.Extra[len(header.Extra)-extraSeal:]

	// Recover the public key and the Ethereum address
	pubkey, err := crypto.Ecrecover(SealHash(header, chainId).Bytes(), signature)
	if err != nil {
		return common.Address{}, err
	}
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

	return signer, nil
}

func main() {
	eth, err := ethclient.Dial("https://rpc.ankr.com/bsc")
	if err != nil {
		panic(err)
	}
	//confirmations := 12
	//for i := 0; i < confirmations; i++ {
	//	block, err := eth.BlockByNumber(context.Background(), big.NewInt(int64(15946200+i)))
	//	if err != nil {
	//		panic(err)
	//	}
	//	payload, _ := rlp.EncodeToBytes(block.Header())
	//	println(hexutil.Encode(payload))
	//}
	block, err := eth.BlockByNumber(context.Background(), big.NewInt(int64(13082000)))
	if err != nil {
		panic(err)
	}
	println("-----------BLOCK HEADER JSON-----------------")
	json, _ := block.Header().MarshalJSON()
	println(string(json))
	println("------------BLOCK HEADER---------------------")
	payload, err := rlp.EncodeToBytes(block.Header())
	println(hexutil.Encode(payload)[2:])
	println("----------EXTRA DATA SHOULD BE---------------")
	println(hexutil.Encode(block.Header().Extra[:len(block.Header().Extra)-65])[2:])
	println("----------SIGNING DATA-----------------------")
	signingData := parlia.ParliaRLP(block.Header(), big.NewInt(56))
	println(hexutil.Encode(signingData)[2:])
	println()
	println()
	addr, err := ecrecover(block.Header(), big.NewInt(56))
	if err != nil {
		panic(err)
	}
	println(addr.Hex())
}
