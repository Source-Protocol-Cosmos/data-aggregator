package main

import (
	"os"

	"github.com/Source-Protocol-Cosmos/juno/v3/cmd/parse/types"

	wasmapp "github.com/CosmWasm/wasmd/app"
	"github.com/Source-Protocol-Cosmos/juno/v3/cmd"
	"github.com/Source-Protocol-Cosmos/juno/v3/modules/messages"
	"github.com/Source-Protocol-Cosmos/juno/v3/modules/registrar"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/std"
)
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeTestEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	wasmapp.ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	wasmapp.ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
func main() {
	// JunoConfig the runner
	config := cmd.NewConfig("juno").
		WithParseConfig(types.NewConfig().
			WithRegistrar(registrar.NewDefaultRegistrar(
				messages.CosmosMessageAddressesParser,
			)).WithEncodingConfigBuilder(MakeEncodingConfig),
		)

	// Run the commands and panic on any error
	exec := cmd.BuildDefaultExecutor(config)
	err := exec.Execute()
	if err != nil {
		os.Exit(1)
	}
}
