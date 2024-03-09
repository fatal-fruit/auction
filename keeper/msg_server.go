package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	at "github.com/fatal-fruit/auction/types"
)

type msgServer struct {
	k Keeper
}

var _ at.MsgServer = msgServer{}

func NewMsgServerImpl(keeper Keeper) at.MsgServer {
	return &msgServer{k: keeper}
}

func (ms msgServer) NewAuction(goCtx context.Context, msg *at.MsgNewAuction) (*at.MsgNewAuctionResponse, error) {
	// Get Next Id
	id, err := ms.k.IDs.Next(goCtx)
	if err != nil {
		return &at.MsgNewAuctionResponse{}, fmt.Errorf("error creating id for auction")
	}

	owner := sdk.MustAccAddressFromBech32(msg.Owner)

	err = ms.k.bk.SendCoinsFromAccountToModule(goCtx, owner, at.ModuleName, msg.Deposit)
	if err != nil {
		return &at.MsgNewAuctionResponse{}, fmt.Errorf("error crediting auction deposit")
	}

	strategy, err := BuildSettleStrategy(goCtx, ms.k.es)
	if err != nil {
		// TODO: Rollback deposit
		return &at.MsgNewAuctionResponse{}, fmt.Errorf("error creating escrow contract for auction")
	}
	auction := at.ReserveAuction{
		Id:           id,
		Status:       at.ACTIVE,
		Owner:        owner.String(),
		AuctionType:  msg.AuctionType,
		Duration:     msg.Duration,
		ReservePrice: msg.ReservePrice,
		Bids:         []*at.Bid{},
		Strategy:     strategy.ToProto(),
	}

	ms.k.Logger().Info(auction.String())
	err = ms.k.Auctions.Set(goCtx, id, auction)
	if err != nil {
		// TODO: Rollback deposit
		return &at.MsgNewAuctionResponse{}, fmt.Errorf("error creating auction")
	}

	hasAuctions, err := ms.k.OwnerAuctions.Has(goCtx, owner)
	if err != nil {
		return &at.MsgNewAuctionResponse{}, err
	}
	var oa at.OwnerAuctions
	if hasAuctions {
		oa, err = ms.k.OwnerAuctions.Get(goCtx, owner)
		if err != nil {
			return &at.MsgNewAuctionResponse{}, err

		}
	}
	oa.Ids = append(oa.Ids, id)

	// Set Auctions by Owner
	err = ms.k.OwnerAuctions.Set(goCtx, owner, oa)
	if err != nil {
		return &at.MsgNewAuctionResponse{}, err
	}

	// Push auction to ActiveAuction Queue
	err = ms.k.ActiveAuctions.Set(goCtx, id)
	if err != nil {
		return &at.MsgNewAuctionResponse{}, err
	}

	return &at.MsgNewAuctionResponse{
		Id: id,
	}, nil
}

func (ms msgServer) StartAuction(goCtx context.Context, msg *at.MsgStartAuction) (*at.MsgStartAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var auction at.ReserveAuction
	hasAuctions, err := ms.k.Auctions.Has(goCtx, msg.Id)
	if err != nil {
		return &at.MsgStartAuctionResponse{}, err
	}
	if hasAuctions {
		auction, err = ms.k.Auctions.Get(goCtx, msg.GetId())
	}
	if err != nil {
		return &at.MsgStartAuctionResponse{}, err
	}

	// Generate start/end time
	start := ctx.BlockTime()
	end := start.Add(auction.Duration)

	// Set Start and End time for auction
	auction.StartTime = start
	auction.EndTime = end

	// Save updated auction
	err = ms.k.Auctions.Set(goCtx, auction.GetId(), auction)
	if err != nil {
		return &at.MsgStartAuctionResponse{}, err
	}

	// Push auction to ActiveAuction Queue
	err = ms.k.ActiveAuctions.Set(goCtx, msg.Id)
	if err != nil {
		return &at.MsgStartAuctionResponse{}, err
	}

	return &at.MsgStartAuctionResponse{}, nil
}

func (ms msgServer) NewBid(goCtx context.Context, msg *at.MsgNewBid) (*at.MsgNewBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// get auction from active auctions
	hasAuction, err := ms.k.ActiveAuctions.Has(goCtx, msg.GetAuctionId())
	if err != nil {
		return &at.MsgNewBidResponse{}, err
	}
	if hasAuction {
		auction, err := ms.k.Auctions.Get(goCtx, msg.GetAuctionId())
		if err != nil {
			return &at.MsgNewBidResponse{}, err
		}

		// Validate bid price is over Reserve Price
		if msg.Bid.IsLT(auction.ReservePrice) {
			return &at.MsgNewBidResponse{}, fmt.Errorf("bid lower than reserve price")
		}

		// Validate auction is active
		if ctx.BlockTime().After(auction.EndTime) {
			return &at.MsgNewBidResponse{}, fmt.Errorf("expired auction")
		}

		// Validate bid price is competitive
		if len(auction.Bids) > 0 && msg.Bid.IsLTE(auction.LastPrice) {
			return &at.MsgNewBidResponse{}, fmt.Errorf("bid lower than latest price")
		}

		auction.Bids = append(auction.Bids, &at.Bid{
			AuctionId: msg.AuctionId,
			Bidder:    msg.Owner,
			BidPrice:  msg.Bid,
			Timestamp: ctx.BlockTime(),
		})

		auction.LastPrice = msg.Bid

		err = ms.k.Auctions.Set(goCtx, auction.GetId(), auction)
		if err != nil {
			return &at.MsgNewBidResponse{}, err
		}
	} else {
		return &at.MsgNewBidResponse{}, fmt.Errorf("invalid auction id")
	}

	return &at.MsgNewBidResponse{}, nil
}

func (ms msgServer) Exec(goCtx context.Context, msg *at.MsgExecAuction) (*at.MsgExecAuctionResponse, error) {
	// Check auction is in pending
	isPending, err := ms.k.PendingAuctions.Has(goCtx, msg.GetAuctionId())
	if err != nil {
		return &at.MsgExecAuctionResponse{}, err
	}
	if !isPending {
		return &at.MsgExecAuctionResponse{}, fmt.Errorf("auction is not executable")
	}
	//
	auction, err := ms.k.Auctions.Get(goCtx, msg.GetAuctionId())
	if err != nil {
		return &at.MsgExecAuctionResponse{}, err
	}

	// execute strategy
	exeuctionStrat := SettleStrategy{auction.Strategy}
	err = exeuctionStrat.ExecuteStrategy(goCtx, auction, ms.k.es, ms.k.bk)
	if err != nil {
		return &at.MsgExecAuctionResponse{}, err
	}
	//auction.Strategy

	// remove from pending
	err = ms.k.PendingAuctions.Remove(goCtx, msg.GetAuctionId())
	if err != nil {
		return &at.MsgExecAuctionResponse{}, err
	}
	//update status
	auction.Status = at.CLOSED
	err = ms.k.Auctions.Set(goCtx, auction.GetId(), auction)
	if err != nil {
		return &at.MsgExecAuctionResponse{}, err
	}

	return &at.MsgExecAuctionResponse{}, nil
}
