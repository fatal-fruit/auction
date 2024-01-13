package client

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	auctiontypes "github.com/fatal-fruit/auction/types"
	"github.com/spf13/cobra"
)

// GetCmdQueryAllAuctions returns the command for querying all auctions
func GetCmdQueryAllAuctions(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "all-auctions",
		Short: "List all auctions",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := auctiontypes.NewQueryClient(clientCtx)
			res, err := queryClient.AllAuctions(context.Background(), &auctiontypes.QueryAllAuctionsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

func CmdListOwnerAuctions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner-auctions [owner-address]",
		Short: "Query auctions by owner address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := auctiontypes.NewQueryClient(clientCtx)
			res, err := queryClient.OwnerAuctions(context.Background(), &auctiontypes.QueryOwnerAuctionsRequest{
				OwnerAddress: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

