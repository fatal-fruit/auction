package client

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	auctiontypes "github.com/fatal-fruit/auction/types"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

// NewAuctionCmd creates a CLI command for MsgNewAuction.
func NewAuctionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-auction [type] [auction-json-file] [deposit] --from [sender]",
		Args:  cobra.ExactArgs(3),
		Short: "create new auction",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s create-auction <type> auction_metadata.json <deposit> --from <sender> --chain-id <chain-id>
		
		Where auction_type is the type url of the auction
		ex: /fatal_fruit.auction.v1.ReserveAuctionMetadata

		and auction_metadata.json contains:
		{
			"@type": "/fatal_fruit.auction.v1.ReserveAuctionMetadata",
			"duration": "1000000ms",  
				"reserve_price": {
					"denom":"stake",
					"amount":"250"
				}
		}
		`, version.AppName, auctiontypes.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if args[0] == "" || args[1] == "" || args[2] == "" {
				return fmt.Errorf("auction type, metadata, and deposit cannot be empty")
			}

			// Validate auction type
			auctionType, err := parseAuctionType(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}

			// Parse auction metadata
			auctionMetadata, err := parseAuctionMetadata(clientCtx.Codec, args[1])
			if err != nil {
				return err
			}

			// Parse deposit
			deposit, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return err
			}

			owner := clientCtx.GetFromAddress().String()

			msg := auctiontypes.MsgNewAuction{
				Owner:       owner,
				AuctionType: auctionType,
				Deposit:     deposit,
			}

			if err = msg.SetMetadata(auctionMetadata); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// StartAuctionCmd creates a CLI command for MsgStartAuction.
func StartAuctionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-auction [auction-id] --from [sender]",
		Args:  cobra.ExactArgs(3),
		Short: "start auction",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s start-auction <auction-id> --from <sender> --chain-id <chain-id>`,
			version.AppName, auctiontypes.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if args[0] == "" {
				return fmt.Errorf("auction id cannot be empty")
			}

			owner := clientCtx.GetFromAddress().String()
			auctionId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("unable to parse auction id %s", args[0])
			}

			msg := auctiontypes.MsgStartAuction{
				Owner: owner,
				Id:    auctionId,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func BidCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bid [auction-id] [bid-price] --from [sender]",
		Args:  cobra.ExactArgs(2),
		Short: "bid on an auction",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s bid <auction-id> <deposit> --from <sender> --chain-id <chain-id>`, version.AppName, auctiontypes.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if args[0] == "" || args[1] == "" {
				return fmt.Errorf("auction-id and deposit cannot be empty")
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			// Parse Price
			bp, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			// TODO: Remove?
			found, bidPrice := bp.Find(sdk.DefaultBondDenom)
			if !found {
				return fmt.Errorf("invalid bid denom")
			}

			owner := clientCtx.GetFromAddress().String()

			msg := auctiontypes.MsgNewBid{
				BidAmount: bidPrice,
				Owner:     owner,
				AuctionId: id,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func ExecuteAuctionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec [auction-id] --from [sender]",
		Args:  cobra.ExactArgs(1),
		Short: "execute pending auction",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s exec <auction-id> --from <sender> --chain-id <chain-id>`, version.AppName, auctiontypes.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if args[0] == "" {
				return fmt.Errorf("contract-id ")
			}

			fmt.Printf("Auction Id :: %s", args[0])

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress().String()

			msg := auctiontypes.MsgExecAuction{
				AuctionId: id,
				Sender:    sender,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CancelAuctionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-auction [auction-id] --from [sender]",
		Args:  cobra.ExactArgs(1),
		Short: "cancel an auction",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s cancel-auction <auction-id> --from <sender> --chain-id <chain-id>`, version.AppName, auctiontypes.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if args[0] == "" {
				return fmt.Errorf("auction-id cannot be empty")
			}

			fmt.Printf("Auction Id :: %s", args[0])

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress().String()

			msg := auctiontypes.MsgCancelAuction{
				AuctionId: id,
				Sender:    sender,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
