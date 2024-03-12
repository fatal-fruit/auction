package keeper_test

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	auctiontestutil "github.com/fatal-fruit/auction/testutil"
	auctiontypes "github.com/fatal-fruit/auction/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestQueryAuction(t *testing.T) {
	f := auctiontestutil.InitFixture(t)
	require := require.New(t)
	contractId := uint64(0)

	metadata := auctiontypes.ReserveAuctionMetadata{
		ReservePrice: sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000),
		Duration:     time.Duration(30) * time.Second,
	}
	anyMd, err := codectypes.NewAnyWithValue(&metadata)
	require.NoError(err)

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
				Owner:           f.Addrs[0].String(),
				Deposit:         sdk.NewCoins(sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000)),
				AuctionType:     f.ReserveAuctionType,
				AuctionMetadata: anyMd,
			},
			setupTest: func(tf *auctiontestutil.TestFixture) struct {
				res *auctiontypes.MsgNewAuctionResponse
			} {
				defaultDep := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)

				contract := &auctiontestutil.EscrowModContract{
					Id:      contractId,
					Address: f.Addrs[2],
				}

				tf.MockEscrowService.EXPECT().NewContract(tf.Ctx, contractId).Return(contract, nil).AnyTimes()
				tf.MockBankKeeper.EXPECT().SendCoinsFromAccountToModule(tf.Ctx, tf.Addrs[0], auctiontypes.ModuleName, sdk.NewCoins(defaultDep))

				msg := auctiontypes.MsgNewAuction{
					Owner:           f.Addrs[0].String(),
					Deposit:         sdk.NewCoins(sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000)),
					AuctionType:     f.ReserveAuctionType,
					AuctionMetadata: anyMd,
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
				require.NotNil(queryRes)

				auction, err := f.K.Auctions.Get(f.Ctx, mr.res.GetId())
				require.NoError(err)

				var act auctiontypes.Auction
				res := queryRes.GetAuction()
				err = f.EnCfg.InterfaceRegistry.UnpackAny(res, &act)
				require.NoError(err)

				switch r := act.(type) {
				case *auctiontypes.ReserveAuction:
					require.EqualValues(r.Id, auction.GetId())
					require.EqualValues(r.GetType(), auction.GetType())
				default:
					t.Errorf("invalid auction type")
				}
			}
		})
	}

}

func TestQueryOwnerAuctions(t *testing.T) {
	f := auctiontestutil.InitFixture(t)
	require := require.New(t)

	metadata := auctiontypes.ReserveAuctionMetadata{
		ReservePrice: sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000),
		Duration:     time.Duration(30) * time.Second,
	}
	anyMd, err := codectypes.NewAnyWithValue(&metadata)
	require.NoError(err)

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
				Owner:           f.Addrs[0].String(),
				Deposit:         sdk.NewCoins(sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000)),
				AuctionType:     f.ReserveAuctionType,
				AuctionMetadata: anyMd,
			},
			setupTest: func(tf *auctiontestutil.TestFixture) struct {
				ownerAuctions []uint64
			} {
				defaultDep := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)

				contract1 := &auctiontestutil.EscrowModContract{
					Id:      uint64(0),
					Address: f.Addrs[2],
				}

				contract2 := &auctiontestutil.EscrowModContract{
					Id:      uint64(1),
					Address: f.Addrs[2],
				}

				tf.MockEscrowService.EXPECT().NewContract(tf.Ctx, contract1.Id).Return(contract1, nil).Times(1)
				tf.MockEscrowService.EXPECT().NewContract(tf.Ctx, contract2.Id).Return(contract2, nil).Times(1)
				tf.MockBankKeeper.EXPECT().SendCoinsFromAccountToModule(tf.Ctx, tf.Addrs[0], auctiontypes.ModuleName, sdk.NewCoins(defaultDep)).Times(2)

				msg1 := auctiontypes.MsgNewAuction{
					Owner:           f.Addrs[0].String(),
					Deposit:         sdk.NewCoins(sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000)),
					AuctionType:     f.ReserveAuctionType,
					AuctionMetadata: anyMd,
				}

				metadata2 := auctiontypes.ReserveAuctionMetadata{
					ReservePrice: sdk.NewInt64Coin(f.K.GetDefaultDenom(), 5000),
					Duration:     time.Duration(20) * time.Second,
				}
				anyMd2, err := codectypes.NewAnyWithValue(&metadata2)
				require.NoError(err)

				msg2 := auctiontypes.MsgNewAuction{
					Owner:           f.Addrs[0].String(),
					Deposit:         sdk.NewCoins(sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000)),
					AuctionType:     f.ReserveAuctionType,
					AuctionMetadata: anyMd2,
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
				var expectedAuctions []auctiontypes.Auction
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
	// TODO: Fix
	t.Skip()
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
		err := f.K.Auctions.Set(f.Ctx, auction.Id, &auction)
		require.NoError(err)
	}

	queryRes, err := f.QueryServer.AllAuctions(f.Ctx, &auctiontypes.QueryAllAuctionsRequest{})
	require.NoError(err)
	require.Len(queryRes.Auctions, len(auctions))
}
