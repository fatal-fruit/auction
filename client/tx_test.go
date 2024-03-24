package client

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	rpcclientmock "github.com/cometbft/cometbft/rpc/client/mock"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	at "github.com/fatal-fruit/auction/auctiontypes"
	auctiontypes "github.com/fatal-fruit/auction/types"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestCreateAuction(t *testing.T) {
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
	kr := keyring.NewInMemory(encConfig.Codec)
	baseContext := client.Context{}.
		WithKeyring(kr).
		WithTxConfig(encConfig.TxConfig).
		WithCodec(encConfig.Codec).
		WithClient(clitestutil.MockCometRPC{Client: rpcclientmock.Client{}}).
		WithAccountRetriever(client.MockAccountRetriever{}).
		WithOutput(io.Discard).WithChainID("test-chain")

	accounts := testutil.CreateKeyringAccounts(t, kr, 1)
	val := accounts[0]

	ctxGen := func() client.Context {
		bz, _ := encConfig.Codec.Marshal(&sdk.TxResponse{})
		c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
			Value: bz,
		})
		return baseContext.WithClient(c)
	}
	clientCtx := ctxGen()

	auctionType := "/fatal_fruit.auction.v1.ReserveAuctionMetadata"
	deposit := "250stake"
	metadata := `{
				"@type": "/fatal_fruit.auction.v1.ReserveAuctionMetadata",
				"duration": "1000000ms",  
				"reserve_price": {
					"denom":"stake",
					"amount":"250"
				}
			}`

	mdFile := testutil.WriteToNewTempFile(t, metadata)
	defer mdFile.Close()

	args := []string{
		auctionType,
		mdFile.Name(),
		deposit,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(10))).String()),
	}

	cmd := NewAuctionCmd()
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
	require.NoError(t, err)
	require.NotNil(t, out)

	msg := &sdk.TxResponse{}
	require.NoError(t, clientCtx.Codec.UnmarshalJSON(out.Bytes(), msg))
}
