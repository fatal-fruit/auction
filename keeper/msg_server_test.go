package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/fatal-fruit/auction/keeper"
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
				contract := &keeper.EscrowModContract{
					Id:      contractId,
					Address: f.addrs[2],
				}

				tf.mockEscrowService.EXPECT().NewContract(tf.ctx, contractId).Return(contract, nil).AnyTimes()
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
				require.NoError(err)
				isActive, err := f.k.ActiveAuctions.Has(f.ctx, res.GetId())
				require.NoError(err)
				require.True(isActive)

				newModAcctBalance := f.k.GetModuleBalance(f.ctx, f.k.GetDefaultDenom())
				require.NotNil(auction)
				require.Equal(tc.req.Duration, auction.EndTime.Sub(auction.StartTime))
				require.Equal(tc.req.ReservePrice, auction.ReservePrice)
				require.Equal(expValues.contractId, auction.EscrowContract)
				require.Equal(tc.req.Owner, auction.Owner)
				require.Equal(newModAcctBalance.Sub(deposit), modAcctBalance)
				require.Equal(len(auction.Bids), 0)
			}
		})
	}
}

func TestNewBid(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)

	testCases := []struct {
		name      string
		owner     sdk.AccAddress
		bid       sdk.Coin
		expErr    bool
		setupTest func(fixture *testFixture) struct {
			contractId uint64
		}
	}{
		{
			name:   "valid bid",
			owner:  f.addrs[1],
			bid:    sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1100),
			expErr: false,
			setupTest: func(tf *testFixture) struct {
				contractId uint64
			} {
				contractId := uint64(1)
				defaultDep := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)

				contract := &keeper.EscrowModContract{
					Id:      contractId,
					Address: f.addrs[2],
				}
				tf.mockEscrowService.EXPECT().NewContract(tf.ctx, contractId).Return(contract, nil).AnyTimes()
				tf.mockBankKeeper.EXPECT().SendCoinsFromAccountToModule(tf.ctx, tf.addrs[0], auctiontypes.ModuleName, sdk.NewCoins(defaultDep)).Times(1)

				msg1 := auctiontypes.MsgNewAuction{
					Owner:        f.addrs[0].String(),
					Deposit:      sdk.NewCoins(sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000)),
					ReservePrice: sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000),
					Duration:     time.Duration(30) * time.Second,
					AuctionType:  auctiontypes.RESERVE,
				}

				auctionRes, err := f.msgServer.NewAuction(f.ctx, &msg1)
				require.NoError(err)

				return struct {
					contractId uint64
				}{
					auctionRes.GetId(),
				}
			},
		},
		{
			name:   "invalid bid price",
			owner:  f.addrs[1],
			bid:    sdk.NewInt64Coin(f.k.GetDefaultDenom(), 900),
			expErr: true,
			setupTest: func(tf *testFixture) struct {
				contractId uint64
			} {
				contractId := uint64(1)
				defaultDep := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)

				contract := &keeper.EscrowModContract{
					Id:      contractId,
					Address: f.addrs[2],
				}
				tf.mockEscrowService.EXPECT().NewContract(tf.ctx, contractId).Return(contract, nil).AnyTimes()
				tf.mockBankKeeper.EXPECT().SendCoinsFromAccountToModule(tf.ctx, tf.addrs[0], auctiontypes.ModuleName, sdk.NewCoins(defaultDep)).Times(1)

				msg1 := auctiontypes.MsgNewAuction{
					Owner:        f.addrs[0].String(),
					Deposit:      sdk.NewCoins(sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000)),
					ReservePrice: sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000),
					Duration:     time.Duration(30) * time.Second,
					AuctionType:  auctiontypes.RESERVE,
				}

				auctionRes, err := f.msgServer.NewAuction(f.ctx, &msg1)
				require.NoError(err)

				return struct {
					contractId uint64
				}{
					auctionRes.GetId(),
				}
			},
		},
		{
			name:   "bid for expired auction",
			owner:  f.addrs[1],
			bid:    sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1100),
			expErr: true,
			setupTest: func(tf *testFixture) struct {
				contractId uint64
			} {
				id, err := f.k.IDs.Next(f.ctx)
				require.NoError(err)
				auction := auctiontypes.ReserveAuction{
					Id:             id,
					Owner:          f.addrs[0].String(),
					AuctionType:    auctiontypes.RESERVE,
					EscrowContract: 1,
					ReservePrice:   sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000),
					StartTime:      time.Now().Add(-30 * time.Second),
					EndTime:        time.Now().Add(-1 * time.Second),
					Bids:           []*auctiontypes.Bid{},
					Strategy: &auctiontypes.SettleStrategy{
						StrategyType:          auctiontypes.SETTLE,
						EscrowContractId:      1,
						EscrowContractAddress: f.addrs[2].String(),
					},
				}
				err = f.k.Auctions.Set(f.ctx, id, auction)
				require.NoError(err)
				err = f.k.ActiveAuctions.Set(f.ctx, id)
				require.NoError(err)

				return struct {
					contractId uint64
				}{
					id,
				}
			},
		},
		{
			name:   "bid lower than competitive bid",
			owner:  f.addrs[2],
			bid:    sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000),
			expErr: true,
			setupTest: func(tf *testFixture) struct {
				contractId uint64
			} {
				id, err := f.k.IDs.Next(f.ctx)
				require.NoError(err)
				auction := auctiontypes.ReserveAuction{
					Id:             id,
					Owner:          f.addrs[0].String(),
					AuctionType:    auctiontypes.RESERVE,
					EscrowContract: 1,
					ReservePrice:   sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000),
					StartTime:      time.Now(),
					EndTime:        time.Now().Add(30 * time.Second),
					LastPrice:      sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1100),
					Bids: []*auctiontypes.Bid{
						{
							AuctionId: id,
							Bidder:    f.addrs[1].String(),
							BidPrice:  sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1100),
							Timestamp: time.Now(),
						},
					},
					Strategy: &auctiontypes.SettleStrategy{
						StrategyType:          auctiontypes.SETTLE,
						EscrowContractId:      1,
						EscrowContractAddress: f.addrs[2].String(),
					},
				}
				err = f.k.Auctions.Set(f.ctx, id, auction)
				require.NoError(err)
				err = f.k.ActiveAuctions.Set(f.ctx, id)
				require.NoError(err)

				return struct {
					contractId uint64
				}{
					id,
				}
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			msgRes := tc.setupTest(f)
			bid := auctiontypes.MsgNewBid{
				AuctionId: msgRes.contractId,
				Owner:     tc.owner.String(),
				Bid:       tc.bid,
			}
			_, err := f.msgServer.NewBid(f.ctx, &bid)

			if tc.expErr {
				require.Error(err)
			} else {
				auction, err := f.k.Auctions.Get(f.ctx, msgRes.contractId)
				require.NoError(err)
				bd := auction.GetBids()[0]
				require.Equal(bd.BidPrice, tc.bid)
				require.Equal(bd.AuctionId, msgRes.contractId)
				require.Equal(bd.Bidder, tc.owner.String())
				require.Equal(auction.LastPrice, bd.BidPrice)
			}
		})
	}
}
