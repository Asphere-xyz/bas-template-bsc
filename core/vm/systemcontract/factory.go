package systemcontract

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/systemcontract"
)

func CreateEvmHook(address common.Address, context EvmHookContext) EvmHook {
	if address == systemcontract.EvmHookRuntimeUpgradeAddress {
		return &evmHookRuntimeUpgrade{context: context}
	}
	return nil
}
