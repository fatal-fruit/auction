package keeper

import (
	"context"
)

type msgServer struct {
	k Keeper
}

func NewMsgServerImpl(keeper Keeper) {
}

func (ms msgServer) Reserve(goCtx context.Context) {

}
