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

	strategy, err := BuildSettleStrategy(goCtx, ms.k.es)
	if err != nil {
		// TODO: Rollback deposit
		return &auctiontypes.MsgNewAuctionResponse{}, fmt.Errorf("error creating escrow contract for auction")
	}
	auction := auctiontypes.ReserveAuction{
		Id:           id,
		Status:       auctiontypes.ACTIVE,
		Owner:        owner.String(),
		AuctionType:  msg.AuctionType,
		ReservePrice: msg.ReservePrice,
		StartTime:    start,
		EndTime:      end,
		Bids:         []*auctiontypes.Bid{},
		Strategy:     strategy.ToProto(),
	}

	ms.k.Logger().Info(auction.String())
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

		// Validate bid price is over Reserve Price
		if msg.Bid.IsLT(auction.ReservePrice) {
			return &auctiontypes.MsgNewBidResponse{}, fmt.Errorf("bid lower than reserve price")
		}

		// Validate auction is active
		if ctx.BlockTime().After(auction.EndTime) {
			return &auctiontypes.MsgNewBidResponse{}, fmt.Errorf("expired auction")
		}

		// Validate bid price is competitive
		if len(auction.Bids) > 0 && msg.Bid.IsLTE(auction.LastPrice) {
			return &auctiontypes.MsgNewBidResponse{}, fmt.Errorf("bid lower than latest price")
		}

		auction.Bids = append(auction.Bids, &auctiontypes.Bid{
			AuctionId: msg.AuctionId,
			Bidder:    msg.Owner,
			BidPrice:  msg.Bid,
			Timestamp: ctx.BlockTime(),
		})

		auction.LastPrice = msg.Bid

		err = ms.k.Auctions.Set(goCtx, auction.GetId(), auction)
		if err != nil {
			return &auctiontypes.MsgNewBidResponse{}, err
		}
	} else {
		return &auctiontypes.MsgNewBidResponse{}, fmt.Errorf("invalid auction id")
	}

	return &auctiontypes.MsgNewBidResponse{}, nil
}

func (ms msgServer) Exec(goCtx context.Context, msg *auctiontypes.MsgExecAuction) (*auctiontypes.MsgExecAuctionResponse, error) {
	// Check auction is in pending
	isPending, err := ms.k.PendingAuctions.Has(goCtx, msg.GetAuctionId())
	if err != nil {
		return &auctiontypes.MsgExecAuctionResponse{}, err
	}
	if !isPending {
		return &auctiontypes.MsgExecAuctionResponse{}, fmt.Errorf("auction is not executable")
	}
	//
	auction, err := ms.k.Auctions.Get(goCtx, msg.GetAuctionId())
	if err != nil {
		return &auctiontypes.MsgExecAuctionResponse{}, err
	}

	// execute strategy
	exeuctionStrat := SettleStrategy{auction.Strategy}
	err = exeuctionStrat.ExecuteStrategy(goCtx, auction, ms.k.es, ms.k.bk)
	//auction.Strategy

	// remove from pending
	err = ms.k.PendingAuctions.Remove(goCtx, msg.GetAuctionId())
	if err != nil {
		return &auctiontypes.MsgExecAuctionResponse{}, err
	}
	//update status
	auction.Status = auctiontypes.CLOSED
	err = ms.k.Auctions.Set(goCtx, auction.GetId(), auction)
	if err != nil {
		return &auctiontypes.MsgExecAuctionResponse{}, err
	}

	return &auctiontypes.MsgExecAuctionResponse{}, nil
}
