package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/fatal-fruit/auction/keeper"
	auctiontypes "github.com/fatal-fruit/auction/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQueryAuction(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)

	testCases := []struct {
		name      string
		req       auctiontypes.MsgNewAuction
		expErr    bool
		setupTest func(fixture *testFixture) struct {
			res *auctiontypes.MsgNewAuctionResponse
		}
	}{
		{
			name: "retrieve valid auction",
			req: auctiontypes.MsgNewAuction{
				Owner:        f.addrs[0].String(),
				Deposit:      sdk.NewCoins(sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000)),
				ReservePrice: sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000),
				Duration:     time.Duration(30) * time.Second,
				AuctionType:  auctiontypes.RESERVE,
			},
			setupTest: func(tf *testFixture) struct {
				res *auctiontypes.MsgNewAuctionResponse
			} {
				contractId := uint64(1)
				defaultDep := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)

				contract := &keeper.EscrowModContract{
					Id:      contractId,
					Address: f.addrs[2],
				}

				tf.mockEscrowService.EXPECT().NewContract(tf.ctx, contractId).Return(contract, nil).AnyTimes()
				tf.mockBankKeeper.EXPECT().SendCoinsFromAccountToModule(tf.ctx, tf.addrs[0], auctiontypes.ModuleName, sdk.NewCoins(defaultDep))

				msg := auctiontypes.MsgNewAuction{
					Owner:        f.addrs[0].String(),
					Deposit:      sdk.NewCoins(sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000)),
					ReservePrice: sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000),
					Duration:     time.Duration(30) * time.Second,
					AuctionType:  auctiontypes.RESERVE,
				}
				msgRes, err := f.msgServer.NewAuction(f.ctx, &msg)
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
			queryRes, err := f.queryServer.Auction(f.ctx, &auctiontypes.QueryAuctionRequest{
				Id: mr.res.GetId(),
			})

			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
				auction, err := f.k.Auctions.Get(f.ctx, mr.res.GetId())
				require.NoError(err)

				require.NotNil(auction)
				require.EqualValues(queryRes.Auction.AuctionType, auction.AuctionType)
				require.EqualValues(queryRes.Auction.EscrowContract, auction.EscrowContract)
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
	f := initFixture(t)
	require := require.New(t)

	testCases := []struct {
		name      string
		owner     sdk.AccAddress
		req       auctiontypes.MsgNewAuction
		expErr    bool
		setupTest func(fixture *testFixture) struct {
			ownerAuctions []uint64
		}
	}{
		{
			name:  "retrieve owner's auctions",
			owner: f.addrs[0],
			req: auctiontypes.MsgNewAuction{
				Owner:        f.addrs[0].String(),
				Deposit:      sdk.NewCoins(sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000)),
				ReservePrice: sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000),
				Duration:     time.Duration(30) * time.Second,
				AuctionType:  auctiontypes.RESERVE,
			},
			setupTest: func(tf *testFixture) struct {
				ownerAuctions []uint64
			} {
				contractId := uint64(1)
				defaultDep := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)

				contract := &keeper.EscrowModContract{
					Id:      contractId,
					Address: f.addrs[2],
				}

				tf.mockEscrowService.EXPECT().NewContract(tf.ctx, contractId).Return(contract, nil).AnyTimes()
				tf.mockBankKeeper.EXPECT().SendCoinsFromAccountToModule(tf.ctx, tf.addrs[0], auctiontypes.ModuleName, sdk.NewCoins(defaultDep)).Times(2)

				msg1 := auctiontypes.MsgNewAuction{
					Owner:        f.addrs[0].String(),
					Deposit:      sdk.NewCoins(sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000)),
					ReservePrice: sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000),
					Duration:     time.Duration(30) * time.Second,
					AuctionType:  auctiontypes.RESERVE,
				}
				msg2 := auctiontypes.MsgNewAuction{
					Owner:        f.addrs[0].String(),
					Deposit:      sdk.NewCoins(sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000)),
					ReservePrice: sdk.NewInt64Coin(f.k.GetDefaultDenom(), 5000),
					Duration:     time.Duration(20) * time.Second,
					AuctionType:  auctiontypes.RESERVE,
				}
				msgRes1, err := f.msgServer.NewAuction(f.ctx, &msg1)
				require.NoError(err)
				msgRes2, err := f.msgServer.NewAuction(f.ctx, &msg2)
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
			queryRes, err := f.queryServer.OwnerAuctions(f.ctx, &auctiontypes.QueryOwnerAuctionsRequest{
				OwnerAddress: tc.owner.String(),
			})

			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
				var expectedAuctions []auctiontypes.ReserveAuction
				for _, aId := range msgRes.ownerAuctions {
					a, err := f.k.Auctions.Get(f.ctx, aId)
					require.NoError(err)
					expectedAuctions = append(expectedAuctions, a)
				}

				require.NotNil(queryRes.Auctions)
				require.Equal(len(queryRes.Auctions), len(expectedAuctions))
			}
		})
	}

}
