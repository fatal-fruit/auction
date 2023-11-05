package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auctiontypes "github.com/fatal-fruit/auction/types"
)

type msgServer struct {
	k Keeper
}

var _ auctiontypes.MsgServer = msgServer{}

func NewMsgServerImpl(keeper Keeper) auctiontypes.MsgServer {
	return &msgServer{k: keeper}
}

func (ms msgServer) NewAuction(goCtx context.Context, msg *auctiontypes.MsgNewAuction) (*auctiontypes.MsgNewAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Get Next Id
	id, err := ms.k.IDs.Next(goCtx)
	if err != nil {
		return &auctiontypes.MsgNewAuctionResponse{}, fmt.Errorf("error creating id for auction")
	}

	owner := sdk.MustAccAddressFromBech32(msg.Owner)

	// Generate start/end time
	start := ctx.BlockTime()
	end := start.Add(msg.Duration)

	err = ms.k.bk.SendCoinsFromAccountToModule(goCtx, owner, auctiontypes.ModuleName, msg.Deposit)
	if err != nil {
		return &auctiontypes.MsgNewAuctionResponse{}, fmt.Errorf("error crediting auction deposit")
	}

	// Generate escrow contract
	contractId, err := ms.k.es.NewContract()
	if err != nil {
		// TODO: Rollback deposit
		return &auctiontypes.MsgNewAuctionResponse{}, fmt.Errorf("error creating escrow contract for auction")
	}

	auction := auctiontypes.ReserveAuction{
		Id:             id,
		Owner:          owner.String(),
		AuctionType:    msg.AuctionType,
		EscrowContract: contractId,
		ReservePrice:   msg.ReservePrice,
		StartTime:      start,
		EndTime:        end,
		Bids:           []*auctiontypes.Bid{},
	}

	ms.k.Logger(ctx).Info(auction.String())
	err = ms.k.Auctions.Set(goCtx, id, auction)
	if err != nil {
		// TODO: Rollback deposit
		return &auctiontypes.MsgNewAuctionResponse{}, fmt.Errorf("error creating auction")
	}

	hasAuctions, err := ms.k.OwnerAuctions.Has(goCtx, owner)
	if err != nil {
		return &auctiontypes.MsgNewAuctionResponse{}, err
	}
	var oa auctiontypes.OwnerAuctions
	if hasAuctions {
		oa, err = ms.k.OwnerAuctions.Get(goCtx, owner)
		if err != nil {
			return &auctiontypes.MsgNewAuctionResponse{}, err

		}
	}
	oa.Ids = append(oa.Ids, id)

	// Set Auctions by Owner
	err = ms.k.OwnerAuctions.Set(goCtx, owner, oa)
	if err != nil {
		return &auctiontypes.MsgNewAuctionResponse{}, err
	}

	// Push auction to ActiveAuction Queue
	err = ms.k.ActiveAuctions.Set(goCtx, id)
	if err != nil {
		return &auctiontypes.MsgNewAuctionResponse{}, err
	}

	return &auctiontypes.MsgNewAuctionResponse{
		Id: id,
	}, nil
}

func (ms msgServer) NewBid(goCtx context.Context, msg *auctiontypes.MsgNewBid) (*auctiontypes.MsgNewBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// get auction from active auctions
	hasAuction, err := ms.k.ActiveAuctions.Has(goCtx, msg.GetAuctionId())
	if err != nil {
		return &auctiontypes.MsgNewBidResponse{}, err
	}
	if hasAuction {
		auction, err := ms.k.Auctions.Get(goCtx, msg.GetAuctionId())
		if err != nil {
			return &auctiontypes.MsgNewBidResponse{}, err
		}

		// Validate bid price is comepetitive
		if msg.Bid.IsLT(auction.ReservePrice) {
			return &auctiontypes.MsgNewBidResponse{}, fmt.Errorf("invalid bid price")
		}

		// Validate auction is active
		if ctx.BlockTime().After(auction.EndTime) {
			return &auctiontypes.MsgNewBidResponse{}, fmt.Errorf("expired auction")
		}

		auction.Bids = append(auction.Bids, &auctiontypes.Bid{
			AuctionId: msg.AuctionId,
			Bidder:    msg.Owner,
			BidPrice:  msg.Bid,
			Timestamp: ctx.BlockTime(),
		})

		err = ms.k.Auctions.Set(goCtx, auction.GetId(), auction)
		if err != nil {
			return &auctiontypes.MsgNewBidResponse{}, err
		}
	} else {
		return &auctiontypes.MsgNewBidResponse{}, fmt.Errorf("invalid auction id")
	}

	return &auctiontypes.MsgNewBidResponse{}, nil
}
