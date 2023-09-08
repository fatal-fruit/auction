package keeper

import (
	types "auction/types"
	"context"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (srv msgServer) MsgCreateAuctionMessage(ctx context.Context, msg *typesMsgCreate) (*types.MsgTripCircuitBreakerResponse, error) {

}

func (srv msgServer) MsgUpdateAuctionMessage(ctx context.Context, msg *typesMsgCreate) (*types.MsgTripCircuitBreakerResponse, error) {

}

func (srv msgServer) MsgDeleteAuctionMessage(ctx context.Context, msg *typesMsgCreate) (*types.MsgTripCircuitBreakerResponse, error) {

}

func (srv msgServer) MsgExecuteAuctionMessage(ctx context.Context, msg *typesMsgCreate) (*types.MsgTripCircuitBreakerResponse, error) {

}
