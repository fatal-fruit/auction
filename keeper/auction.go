package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auctiontypes "github.com/fatal-fruit/auction/types"
)

func (k *Keeper) CreateAuction(ctx context.Context, auctionType string, owner sdk.AccAddress, md auctiontypes.AuctionMetadata) (auctiontypes.Auction, error) {
	// Check if keeper has registered auction type
	if !k.resolver.HasType(auctionType) {
		return nil, fmt.Errorf("proposal type %s is not registered", auctionType)
	}

	handler := k.resolver.GetHandler(auctionType)

	// Get Next Id
	id, err := k.IDs.Next(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating id for auction")
	}

	auction, err := handler.CreateAuction(ctx, id, md)
	if err != nil {
		return nil, fmt.Errorf("error creating auction")
	}
	auction.SetOwner(owner)

	return auction, nil
}

func (k *Keeper) ExecuteAuction(ctx context.Context, auction auctiontypes.Auction) error {
	if !k.resolver.HasType(auction.GetType()) {
		return fmt.Errorf("proposal type %s is not registered", auction.GetType())
	}

	handler := k.resolver.GetHandler(auction.GetType())

	return handler.ExecAuction(ctx, auction)
}
