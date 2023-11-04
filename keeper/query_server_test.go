package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
				//defaultModBalance := sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000)
				defaultDep := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)

				tf.mockEscrowService.EXPECT().NewContract().Return(contractId, nil).AnyTimes()
				//tf.mockAcctKeeper.EXPECT().GetAccount(tf.ctx, tf.modAddr).Return(tf.modAccount).AnyTimes()
				//tf.mockAcctKeeper.EXPECT().GetModuleAddress(auctiontypes.ModuleName).Return(tf.modAddr).AnyTimes()
				//tf.mockBankKeeper.EXPECT().GetBalance(tf.ctx, tf.modAddr, tf.k.GetDefaultDenom()).Return(defaultModBalance)
				tf.mockBankKeeper.EXPECT().SendCoinsFromAccountToModule(tf.ctx, tf.addrs[0], auctiontypes.ModuleName, sdk.NewCoins(defaultDep))
				//tf.mockBankKeeper.EXPECT().GetBalance(tf.ctx, tf.modAddr, tf.k.GetDefaultDenom()).Return(defaultModBalance.Add(defaultDep))

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
