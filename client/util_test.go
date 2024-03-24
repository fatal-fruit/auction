package client

import (
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	at "github.com/fatal-fruit/auction/auctiontypes"
	auctiontypes "github.com/fatal-fruit/auction/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestParseAuctionType(t *testing.T) {
	encConfig := moduletestutil.MakeTestEncodingConfig()
	encConfig.InterfaceRegistry.RegisterInterface(
		"fatal_fruit.auction.v1.AuctionMetadata",
		(*auctiontypes.AuctionMetadata)(nil),
		&at.ReserveAuctionMetadata{},
	)
	encConfig.InterfaceRegistry.RegisterInterface(
		"fatal_fruit.auction.v1.Auction",
		(*auctiontypes.Auction)(nil),
		&at.ReserveAuction{},
	)
	cdc := encConfig.Codec
	auctionType := "/fatal_fruit.auction.v1.ReserveAuction"

	res, err := parseAuctionType(cdc, auctionType)
	require.NoError(t, err)
	require.Equal(t, res, auctionType)
}

func TestParseMetadata(t *testing.T) {
	encConfig := moduletestutil.MakeTestEncodingConfig()
	encConfig.InterfaceRegistry.RegisterInterface(
		"fatal_fruit.auction.v1.AuctionMetadata",
		(*auctiontypes.AuctionMetadata)(nil),
		&at.ReserveAuctionMetadata{},
	)
	encConfig.InterfaceRegistry.RegisterInterface(
		"fatal_fruit.auction.v1.Auction",
		(*auctiontypes.Auction)(nil),
		&at.ReserveAuction{},
	)
	cdc := encConfig.Codec
	md := `{
				"@type": "/fatal_fruit.auction.v1.ReserveAuctionMetadata",
				"duration": "1000000ms",  
				"reserve_price": {
					"denom":"stake",
					"amount":"250"
				}
			}`

	mdFile := testutil.WriteToNewTempFile(t, md)
	defer mdFile.Close()

	expMsg := &at.ReserveAuctionMetadata{
		Duration:     time.Duration(1000000000000),
		ReservePrice: sdk.NewInt64Coin("stake", int64(250)),
	}

	res, err := parseAuctionMetadata(cdc, mdFile.Name())
	require.NoError(t, err)
	require.Equal(t, res.String(), expMsg.String())
}
