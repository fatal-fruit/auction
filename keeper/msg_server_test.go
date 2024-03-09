package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/fatal-fruit/auction/keeper"
	auctiontestutil "github.com/fatal-fruit/auction/testutil"
	auctiontypes "github.com/fatal-fruit/auction/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewAuction(t *testing.T) {
	f := auctiontestutil.InitFixture(t)
	require := require.New(t)

	testCases := []struct {
		name      string
		req       auctiontypes.MsgNewAuction
		expErr    bool
		setupTest func(fixture *auctiontestutil.TestFixture) struct {
			contractId uint64
		}
	}{
		{
			name: "valid auction",
			req: auctiontypes.MsgNewAuction{
				Owner:        f.Addrs[0].String(),
				Deposit:      sdk.NewCoins(sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000)),
				ReservePrice: sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000),
				Duration:     time.Duration(30) * time.Second,
				AuctionType:  auctiontypes.RESERVE,
			},
			setupTest: func(tf *auctiontestutil.TestFixture) struct {
				contractId uint64
			} {
				contractId := uint64(1)
				defaultModBalance := sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000)
				defaultDep := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)
				contract := &keeper.EscrowModContract{
					Id:      contractId,
					Address: f.Addrs[2],
				}

				tf.MockEscrowService.EXPECT().NewContract(tf.Ctx, contractId).Return(contract, nil).AnyTimes()
				tf.MockAcctKeeper.EXPECT().GetAccount(tf.Ctx, tf.ModAddr).Return(tf.ModAccount).AnyTimes()
				tf.MockAcctKeeper.EXPECT().GetModuleAddress(auctiontypes.ModuleName).Return(tf.ModAddr).AnyTimes()
				tf.MockBankKeeper.EXPECT().GetBalance(tf.Ctx, tf.ModAddr, tf.K.GetDefaultDenom()).Return(defaultModBalance)
				tf.MockBankKeeper.EXPECT().SendCoinsFromAccountToModule(tf.Ctx, tf.Addrs[0], auctiontypes.ModuleName, sdk.NewCoins(defaultDep))
				tf.MockBankKeeper.EXPECT().GetBalance(tf.Ctx, tf.ModAddr, tf.K.GetDefaultDenom()).Return(defaultModBalance.Add(defaultDep))
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
			modAcctBalance := f.K.GetModuleBalance(f.Ctx, f.K.GetDefaultDenom())
			found, deposit := tc.req.Deposit.Find(f.K.GetDefaultDenom())
			require.True(found)
			res, err := f.MsgServer.NewAuction(f.Ctx, &tc.req)
			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
				auction, err := f.K.Auctions.Get(f.Ctx, res.GetId())
				require.NoError(err)
				isActive, err := f.K.ActiveAuctions.Has(f.Ctx, res.GetId())
				require.NoError(err)
				require.True(isActive)

				newModAcctBalance := f.K.GetModuleBalance(f.Ctx, f.K.GetDefaultDenom())
				require.NotNil(auction)
				require.Equal(tc.req.Duration, auction.Duration)
				require.Equal(tc.req.ReservePrice, auction.ReservePrice)
				require.Equal(expValues.contractId, auction.Strategy.EscrowContractId)
				require.Equal(tc.req.Owner, auction.Owner)
				require.Equal(newModAcctBalance.Sub(deposit), modAcctBalance)
				require.Equal(len(auction.Bids), 0)
			}
		})
	}
}

func TestNewBid(t *testing.T) {
	f := auctiontestutil.InitFixture(t)
	require := require.New(t)

	testCases := []struct {
		name      string
		owner     sdk.AccAddress
		bid       sdk.Coin
		expErr    bool
		setupTest func(fixture *auctiontestutil.TestFixture) struct {
			contractId uint64
		}
	}{
		{
			name:   "valid bid",
			owner:  f.Addrs[1],
			bid:    sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1100),
			expErr: false,
			setupTest: func(tf *auctiontestutil.TestFixture) struct {
				contractId uint64
			} {
				contractId := uint64(1)
				defaultDep := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)

				contract := &keeper.EscrowModContract{
					Id:      contractId,
					Address: f.Addrs[2],
				}
				tf.MockEscrowService.EXPECT().NewContract(tf.Ctx, contractId).Return(contract, nil).AnyTimes()
				tf.MockBankKeeper.EXPECT().SendCoinsFromAccountToModule(tf.Ctx, tf.Addrs[0], auctiontypes.ModuleName, sdk.NewCoins(defaultDep)).Times(1)

				msg1 := auctiontypes.MsgNewAuction{
					Owner:        f.Addrs[0].String(),
					Deposit:      sdk.NewCoins(sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000)),
					ReservePrice: sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000),
					Duration:     time.Duration(30) * time.Second,
					AuctionType:  auctiontypes.RESERVE,
				}

				auctionRes, err := f.MsgServer.NewAuction(f.Ctx, &msg1)
				require.NoError(err)

				_, err = f.MsgServer.StartAuction(f.Ctx, &auctiontypes.MsgStartAuction{
					Owner: f.Addrs[0].String(),
					Id:    auctionRes.GetId(),
				})
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
			owner:  f.Addrs[1],
			bid:    sdk.NewInt64Coin(f.K.GetDefaultDenom(), 900),
			expErr: true,
			setupTest: func(tf *auctiontestutil.TestFixture) struct {
				contractId uint64
			} {
				contractId := uint64(1)
				defaultDep := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)

				contract := &keeper.EscrowModContract{
					Id:      contractId,
					Address: f.Addrs[2],
				}
				tf.MockEscrowService.EXPECT().NewContract(tf.Ctx, contractId).Return(contract, nil).AnyTimes()
				tf.MockBankKeeper.EXPECT().SendCoinsFromAccountToModule(tf.Ctx, tf.Addrs[0], auctiontypes.ModuleName, sdk.NewCoins(defaultDep)).Times(1)

				msg1 := auctiontypes.MsgNewAuction{
					Owner:        f.Addrs[0].String(),
					Deposit:      sdk.NewCoins(sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000)),
					ReservePrice: sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000),
					Duration:     time.Duration(30) * time.Second,
					AuctionType:  auctiontypes.RESERVE,
				}

				auctionRes, err := f.MsgServer.NewAuction(f.Ctx, &msg1)
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
			owner:  f.Addrs[1],
			bid:    sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1100),
			expErr: true,
			setupTest: func(tf *auctiontestutil.TestFixture) struct {
				contractId uint64
			} {
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

				return struct {
					contractId uint64
				}{
					id,
				}
			},
		},
		{
			name:   "bid lower than competitive bid",
			owner:  f.Addrs[2],
			bid:    sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000),
			expErr: true,
			setupTest: func(tf *auctiontestutil.TestFixture) struct {
				contractId uint64
			} {
				id, err := f.K.IDs.Next(f.Ctx)
				require.NoError(err)
				auction := auctiontypes.ReserveAuction{
					Id:           id,
					Status:       auctiontypes.ACTIVE,
					Owner:        f.Addrs[0].String(),
					AuctionType:  auctiontypes.RESERVE,
					ReservePrice: sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1000),
					StartTime:    time.Now(),
					EndTime:      time.Now().Add(30 * time.Second),
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
			_, err := f.MsgServer.NewBid(f.Ctx, &bid)

			if tc.expErr {
				require.Error(err)
			} else {
				auction, err := f.K.Auctions.Get(f.Ctx, msgRes.contractId)
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

func TestExecAuction(t *testing.T) {
	f := auctiontestutil.InitFixture(t)
	require := require.New(t)

	testCases := []struct {
		name      string
		req       auctiontypes.MsgExecAuction
		expErr    bool
		setupTest func(fixture *auctiontestutil.TestFixture) struct {
			auctionId uint64
		}
	}{
		{
			name: "execute pending auction",
			req: auctiontypes.MsgExecAuction{
				Sender: f.Addrs[0].String(),
			},
			setupTest: func(tf *auctiontestutil.TestFixture) struct {
				auctionId uint64
			} {
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
						EscrowContractId:      uint64(1),
						EscrowContractAddress: f.Addrs[2].String(),
					},
				}
				err = f.K.Auctions.Set(f.Ctx, id, auction)

				require.NoError(err)
				err = f.K.PendingAuctions.Set(f.Ctx, id)
				require.NoError(err)
				f.MockBankKeeper.EXPECT().SendCoins(f.Ctx, f.Addrs[1], f.Addrs[0], sdk.Coins{sdk.NewInt64Coin(f.K.GetDefaultDenom(), 1100)}).Times(1)
				f.MockEscrowService.EXPECT().Release(f.Ctx, uint64(1), f.Addrs[2], f.Addrs[1]).Times(1)

				return struct {
					auctionId uint64
				}{
					auctionId: id,
				}
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			msgRes := tc.setupTest(f)
			tc.req.AuctionId = msgRes.auctionId

			_, err := f.MsgServer.Exec(f.Ctx, &tc.req)

			if tc.expErr {
				require.Error(err)
			} else {
				// set status to executed
				// remove from pending
				// expect escrow service release to have been called
				isPending, err := f.K.PendingAuctions.Has(f.Ctx, msgRes.auctionId)
				require.NoError(err)
				require.False(isPending)

				//expect bank send coins called two more times

				auction, err := f.K.Auctions.Get(f.Ctx, msgRes.auctionId)
				require.Equal(auction.Status, auctiontypes.CLOSED)
				require.NoError(err)

			}
		})
	}
}
