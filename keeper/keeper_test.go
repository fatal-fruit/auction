package keeper_test

import (
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/fatal-fruit/auction/keeper"
	auctiontestutil "github.com/fatal-fruit/auction/testutil"
	auctiontypes "github.com/fatal-fruit/auction/types"
	"go.uber.org/mock/gomock"
	"testing"
)

type testFixture struct {
	ctx         sdk.Context
	k           keeper.Keeper
	msgServer   auctiontypes.MsgServer
	queryServer auctiontypes.QueryServer

	mockAcctKeeper    *auctiontestutil.MockAccountKeeper
	mockBankKeeper    *auctiontestutil.MockBankKeeper
	mockEscrowService *auctiontestutil.MockEscrowService

	addrs      []sdk.AccAddress
	modAccount *authtypes.ModuleAccount
	modAddr    sdk.AccAddress
}

func initFixture(t *testing.T) *testFixture {
	encConfig := moduletestutil.MakeTestEncodingConfig()
	storeKey := storetypes.NewKVStoreKey(auctiontypes.ModuleName)
	testCtx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("t_test"))
	storeService := runtime.NewKVStoreService(storeKey)
	addrs := simtestutil.CreateIncrementalAccounts(3)
	authority := authtypes.NewModuleAddress("gov")
	auctionModAddr := authtypes.NewModuleAddress(auctiontypes.ModuleName)
	auctionAcct := authtypes.NewEmptyModuleAccount(auctiontypes.ModuleName, authtypes.Minter)

	ctrl := gomock.NewController(t)
	mockAcctKeeper := auctiontestutil.NewMockAccountKeeper(ctrl)
	mockBankKeeper := auctiontestutil.NewMockBankKeeper(ctrl)
	mockEscrowService := auctiontestutil.NewMockEscrowService(ctrl)

	k := keeper.NewKeeper(
		encConfig.Codec,
		addresscodec.NewBech32Codec("cosmos"),
		storeService,
		authority.String(),
		mockAcctKeeper,
		mockBankKeeper,
		mockEscrowService,
		sdk.DefaultBondDenom,
	)
	err := k.InitGenesis(testCtx.Ctx, auctiontypes.NewGenesisState())
	if err != nil {
		panic(err)
	}

	return &testFixture{
		ctx:               testCtx.Ctx,
		k:                 k,
		msgServer:         keeper.NewMsgServerImpl(k),
		queryServer:       keeper.NewQueryServerImpl(k),
		addrs:             addrs,
		modAccount:        auctionAcct,
		modAddr:           auctionModAddr,
		mockAcctKeeper:    mockAcctKeeper,
		mockBankKeeper:    mockBankKeeper,
		mockEscrowService: mockEscrowService,
	}
}
