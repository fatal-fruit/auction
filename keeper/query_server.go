package keeper

import (
	"context"
	"fmt"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
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

	aa, err := codectypes.NewAnyWithValue(auction)
	if err != nil {
		return &auctiontypes.QueryAuctionResponse{}, err
	}

	return &auctiontypes.QueryAuctionResponse{
		Auction: aa,
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

	anyAuctions := make([]*codectypes.Any, 0, len(ownerAuctions.Ids))
	for _, id := range ownerAuctions.Ids {
		a, err := qs.k.Auctions.Get(goCtx, id)
		if err != nil {
			return &auctiontypes.QueryOwnerAuctionsResponse{}, fmt.Errorf(fmt.Sprintf("unable to retrieve owner auctions with address :: %s", ownerAddress))
		}

		aa, err := codectypes.NewAnyWithValue(a)
		if err != nil {
			return &auctiontypes.QueryOwnerAuctionsResponse{}, fmt.Errorf(fmt.Sprintf("unable to retrieve owner auctions with address :: %s", ownerAddress))
		}
		anyAuctions = append(anyAuctions, aa)
	}

	return &auctiontypes.QueryOwnerAuctionsResponse{
		Auctions: anyAuctions,
	}, nil
}

func (qs queryServer) AllAuctions(ctx context.Context, _ *auctiontypes.QueryAllAuctionsRequest) (*auctiontypes.QueryAllAuctionsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	auctions := qs.k.GetAllAuctions(sdkCtx)

	// Convert to slice of pointers
	auctionsPtrs := make([]*auctiontypes.ReserveAuction, len(auctions))
	for i, auction := range auctions {
		auctionsPtrs[i] = &auction
	}

	return &auctiontypes.QueryAllAuctionsResponse{Auctions: auctionsPtrs}, nil
}
