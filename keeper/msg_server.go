package keeper

import (
	"context"

	"github.com/fatal-fruit/auction/types"
)

type msgServer struct {
	k Keeper
}

var _ types.MsgServer = msgServer{}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}

func (ms msgServer) CreateAuction(ctx context.Context, msg *types.MsgCreateAuction) (*types.MsgCreateAuctionResponse, error) {
	return nil, nil
}
