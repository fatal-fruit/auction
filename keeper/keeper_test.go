package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	auctiontestutil "github.com/fatal-fruit/auction/testutil"
	auctiontypes "github.com/fatal-fruit/auction/types"
	"github.com/stretchr/testify/require"
)

func TestProcessActiveAuctions(t *testing.T) {
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
	}
	err = f.K.Auctions.Set(f.Ctx, id, auction)
	require.NoError(err)
	err = f.K.ActiveAuctions.Set(f.Ctx, id)
	require.NoError(err)

	err = f.K.ProcessActiveAuctions(f.Ctx)
	require.NoError(err)
	isActive, err := f.K.ActiveAuctions.Has(f.Ctx, id)
	require.False(isActive)
	require.NoError(err)
	isExpired, err := f.K.ExpiredAuctions.Has(f.Ctx, id)
	require.NoError(err)
	require.True(isExpired)
}

func TestProcessExpiredAuctions(t *testing.T) {
	f := auctiontestutil.InitFixture(t)
	require := require.New(t)

	id1, err := f.K.IDs.Next(f.Ctx)
	require.NoError(err)
	id2, err := f.K.IDs.Next(f.Ctx)
	require.NoError(err)

	auctions := []auctiontypes.ReserveAuction{
		{
			Id:           id1,
			Status:       auctiontypes.ACTIVE,
			Owner:        f.Addrs[0].String(),
			AuctionType:  auctiontypes.RESERVE,
			ReservePrice: sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000),
			StartTime:    time.Now().Add(-30 * time.Second),
			EndTime:      time.Now().Add(-1 * time.Second),
			LastPrice:    sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1100),
			Bids: []*auctiontypes.Bid{
				{
					AuctionId: id1,
					Bidder:    f.Addrs[1].String(),
					BidPrice:  sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1100),
					Timestamp: time.Now(),
				},
			},
		},
		{
			Id:           id2,
			Status:       auctiontypes.ACTIVE,
			Owner:        f.Addrs[0].String(),
			AuctionType:  auctiontypes.RESERVE,
			ReservePrice: sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000),
			StartTime:    time.Now().Add(-30 * time.Second),
			EndTime:      time.Now().Add(-1 * time.Second),
			Bids:         []*auctiontypes.Bid{},
		},
	}

	for _, a := range auctions {
		err = f.K.Auctions.Set(f.Ctx, a.GetId(), a)
		require.NoError(err)
		err = f.K.ExpiredAuctions.Set(f.Ctx, a.GetId())
		require.NoError(err)
	}

	for _, a := range auctions {
		isExpired, err := f.K.ExpiredAuctions.Has(f.Ctx, a.GetId())
		require.True(isExpired)
		require.NoError(err)
	}

	err = f.K.ProcessExpiredAuctions(f.Ctx)
	require.NoError(err)

	for _, a := range auctions {
		isExpired, err := f.K.ExpiredAuctions.Has(f.Ctx, a.GetId())
		require.False(isExpired)
		require.NoError(err)
	}

	isPending, err := f.K.PendingAuctions.Has(f.Ctx, id1)
	require.True(isPending)
	require.NoError(err)

	isCancelled, err := f.K.CancelledAuctions.Has(f.Ctx, id2)
	require.True(isCancelled)
	require.NoError(err)
}

func TestGetAllAuctions(t *testing.T) {
	f := auctiontestutil.InitFixture(t)
	require := require.New(t)

	id1, err := f.K.IDs.Next(f.Ctx)
	require.NoError(err)
	id2, err := f.K.IDs.Next(f.Ctx)
	require.NoError(err)

	auctions := []auctiontypes.ReserveAuction{
		{
			Id:           id1,
			Status:       auctiontypes.ACTIVE,
			Owner:        f.Addrs[0].String(),
			AuctionType:  auctiontypes.RESERVE,
			ReservePrice: sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000),
			StartTime:    time.Now().Add(-30 * time.Second),
			EndTime:      time.Now().Add(-1 * time.Second),
			LastPrice:    sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1100),
			Bids: []*auctiontypes.Bid{
				{
					AuctionId: id1,
					Bidder:    f.Addrs[1].String(),
					BidPrice:  sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1100),
					Timestamp: time.Now(),
				},
			},
		},
		{
			Id:           id2,
			Status:       auctiontypes.ACTIVE,
			Owner:        f.Addrs[0].String(),
			AuctionType:  auctiontypes.RESERVE,
			ReservePrice: sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000),
			StartTime:    time.Now().Add(-30 * time.Second),
			EndTime:      time.Now().Add(-1 * time.Second),
			Bids:         []*auctiontypes.Bid{},
		},
	}

	for _, a := range auctions {
		err = f.K.Auctions.Set(f.Ctx, a.GetId(), a)
		require.NoError(err)
		err = f.K.ExpiredAuctions.Set(f.Ctx, a.GetId())
		require.NoError(err)
	}

	auctions = f.K.GetAllAuctions(f.Ctx)
	require.Equal(2, len(auctions))
}

func TestPurgeCancelledAuctions(t *testing.T) {
	f := auctiontestutil.InitFixture(t)
	require := require.New(t)

	for i := 0; i < 3; i++ {
		id, err := f.K.IDs.Next(f.Ctx)
		require.NoError(err)
		err = f.K.CancelledAuctions.Set(f.Ctx, id)
		require.NoError(err)
	}

	err := f.K.PurgeCancelledAuctions(f.Ctx)
	require.NoError(err)

	cancelledAuctions := f.K.GetCancelledAuctions(f.Ctx)
	require.Empty(cancelledAuctions)
}

func TestGetCancelledAuctions(t *testing.T) {
	f := auctiontestutil.InitFixture(t)
	require := require.New(t)

	numAuctions := 3
	for i := 0; i < numAuctions; i++ {
		id, err := f.K.IDs.Next(f.Ctx)
		require.NoError(err)

		auction := auctiontypes.ReserveAuction{Id: id}
		err = f.K.Auctions.Set(f.Ctx, id, auction)
		require.NoError(err)

		err = f.K.CancelAuction(f.Ctx, id)
		require.NoError(err)
	}
}
