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
	owner := sdk.MustAccAddressFromBech32(msg.Owner)

	var md at.AuctionMetadata
	err := ms.k.cdc.UnpackAny(msg.GetAuctionMetadata(), &md)
	if err != nil {
		return &at.MsgNewAuctionResponse{}, fmt.Errorf("error serializing auction metadata")
	}

	auction, err := ms.k.CreateAuction(goCtx, msg.AuctionType, owner, md)
	if err != nil {
		return &at.MsgNewAuctionResponse{}, fmt.Errorf("error creating auction")
	}

	err = ms.k.bk.SendCoinsFromAccountToModule(goCtx, owner, at.ModuleName, msg.Deposit)
	if err != nil {
		return &at.MsgNewAuctionResponse{}, fmt.Errorf("error crediting auction deposit")
	}

	ms.k.Logger().Info(auction.String())
	err = ms.k.Auctions.Set(goCtx, auction.GetId(), auction)
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
	oa.Ids = append(oa.Ids, auction.GetId())

	// Set Auctions by Owner
	err = ms.k.OwnerAuctions.Set(goCtx, owner, oa)
	if err != nil {
		return &at.MsgNewAuctionResponse{}, err
	}

	// Push auction to ActiveAuction Queue
	err = ms.k.ActiveAuctions.Set(goCtx, auction.GetId())
	if err != nil {
		return &at.MsgNewAuctionResponse{}, err
	}

	return &at.MsgNewAuctionResponse{
		Id: auction.GetId(),
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

		auction, err = ms.k.SubmitBid(ctx, auction.GetType(), auction, msg)
		if err != nil {
			return &at.MsgNewBidResponse{}, fmt.Errorf("error creating auction")
		}

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
