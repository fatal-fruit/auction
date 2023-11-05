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
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
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

func TestProcessActiveAuctions(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)

	id, err := f.k.IDs.Next(f.ctx)
	require.NoError(err)
	auction := auctiontypes.ReserveAuction{
		Id:           id,
		Status:       auctiontypes.ACTIVE,
		Owner:        f.addrs[0].String(),
		AuctionType:  auctiontypes.RESERVE,
		ReservePrice: sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000),
		StartTime:    time.Now().Add(-30 * time.Second),
		EndTime:      time.Now().Add(-1 * time.Second),
		Bids:         []*auctiontypes.Bid{},
	}
	err = f.k.Auctions.Set(f.ctx, id, auction)
	require.NoError(err)
	err = f.k.ActiveAuctions.Set(f.ctx, id)
	require.NoError(err)

	f.k.ProcessActiveAuctions(f.ctx)
	isActive, err := f.k.ActiveAuctions.Has(f.ctx, id)
	require.False(isActive)
	require.NoError(err)
	isExpired, err := f.k.ExpiredAuctions.Has(f.ctx, id)
	require.NoError(err)
	require.True(isExpired)
}

func TestProcessExpiredAuctions(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)

	id1, err := f.k.IDs.Next(f.ctx)
	require.NoError(err)
	id2, err := f.k.IDs.Next(f.ctx)
	require.NoError(err)

	auctions := []auctiontypes.ReserveAuction{
		{
			Id:           id1,
			Status:       auctiontypes.ACTIVE,
			Owner:        f.addrs[0].String(),
			AuctionType:  auctiontypes.RESERVE,
			ReservePrice: sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000),
			StartTime:    time.Now().Add(-30 * time.Second),
			EndTime:      time.Now().Add(-1 * time.Second),
			LastPrice:    sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1100),
			Bids: []*auctiontypes.Bid{
				{
					AuctionId: id1,
					Bidder:    f.addrs[1].String(),
					BidPrice:  sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1100),
					Timestamp: time.Now(),
				},
			},
		},
		{
			Id:           id2,
			Status:       auctiontypes.ACTIVE,
			Owner:        f.addrs[0].String(),
			AuctionType:  auctiontypes.RESERVE,
			ReservePrice: sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000),
			StartTime:    time.Now().Add(-30 * time.Second),
			EndTime:      time.Now().Add(-1 * time.Second),
			Bids:         []*auctiontypes.Bid{},
		},
	}

	for _, a := range auctions {
		err = f.k.Auctions.Set(f.ctx, a.GetId(), a)
		require.NoError(err)
		err = f.k.ExpiredAuctions.Set(f.ctx, a.GetId())
		require.NoError(err)
	}

	for _, a := range auctions {
		isExpired, err := f.k.ExpiredAuctions.Has(f.ctx, a.GetId())
		require.True(isExpired)
		require.NoError(err)
	}

	f.k.ProcessExpiredAuctions(f.ctx)

	for _, a := range auctions {
		isExpired, err := f.k.ExpiredAuctions.Has(f.ctx, a.GetId())
		require.False(isExpired)
		require.NoError(err)
	}

	isPending, err := f.k.PendingAuctions.Has(f.ctx, id1)
	require.True(isPending)
	require.NoError(err)

	isCancelled, err := f.k.CancelledAuctions.Has(f.ctx, id2)
	require.True(isCancelled)
	require.NoError(err)
}
