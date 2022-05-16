package systemcontract

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

type EvmHook interface {
	RequiredGas(input []byte) uint64
	Run(input []byte) ([]byte, error)
}

type StateDB interface {
	GetCodeHash(common.Address) common.Hash
	GetCode(common.Address) []byte
	SetCode(common.Address, []byte)
	GetCodeSize(common.Address) int
}

type EVM interface {
	CreateWithAddress(caller common.Address, code []byte, gas uint64, value *big.Int, contractAddr common.Address) (ret []byte, leftOverGas uint64, err error)
}

type EvmHookContext struct {
	CallerAddress common.Address
	StateDb       StateDB
	Evm           EVM
	ChainConfig   *params.ChainConfig
	ChainRules    params.Rules
}
