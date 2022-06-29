package wasm

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/Source-Protocol-Cosmos/data-aggregator/v3/database"
	"github.com/Source-Protocol-Cosmos/data-aggregator/v3/modules"
)

var (
	_ modules.Module        = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/profiles module handler
type Module struct {
	db     database.Database
	client wasmtypes.QueryClient
}

// NewModule allows to build a new Module instance
func NewModule(db database.Database, client wasmtypes.QueryClient) *Module {
	return &Module{
		db:     db,
		client: client,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "wasm"
}
