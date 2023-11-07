package abci

import (
	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auctiontestutil "github.com/fatal-fruit/auction/testutil"
	auctiontypes "github.com/fatal-fruit/auction/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestEndBlocker_ActiveToCancelled(t *testing.T) {
	f := auctiontestutil.InitFixture(t)
	require := require.New(t)

	id, err := f.K.IDs.Next(f.Ctx)
	require.NoError(err)
	auction := auctiontypes.ReserveAuction{
		Id:           id,
		Status:       auctiontypes.ACTIVE,
		Owner:        f.Addrs[0].String(),
		AuctionType:  auctiontypes.RESERVE,
		ReservePrice: sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000),
		StartTime:    time.Now().Add(-30 * time.Second),
		EndTime:      time.Now().Add(-1 * time.Second),
		Bids:         []*auctiontypes.Bid{},
		Strategy: &auctiontypes.SettleStrategy{
			StrategyType:          auctiontypes.SETTLE,
			EscrowContractId:      1,
			EscrowContractAddress: f.Addrs[2].String(),
		},
	}
	err = f.K.Auctions.Set(f.Ctx, id, auction)
	require.NoError(err)
	err = f.K.ActiveAuctions.Set(f.Ctx, id)
	require.NoError(err)

	logger := log.NewNopLogger()
	err = EndBlocker(f.Ctx, f.K, logger)
	require.NoError(err)

	inActive, err := f.K.ActiveAuctions.Has(f.Ctx, id)
	require.NoError(err)
	require.False(inActive)

	inExpired, err := f.K.ExpiredAuctions.Has(f.Ctx, id)
	require.NoError(err)
	require.False(inExpired)

	inPending, err := f.K.PendingAuctions.Has(f.Ctx, id)
	require.NoError(err)
	require.False(inPending)

	inCancelled, err := f.K.CancelledAuctions.Has(f.Ctx, id)
	require.NoError(err)
	require.True(inCancelled)
}

func TestEndBlocker_ActiveToPending(t *testing.T) {
	f := auctiontestutil.InitFixture(t)
	require := require.New(t)

	id, err := f.K.IDs.Next(f.Ctx)
	require.NoError(err)
	auction := auctiontypes.ReserveAuction{
		Id:           id,
		Status:       auctiontypes.ACTIVE,
		Owner:        f.Addrs[0].String(),
		AuctionType:  auctiontypes.RESERVE,
		ReservePrice: sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000),
		StartTime:    time.Now().Add(-30 * time.Second),
		EndTime:      time.Now().Add(-1 * time.Second),
		LastPrice:    sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1100),
		Bids: []*auctiontypes.Bid{
			{
				AuctionId: id,
				Bidder:    f.Addrs[1].String(),
				BidPrice:  sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1100),
				Timestamp: time.Now(),
			},
		},
		Strategy: &auctiontypes.SettleStrategy{
			StrategyType:          auctiontypes.SETTLE,
			EscrowContractId:      1,
			EscrowContractAddress: f.Addrs[2].String(),
		},
	}
	err = f.K.Auctions.Set(f.Ctx, id, auction)
	require.NoError(err)
	err = f.K.ActiveAuctions.Set(f.Ctx, id)
	require.NoError(err)

	logger := log.NewNopLogger()
	err = EndBlocker(f.Ctx, f.K, logger)
	require.NoError(err)

	inActive, err := f.K.ActiveAuctions.Has(f.Ctx, id)
	require.NoError(err)
	require.False(inActive)

	inExpired, err := f.K.ExpiredAuctions.Has(f.Ctx, id)
	require.NoError(err)
	require.False(inExpired)

	inPending, err := f.K.PendingAuctions.Has(f.Ctx, id)
	require.NoError(err)
	require.True(inPending)

	inCancelled, err := f.K.CancelledAuctions.Has(f.Ctx, id)
	require.NoError(err)
	require.False(inCancelled)
}
