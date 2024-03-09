package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/fatal-fruit/auction/keeper"
	auctiontestutil "github.com/fatal-fruit/auction/testutil"
	auctiontypes "github.com/fatal-fruit/auction/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestQueryAuction(t *testing.T) {
	f := auctiontestutil.InitFixture(t)
	require := require.New(t)
	contractId := uint64(0)

	testCases := []struct {
		name      string
		req       auctiontypes.MsgNewAuction
		expErr    bool
		setupTest func(fixture *auctiontestutil.TestFixture) struct {
			res *auctiontypes.MsgNewAuctionResponse
		}
	}{
		{
			name: "retrieve valid auction",
			req: auctiontypes.MsgNewAuction{
				Owner:        f.Addrs[0].String(),
				Deposit:      sdk.NewCoins(sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000)),
				ReservePrice: sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000),
				Duration:     time.Duration(30) * time.Second,
				AuctionType:  auctiontypes.RESERVE,
			},
			setupTest: func(tf *auctiontestutil.TestFixture) struct {
				res *auctiontypes.MsgNewAuctionResponse
			} {
				defaultDep := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)

				contract := &keeper.EscrowModContract{
					Id:      contractId,
					Address: f.Addrs[2],
				}

				tf.MockEscrowService.EXPECT().NewContract(tf.Ctx, contractId).Return(contract, nil).AnyTimes()
				tf.MockBankKeeper.EXPECT().SendCoinsFromAccountToModule(tf.Ctx, tf.Addrs[0], auctiontypes.ModuleName, sdk.NewCoins(defaultDep))

				msg := auctiontypes.MsgNewAuction{
					Owner:        f.Addrs[0].String(),
					Deposit:      sdk.NewCoins(sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000)),
					ReservePrice: sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000),
					Duration:     time.Duration(30) * time.Second,
					AuctionType:  auctiontypes.RESERVE,
				}
				msgRes, err := f.MsgServer.NewAuction(f.Ctx, &msg)
				require.NoError(err)
				return struct {
					res *auctiontypes.MsgNewAuctionResponse
				}{
					msgRes,
				}
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mr := tc.setupTest(f)
			queryRes, err := f.QueryServer.Auction(f.Ctx, &auctiontypes.QueryAuctionRequest{
				Id: mr.res.GetId(),
			})

			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
				auction, err := f.K.Auctions.Get(f.Ctx, mr.res.GetId())
				require.NoError(err)

				require.NotNil(auction)
				require.EqualValues(queryRes.Auction.AuctionType, auction.AuctionType)
				require.EqualValues(queryRes.Auction.Strategy.EscrowContractId, auction.Strategy.GetEscrowContractId())
				require.EqualValues(queryRes.Auction.Id, auction.Id)
				require.EqualValues(queryRes.Auction.Owner, auction.Owner)
				require.EqualValues(queryRes.Auction.ReservePrice, auction.ReservePrice)
				require.EqualValues(queryRes.Auction.StartTime, auction.StartTime)
				require.EqualValues(queryRes.Auction.EndTime, auction.EndTime)
			}
		})
	}

}

func TestQueryOwnerAuctions(t *testing.T) {
	f := auctiontestutil.InitFixture(t)
	require := require.New(t)

	testCases := []struct {
		name      string
		owner     sdk.AccAddress
		req       auctiontypes.MsgNewAuction
		expErr    bool
		setupTest func(fixture *auctiontestutil.TestFixture) struct {
			ownerAuctions []uint64
		}
	}{
		{
			name:  "retrieve owner's auctions",
			owner: f.Addrs[0],
			req: auctiontypes.MsgNewAuction{
				Owner:        f.Addrs[0].String(),
				Deposit:      sdk.NewCoins(sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000)),
				ReservePrice: sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000),
				Duration:     time.Duration(30) * time.Second,
				AuctionType:  auctiontypes.RESERVE,
			},
			setupTest: func(tf *auctiontestutil.TestFixture) struct {
				ownerAuctions []uint64
			} {
				defaultDep := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)

				contract1 := &keeper.EscrowModContract{
					Id:      uint64(0),
					Address: f.Addrs[2],
				}

				contract2 := &keeper.EscrowModContract{
					Id:      uint64(1),
					Address: f.Addrs[2],
				}

				tf.MockEscrowService.EXPECT().NewContract(tf.Ctx, contract1.Id).Return(contract1, nil).Times(1)
				tf.MockEscrowService.EXPECT().NewContract(tf.Ctx, contract2.Id).Return(contract2, nil).Times(1)
				tf.MockBankKeeper.EXPECT().SendCoinsFromAccountToModule(tf.Ctx, tf.Addrs[0], auctiontypes.ModuleName, sdk.NewCoins(defaultDep)).Times(2)

				msg1 := auctiontypes.MsgNewAuction{
					Owner:        f.Addrs[0].String(),
					Deposit:      sdk.NewCoins(sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000)),
					ReservePrice: sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000),
					Duration:     time.Duration(30) * time.Second,
					AuctionType:  auctiontypes.RESERVE,
				}
				msg2 := auctiontypes.MsgNewAuction{
					Owner:        f.Addrs[0].String(),
					Deposit:      sdk.NewCoins(sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000)),
					ReservePrice: sdk.NewInt64Coin(f.K.GetDefaultDenom(), 5000),
					Duration:     time.Duration(20) * time.Second,
					AuctionType:  auctiontypes.RESERVE,
				}
				msgRes1, err := f.MsgServer.NewAuction(f.Ctx, &msg1)
				require.NoError(err)
				msgRes2, err := f.MsgServer.NewAuction(f.Ctx, &msg2)
				require.NoError(err)
				return struct {
					ownerAuctions []uint64
				}{
					[]uint64{msgRes1.GetId(), msgRes2.GetId()},
				}
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			msgRes := tc.setupTest(f)
			queryRes, err := f.QueryServer.OwnerAuctions(f.Ctx, &auctiontypes.QueryOwnerAuctionsRequest{
				OwnerAddress: tc.owner.String(),
			})

			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
				var expectedAuctions []auctiontypes.ReserveAuction
				for _, aId := range msgRes.ownerAuctions {
					a, err := f.K.Auctions.Get(f.Ctx, aId)
					require.NoError(err)
					expectedAuctions = append(expectedAuctions, a)
				}

				require.NotNil(queryRes.Auctions)
				require.Equal(len(queryRes.Auctions), len(expectedAuctions))
			}
		})
	}

}

func TestQueryGetAllAuctions(t *testing.T) {
	f := auctiontestutil.InitFixture(t)
	require := require.New(t)

	f.MockBankKeeper.EXPECT().
		SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Eq(auctiontypes.ModuleName), gomock.Any()).
		Return(nil).
		AnyTimes()

	auctions := []auctiontypes.ReserveAuction{
		{Id: 1, Owner: f.Addrs[0].String()},
		{Id: 2, Owner: f.Addrs[1].String()},
	}
	for _, auction := range auctions {
		err := f.K.Auctions.Set(f.Ctx, auction.Id, auction)
		require.NoError(err)
	}

	queryRes, err := f.QueryServer.AllAuctions(f.Ctx, &auctiontypes.QueryAllAuctionsRequest{})
	require.NoError(err)
	require.Len(queryRes.Auctions, len(auctions))
}
