package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auctiontypes "github.com/fatal-fruit/auction/types"
)

var _ auctiontypes.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k Keeper) auctiontypes.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

func (qs queryServer) Auction(goCtx context.Context, r *auctiontypes.QueryAuctionRequest) (*auctiontypes.QueryAuctionResponse, error) {
	auction, err := qs.k.Auctions.Get(goCtx, r.GetId())
	if err != nil {
		return &auctiontypes.QueryAuctionResponse{}, fmt.Errorf(fmt.Sprintf("unable to retrieve auction with id :: %d", r.GetId()))
	}
	return &auctiontypes.QueryAuctionResponse{
		Auction: &auction,
	}, nil
}

func (qs queryServer) OwnerAuctions(goCtx context.Context, r *auctiontypes.QueryOwnerAuctionsRequest) (*auctiontypes.QueryOwnerAuctionsResponse, error) {
	ownerAddress, err := sdk.AccAddressFromBech32(r.GetOwnerAddress())
	if err != nil {
		return &auctiontypes.QueryOwnerAuctionsResponse{}, fmt.Errorf(fmt.Sprintf("unable to retrieve owner address :: %s", r.GetOwnerAddress()))

	}
	ownerAuctions, err := qs.k.OwnerAuctions.Get(goCtx, ownerAddress)
	if err != nil {
		return &auctiontypes.QueryOwnerAuctionsResponse{}, fmt.Errorf(fmt.Sprintf("unable to retrieve owner auctions with address :: %s", ownerAddress))
	}

	var auctions []*auctiontypes.ReserveAuction

	for _, id := range ownerAuctions.Ids {
		a, err := qs.k.Auctions.Get(goCtx, id)
		if err != nil {
			return &auctiontypes.QueryOwnerAuctionsResponse{}, fmt.Errorf(fmt.Sprintf("unable to retrieve owner auctions with address :: %s", ownerAddress))
		}
		auctions = append(auctions, &a)
	}

	return &auctiontypes.QueryOwnerAuctionsResponse{
		auctions,
	}, nil
}
