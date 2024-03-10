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

// NewContractCmd creates a CLI command for MsgNewContract.
func NewAuctionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-auction [auction-json-file] [deposit] --from [sender]",
		Args:  cobra.ExactArgs(2),
		Short: "create new auction",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s create-auction <deposit> auction_metadata.json --from <sender> --chain-id <chain-id>

		Where auction_metadata.json contains:
		{
			"@type": "/fatal_fruit.auction.v1.ReserveAuctionMetadata",
			"threshold": "1",
			"windows": {
				"voting_period": "120h",
				"min_execution_period": "0s"
			}
		}
		`, version.AppName, auctiontypes.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if args[0] == "" || args[1] == "" {
				return fmt.Errorf("auction metadata and deposit cannot be empty")
			}

			fmt.Printf("Auction Metadata :: %s", args[0])
			fmt.Printf("Deposit :: %s", args[1])
			fmt.Printf("Owner :: %s", clientCtx.GetFromAddress().String())

			auctionMetadata, err := parseAuctionMetadata(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}
			fmt.Printf("Auction :: %s", auctionMetadata)

			// Parse Deposit
			deposit, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return err
			}

			// TODO: Validate auction
			//if err := auction.ValidateBasic(); err != nil {
			//	return err
			//}

			// Parse Reserve Price
			//rp, err := sdk.ParseCoinsNormalized(args[0])
			//if err != nil {
			//	return err
			//}
			//found, reservePrice := rp.Find(sdk.DefaultBondDenom)
			//if !found {
			//	return fmt.Errorf("Invalid reserve price")
			//}

			// Parse Auction Duration
			//seconds, err := strconv.Atoi(args[1])
			//if err != nil {
			//	return fmt.Errorf("invalid duration")
			//}
			//duration := time.Duration(seconds) * time.Second

			owner := clientCtx.GetFromAddress().String()

			msg := auctiontypes.MsgNewAuction{
				Owner:   owner,
				Deposit: deposit,
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
				return fmt.Errorf("contract-id and deposit cannot be empty")
			}

			fmt.Printf("Auction Id :: %s", args[0])
			fmt.Printf("Bid Price :: %s", args[1])

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			// Parse Price
			bp, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}
			found, bidPrice := bp.Find(sdk.DefaultBondDenom)
			if !found {
				return fmt.Errorf("Invalid bid price")
			}

			owner := clientCtx.GetFromAddress().String()

			msg := auctiontypes.MsgNewBid{
				Bid:       bidPrice,
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
