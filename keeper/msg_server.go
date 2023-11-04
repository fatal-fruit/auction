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
	err = ms.k.OwnerAuctions.Set(goCtx, owner, oa)

	return &auctiontypes.MsgNewAuctionResponse{
		Id: id,
	}, nil
}
