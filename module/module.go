package module

import (
	"context"
	"cosmossdk.io/core/appmodule"
	"encoding/json"
	"fmt"
	auctiontypes "github.com/fatal-fruit/auction/types"
	"github.com/spf13/cobra"

	abci "github.com/cometbft/cometbft/abci/types"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	auctionabci "github.com/fatal-fruit/auction/abci"
	auctioncli "github.com/fatal-fruit/auction/client"
	"github.com/fatal-fruit/auction/keeper"
)

const ConsensusVersion = 1

var (
	_ module.AppModuleBasic = AppModule{}
	//_ module.HasGenesis     = AppModule{}
	_ module.HasServices = AppModule{}

	_ appmodule.AppModule     = AppModule{}
	_ appmodule.HasEndBlocker = AppModule{}
)

type AppModule struct {
	cdc    codec.Codec
	keeper keeper.Keeper
}

func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
	return AppModule{
		cdc:    cdc,
		keeper: keeper,
	}
}

func NewAppModuleBasic(m AppModule) module.AppModuleBasic {
	return module.CoreAppModuleBasicAdaptor(m.Name(), m)
}

func (AppModule) Name() string { return auctiontypes.ModuleName }

func (AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	auctiontypes.RegisterLegacyAminoCodec(cdc)
}

func (AppModule) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	if err := auctiontypes.RegisterQueryHandlerClient(context.Background(), mux, auctiontypes.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

func (AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	auctiontypes.RegisterInterfaces(registry)
}

func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

func (am AppModule) EndBlock(ctx context.Context) error {
	return auctionabci.EndBlocker(ctx, am.keeper, am.keeper.Logger())
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	auctiontypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	auctiontypes.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))

	// Register in place module state migration migrations
	// m := keeper.NewMigrator(am.keeper)
	// if err := cfg.RegisterMigration(ns.ModuleName, 1, m.Migrate1to2); err != nil {
	// 	panic(fmt.Sprintf("failed to migrate x/%s from version 1 to 2: %v", ns.ModuleName, err))
	// }
}

func (AppModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(auctiontypes.NewGenesisState())
}

func (AppModule) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data auctiontypes.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", auctiontypes.ModuleName, err)
	}

	return data.Validate()
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState auctiontypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	if err := am.keeper.InitGenesis(ctx, &genesisState); err != nil {
		panic(fmt.Sprintf("failed to initialize %s genesis state: %v", auctiontypes.ModuleName, err))
	}

	return nil
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs, err := am.keeper.ExportGenesis(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to export %s genesis state: %v", auctiontypes.ModuleName, err))
	}

	return cdc.MustMarshalJSON(gs)
}

func (AppModule) GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        auctiontypes.ModuleName,
		Short:                      "auction transaction subcommands",
		Long:                       "Commands for creating and executing auctions",
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		auctioncli.NewAuctionCmd(),
		auctioncli.StartAuctionCmd(),
		auctioncli.BidCmd(),
		auctioncli.ExecuteAuctionCmd(),
		auctioncli.CancelAuctionCmd(),
	)
	return cmd
}
