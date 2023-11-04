package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	auctiontypes "github.com/fatal-fruit/auction/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewAuction(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)

	testCases := []struct {
		name      string
		req       auctiontypes.MsgNewAuction
		expErr    bool
		setupTest func(fixture *testFixture) struct {
			contractId uint64
		}
	}{
		{
			name: "valid auction",
			req: auctiontypes.MsgNewAuction{
				Owner:        f.addrs[0].String(),
				Deposit:      sdk.NewCoins(sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000)),
				ReservePrice: sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000),
				Duration:     time.Duration(30) * time.Second,
				AuctionType:  auctiontypes.RESERVE,
			},
			setupTest: func(tf *testFixture) struct {
				contractId uint64
			} {
				contractId := uint64(1)
				defaultModBalance := sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000)
				defaultDep := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)

				tf.mockEscrowService.EXPECT().NewContract().Return(contractId, nil).AnyTimes()
				tf.mockAcctKeeper.EXPECT().GetAccount(tf.ctx, tf.modAddr).Return(tf.modAccount).AnyTimes()
				tf.mockAcctKeeper.EXPECT().GetModuleAddress(auctiontypes.ModuleName).Return(tf.modAddr).AnyTimes()
				tf.mockBankKeeper.EXPECT().GetBalance(tf.ctx, tf.modAddr, tf.k.GetDefaultDenom()).Return(defaultModBalance)
				tf.mockBankKeeper.EXPECT().SendCoinsFromAccountToModule(tf.ctx, tf.addrs[0], auctiontypes.ModuleName, sdk.NewCoins(defaultDep))
				tf.mockBankKeeper.EXPECT().GetBalance(tf.ctx, tf.modAddr, tf.k.GetDefaultDenom()).Return(defaultModBalance.Add(defaultDep))
				return struct {
					contractId uint64
				}{
					contractId: contractId,
				}
			},
		},
		//{
		//	name: "insufficient owner balance",
		//	req: auctiontypes.MsgNewAuction{
		//		Owner:        f.addrs[0].String(),
		//		Deposit:      sdk.NewCoins(sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000)),
		//		ReservePrice: sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000),
		//		Duration:     time.Duration(30) * time.Second,
		//		AuctionType:  auctiontypes.RESERVE,
		//	},
		//},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			expValues := tc.setupTest(f)
			modAcctBalance := f.k.GetModuleBalance(f.ctx, f.k.GetDefaultDenom())
			found, deposit := tc.req.Deposit.Find(f.k.GetDefaultDenom())
			require.True(found)
			res, err := f.msgServer.NewAuction(f.ctx, &tc.req)
			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
				auction, err := f.k.Auctions.Get(f.ctx, res.GetId())
				newModAcctBalance := f.k.GetModuleBalance(f.ctx, f.k.GetDefaultDenom())
				require.NoError(err)
				require.NotNil(auction)
				require.Equal(tc.req.Duration, auction.EndTime.Sub(auction.StartTime))
				require.Equal(tc.req.ReservePrice, auction.ReservePrice)
				require.Equal(expValues.contractId, auction.EscrowContract)
				require.Equal(tc.req.Owner, auction.Owner)
				require.Equal(newModAcctBalance.Sub(deposit), modAcctBalance)
			}
		})
	}

}
