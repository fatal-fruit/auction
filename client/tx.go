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
	"time"
)

// NewContractCmd creates a CLI command for MsgNewContract.
func NewAuctionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-auction [reserve-price] [duration] [deposit] --from [sender]",
		Args:  cobra.ExactArgs(3),
		Short: "create new auction",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s create-contract <reserve-price> <duration> <deposit> --from <sender> --chain-id <chain-id>`, version.AppName, auctiontypes.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if args[0] == "" || args[1] == "" || args[2] == "" {
				return fmt.Errorf("reserve-price, deposit, and duration cannot be empty")
			}

			fmt.Println(fmt.Sprintf("Reserve Price :: %s", args[0]))
			fmt.Println(fmt.Sprintf("Duration :: %w", args[1]))
			fmt.Println(fmt.Sprintf("Deposit :: %w", args[2]))
			fmt.Println(fmt.Sprintf("Owner :: %w", clientCtx.GetFromAddress().String()))

			// Parse Reserve Price
			rp, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}
			found, reservePrice := rp.Find(sdk.DefaultBondDenom)
			if !found {
				return fmt.Errorf("Invalid reserve price")
			}

			// Parse Deposit
			deposit, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return err
			}

			// Parse Auction Duration
			seconds, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid duration")
			}
			duration := time.Duration(seconds) * time.Second

			owner := clientCtx.GetFromAddress().String()

			msg := auctiontypes.MsgNewAuction{
				Owner:        owner,
				Deposit:      deposit,
				ReservePrice: reservePrice,
				Duration:     duration,
				AuctionType:  auctiontypes.RESERVE,
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

			fmt.Println(fmt.Sprintf("Auction Id :: %s", args[0]))
			fmt.Println(fmt.Sprintf("Bid Price :: %w", args[1]))

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

			fmt.Println(fmt.Sprintf("Auction Id :: %s", args[0]))

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

			fmt.Println(fmt.Sprintf("Auction Id :: %s", args[0]))

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