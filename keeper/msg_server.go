package keeper

import (
	"context"
	"fmt"
	"log"

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
	owner := sdk.MustAccAddressFromBech32(msg.Owner)

	var md at.AuctionMetadata
	err := ms.k.cdc.UnpackAny(msg.GetAuctionMetadata(), &md)
	if err != nil {
		return &at.MsgNewAuctionResponse{}, fmt.Errorf("error serializing auction metadata: %v", err)
	}

	types := ms.k.Resolver.ListTypes()

	auction, err := ms.k.CreateAuction(goCtx, msg.AuctionType, owner, md)
	if err != nil {
		return &at.MsgNewAuctionResponse{}, fmt.Errorf("error creating auction: %v, currenttype: %v, available types: %v", err, msg.AuctionType, types)
	}

	err = ms.k.bk.SendCoinsFromAccountToModule(goCtx, owner, at.ModuleName, msg.Deposit)
	if err != nil {
		return &at.MsgNewAuctionResponse{}, fmt.Errorf("failed to credit auction deposit for owner %s in module %s with deposit %s: %v", owner, at.ModuleName, msg.Deposit, err)
	}

	ms.k.Logger().Info(auction.String())
	err = ms.k.Auctions.Set(goCtx, auction.GetId(), auction)
	if err != nil {
		// TODO: Rollback deposit
		return &at.MsgNewAuctionResponse{}, fmt.Errorf("error creating auction: %v", err)
	}

	hasAuctions, err := ms.k.OwnerAuctions.Has(goCtx, owner)
	if err != nil {
		return &at.MsgNewAuctionResponse{}, fmt.Errorf("error retrieving owner auctions: %v", err)
	}

	var oa at.OwnerAuctions

	if !hasAuctions {
		oa = at.OwnerAuctions{Ids: []uint64{}}
	}

	if hasAuctions {
		oa, err = ms.k.OwnerAuctions.Get(goCtx, owner)
		if err != nil {
			return &at.MsgNewAuctionResponse{}, fmt.Errorf("error retrieving owner auctions: %v", err)
		}
	}

	oa.Ids = append(oa.Ids, auction.GetId())

	// Set Auctions by Owner
	err = ms.k.OwnerAuctions.Set(goCtx, owner, oa)
	if err != nil {
		return &at.MsgNewAuctionResponse{}, fmt.Errorf("error setting owner auctions: %v", err)
	}

	// Push auction to ActiveAuction Queue
	err = ms.k.ActiveAuctions.Set(goCtx, auction.GetId())
	if err != nil {
		return &at.MsgNewAuctionResponse{}, fmt.Errorf("error setting Active auctions: %v", err)
	}

	id := auction.GetId()

	return &at.MsgNewAuctionResponse{
		Id: id,
	}, nil
}

func (ms msgServer) StartAuction(goCtx context.Context, msg *at.MsgStartAuction) (*at.MsgStartAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var auction at.Auction
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
	// TODO: Pass context instead
	auction.StartAuction(ctx.BlockTime())

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
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in NewBid: %v", r)
		}
	}()

	ctx := sdk.UnwrapSDKContext(goCtx)
	owner := sdk.MustAccAddressFromBech32(msg.Owner)

	if err := ms.k.Validate(); err != nil {
    return nil, fmt.Errorf("keeper validation error: %v", err)
    }

	hasAuction, err := ms.k.ActiveAuctions.Has(goCtx, msg.GetAuctionId())
	if err != nil {
		return &at.MsgNewBidResponse{}, fmt.Errorf("error checking active auctions for auction ID %d: %v", msg.GetAuctionId(), err)
	}
	if !hasAuction {
		return &at.MsgNewBidResponse{}, fmt.Errorf("invalid auction ID %d: auction does not exist or is not active", msg.GetAuctionId())
	}
	auction, err := ms.k.Auctions.Get(goCtx, msg.GetAuctionId())
	if err != nil {
		return &at.MsgNewBidResponse{}, fmt.Errorf("error retrieving auction with ID %d: %v", msg.GetAuctionId(), err)
	}

	if auction == nil {
		log.Printf("Auction object is nil")
		return &at.MsgNewBidResponse{}, fmt.Errorf("auction object is nil")
	}

	if auction.GetType() == "" {
		log.Printf("Auction type is empty")
		return &at.MsgNewBidResponse{}, fmt.Errorf("auction type is empty")
	}

	if msg == nil {
		log.Printf("Message object is nil")
		return &at.MsgNewBidResponse{}, fmt.Errorf("msg object is nil")
	}

	err = ms.k.bk.SendCoinsFromAccountToModule(ctx, owner, at.ModuleName, sdk.NewCoins(msg.BidAmount))
	if err != nil {
		return &at.MsgNewBidResponse{}, fmt.Errorf("failed to transfer funds: %v", err)
	}

	auction, err = ms.k.SubmitBid(ctx, auction.GetType(), auction, msg)
	if err != nil {
		return &at.MsgNewBidResponse{}, fmt.Errorf("error submitting bid for auction type %v: %v", auction.GetType(), err)
	}

	err = ms.k.Auctions.Set(goCtx, auction.GetId(), auction)
	if err != nil {
		return &at.MsgNewBidResponse{}, fmt.Errorf("error updating auction with ID %d in store: %v", auction.GetId(), err)
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

	err = ms.k.ExecuteAuction(goCtx, auction)
	if err != nil {
		return &at.MsgExecAuctionResponse{}, err
	}

	// remove from pending
	err = ms.k.PendingAuctions.Remove(goCtx, msg.GetAuctionId())
	if err != nil {
		return &at.MsgExecAuctionResponse{}, err
	}
	//update status
	auction.UpdateStatus(at.CLOSED)
	err = ms.k.Auctions.Set(goCtx, auction.GetId(), auction)
	if err != nil {
		return &at.MsgExecAuctionResponse{}, err
	}

	return &at.MsgExecAuctionResponse{}, nil
}
