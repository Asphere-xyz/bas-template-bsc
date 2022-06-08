package vm

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"strings"
)

// verifyParliaBlock implements precompile for fas verify of parlia block
type verifyParliaBlock struct{}

func NewVerifyParliaBlock() *verifyParliaBlock {
	return &verifyParliaBlock{}
}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *verifyParliaBlock) RequiredGas([]byte) uint64 {
	return params.VerifyParliaBlockGas
}

var errBadParliaBlock = fmt.Errorf("bad parlia block")
var errMalformedInputParams = fmt.Errorf("malformed input params")

func mustNewArguments(types ...string) (result abi.Arguments) {
	var err error
	for _, t := range types {
		var typ abi.Type
		items := strings.Split(t, " ")
		var name string
		if len(items) == 2 {
			name = items[1]
		} else {
			name = items[0]
		}
		typ, err = abi.NewType(items[0], items[0], nil)
		if err != nil {
			panic(err)
		}
		result = append(result, abi.Argument{Type: typ, Name: name})
	}
	return result
}

var verifyParliaBlockInput = mustNewArguments(
	"uint256 chainId",
	"bytes blockProof",
	"uint32 epochInterval",
)

var verifyParliaBlockOutput = mustNewArguments(
	"bytes32 blockHash",
	"uint64 blockNumber",
	"address signer",
	"address[] validators",
	"bytes32 parentHash",
)

func (c *verifyParliaBlock) Run(input []byte) (result []byte, err error) {
	var chainId *big.Int
	var blockProof []byte
	var epochInterval uint32
	{
		if input == nil || len(input) > 65536 {
			return nil, errBadParliaBlock
		}
		items, err := verifyParliaBlockInput.UnpackValues(input)
		if err != nil {
			return nil, err
		} else if len(items) != 3 {
			return nil, errMalformedInputParams
		}
		var ok bool
		chainId, ok = items[0].(*big.Int)
		if !ok {
			return nil, errMalformedInputParams
		}
		blockProof, ok = items[1].([]byte)
		if !ok {
			return nil, errMalformedInputParams
		}
		epochInterval, ok = items[2].(uint32)
		if !ok {
			return nil, errMalformedInputParams
		}
	}
	header := &types.Header{}
	if err := rlp.Decode(bytes.NewReader(blockProof), header); err != nil {
		return nil, err
	}
	var signer common.Address
	if header.Number.Uint64() != 0 {
		signer, err = recoverParliaBlockSigner(header, chainId)
		if err != nil {
			return nil, err
		}
	}
	var validators []common.Address
	if header.Number.Uint64()%uint64(epochInterval) == 0 {
		validators, err = extractParliaValidators(header)
		if err != nil {
			return nil, err
		}
	}
	return verifyParliaBlockOutput.Pack(
		// block hash
		header.Hash(),
		// block number
		header.Number.Uint64(),
		// signing data
		signer,
		// validators
		validators,
		// parent hash
		header.ParentHash,
	)
}

const (
	extraVanity = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal   = 65 // Fixed number of extra-data suffix bytes reserved for signer seal
)

func recoverParliaBlockSigner(header *types.Header, chainId *big.Int) (signer common.Address, err error) {
	if len(header.Extra) < extraSeal {
		return signer, errBadParliaBlock
	}
	signature := header.Extra[len(header.Extra)-extraSeal:]
	b := new(bytes.Buffer)
	err = rlp.Encode(b, []interface{}{
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
	signingData := b.Bytes()
	publicKey, err := crypto.Ecrecover(crypto.Keccak256(signingData), signature)
	if err != nil {
		return signer, err
	}
	copy(signer[:], crypto.Keccak256(publicKey[1:])[12:])
	return signer, nil
}

func extractParliaValidators(header *types.Header) ([]common.Address, error) {
	validatorBytes := header.Extra[extraVanity : len(header.Extra)-extraSeal]
	if len(validatorBytes)%common.AddressLength != 0 {
		return nil, errBadParliaBlock
	}
	n := len(validatorBytes) / common.AddressLength
	result := make([]common.Address, n)
	for i := 0; i < n; i++ {
		address := make([]byte, common.AddressLength)
		copy(address, validatorBytes[i*common.AddressLength:(i+1)*common.AddressLength])
		result[i] = common.BytesToAddress(address)
	}
	return result, nil
}
