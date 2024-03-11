package testutil

import (
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/fatal-fruit/auction/keeper"
	auctiontypes "github.com/fatal-fruit/auction/types"
	"go.uber.org/mock/gomock"
	"testing"
)

type TestFixture struct {
	Ctx         sdk.Context
	EnCfg       moduletestutil.TestEncodingConfig
	K           keeper.Keeper
	MsgServer   auctiontypes.MsgServer
	QueryServer auctiontypes.QueryServer
	Resolver    auctiontypes.AuctionResolver

	MockAcctKeeper    *MockAccountKeeper
	MockBankKeeper    *MockBankKeeper
	MockEscrowService *MockEscrowService

	Addrs              []sdk.AccAddress
	ModAccount         *authtypes.ModuleAccount
	ModAddr            sdk.AccAddress
	Logger             log.Logger
	ReserveAuctionType string
}

func InitFixture(t *testing.T) *TestFixture {
	encConfig := moduletestutil.MakeTestEncodingConfig()
	encConfig.InterfaceRegistry.RegisterInterface(
		"fatal_fruit.auction.v1.AuctionMetadata",
		(*auctiontypes.AuctionMetadata)(nil),
		&auctiontypes.ReserveAuctionMetadata{},
	)
	encConfig.InterfaceRegistry.RegisterInterface(
		"fatal_fruit.auction.v1.Auction",
		(*auctiontypes.Auction)(nil),
		&auctiontypes.ReserveAuction{},
	)
	storeKey := storetypes.NewKVStoreKey(auctiontypes.ModuleName)
	testCtx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("t_test"))
	storeService := runtime.NewKVStoreService(storeKey)
	addrs := simtestutil.CreateIncrementalAccounts(3)
	authority := authtypes.NewModuleAddress("gov")
	auctionModAddr := authtypes.NewModuleAddress(auctiontypes.ModuleName)
	auctionAcct := authtypes.NewEmptyModuleAccount(auctiontypes.ModuleName, authtypes.Minter)

	ctrl := gomock.NewController(t)
	mockAcctKeeper := NewMockAccountKeeper(ctrl)
	mockBankKeeper := NewMockBankKeeper(ctrl)
	mockEscrowService := NewMockEscrowService(ctrl)

	resolver := auctiontypes.NewResolver()
	handler := auctiontypes.NewReserveAuctionHandler(mockEscrowService, mockBankKeeper)
	reserveType := sdk.MsgTypeURL(&auctiontypes.ReserveAuction{})
	resolver.AddType(reserveType, handler)
	resolver.Seal()

	k := keeper.NewKeeper(
		encConfig.Codec,
		addresscodec.NewBech32Codec("cosmos"),
		storeService,
		authority.String(),
		mockAcctKeeper,
		mockBankKeeper,
		sdk.DefaultBondDenom,
		log.NewNopLogger(),
	)
	k.SetAuctionTypesResolver(resolver)
	err := k.InitGenesis(testCtx.Ctx, auctiontypes.NewGenesisState())
	if err != nil {
		panic(err)
	}

	return &TestFixture{
		Ctx:                testCtx.Ctx,
		EnCfg:              encConfig,
		K:                  k,
		MsgServer:          keeper.NewMsgServerImpl(k),
		QueryServer:        keeper.NewQueryServerImpl(k),
		Addrs:              addrs,
		ModAccount:         auctionAcct,
		ModAddr:            auctionModAddr,
		MockAcctKeeper:     mockAcctKeeper,
		MockBankKeeper:     mockBankKeeper,
		MockEscrowService:  mockEscrowService,
		Resolver:           resolver,
		ReserveAuctionType: sdk.MsgTypeURL(&auctiontypes.ReserveAuction{}),
	}
}
