package types

import (
	wasmvmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

// Contract type
type Contract struct {
	*wasmvmtypes.ContractInfo
	Address     string
	CreatedTime string
}

// NewContract instance
func NewContract(contract *wasmvmtypes.ContractInfo, address string, created string) Contract {
	return Contract{
		ContractInfo: contract,
		Address:      address,
		CreatedTime:  created,
	}
}
