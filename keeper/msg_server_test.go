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
		name   string
		req    auctiontypes.MsgNewAuction
		expErr bool
	}{
		{
			"valid auction",
			auctiontypes.MsgNewAuction{
				Owner:       f.addrs[0].String(),
				Deposit:     sdk.NewCoins(sdk.NewInt64Coin(f.k.GetDefaultDenom(), 1000)),
				Duration:    time.Duration(30) * time.Second,
				AuctionType: auctiontypes.RESERVE,
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res, err := f.msgServer.NewAuction(f.ctx, &tc.req)
			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
				auction, err := f.k.Auctions.Get(f.ctx, res.GetId())
				require.NoError(err)
				require.NotNil(auction)
			}
		})
	}

}
