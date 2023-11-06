package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auctiontypes "github.com/fatal-fruit/auction/types"
)

type Strategy interface {
	ExecuteStrategy(context.Context, auctiontypes.ReserveAuction, auctiontypes.EscrowService, auctiontypes.BankKeeper) error
}

type SettleStrategy struct {
	*auctiontypes.SettleStrategy
}

func (s *SettleStrategy) ExecuteStrategy(ctx context.Context, auction auctiontypes.ReserveAuction, es auctiontypes.EscrowService, bk auctiontypes.BankKeeper) error {
	// Select Winner
	winningBid, err := GetWinner(auction)
	if err != nil {
		return err
	}

	bidder := sdk.MustAccAddressFromBech32(winningBid.Bidder)
	auctioneer := sdk.MustAccAddressFromBech32(auction.Owner)

	// Send bid amount to auction owner
	err = bk.SendCoins(ctx, bidder, auctioneer, sdk.Coins{winningBid.BidPrice})
	if err != nil {
		return err
	}

	// Release escrowed auction bounty
	escrowAddr := sdk.MustAccAddressFromBech32(s.EscrowContractAddress)
	err = es.Release(ctx, s.EscrowContractId, escrowAddr, bidder)
	if err != nil {
		return err
	}
	return nil
}

func (s *SettleStrategy) ToProto() *auctiontypes.SettleStrategy {
	return &auctiontypes.SettleStrategy{
		StrategyType:          s.GetStrategyType(),
		EscrowContractId:      s.GetEscrowContractId(),
		EscrowContractAddress: s.GetEscrowContractAddress(),
	}
}

// Use generics
func BuildSettleStrategy(ctx context.Context, es auctiontypes.EscrowService) (*SettleStrategy, error) {
	contract, err := es.NewContract(ctx, 1)
	if err != nil {
		return nil, fmt.Errorf("Unable to create escrow contract")
	}
	s := &auctiontypes.SettleStrategy{
		StrategyType:          auctiontypes.SETTLE,
		EscrowContractId:      contract.GetId(),
		EscrowContractAddress: contract.GetAddress().String(),
	}
	return &SettleStrategy{s}, nil

}
