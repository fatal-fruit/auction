package keeper_test

import (
	"go.uber.org/mock/gomock"
	"testing"

	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/fatal-fruit/auction/keeper"
	auctiontypes "github.com/fatal-fruit/auction/types"

	auctiontestutil "github.com/fatal-fruit/auction/testutil"
)

type testFixture struct {
	ctx         sdk.Context
	k           keeper.Keeper
	msgServer   auctiontypes.MsgServer
	queryServer auctiontypes.QueryServer

	addrs []sdk.AccAddress
}

func initFixture(t *testing.T) *testFixture {
	encCfg := moduletestutil.MakeTestEncodingConfig()
	key := storetypes.NewKVStoreKey(auctiontypes.ModuleName)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	storeService := runtime.NewKVStoreService(key)
	addrs := simtestutil.CreateIncrementalAccounts(3)
	// gomock initializations
	ctrl := gomock.NewController(t)
	bankKeeper := auctiontestutil.NewMockBankKeeper(ctrl)
	defaultDenom := "stake"
	k := keeper.NewKeeper(
		encCfg.Codec,
		addresscodec.NewBech32Codec("cosmos"),
		storeService,
		addrs[0].String(),
		bankKeeper,
		defaultDenom,
	)

	err := k.InitGenesis(testCtx.Ctx, auctiontypes.NewGenesisState())
	if err != nil {
		panic(err)
	}

	return &testFixture{
		ctx:         testCtx.Ctx,
		k:           k,
		msgServer:   keeper.NewMsgServerImpl(k),
		queryServer: keeper.NewQueryServerImpl(k),
		addrs:       addrs,
	}
}
