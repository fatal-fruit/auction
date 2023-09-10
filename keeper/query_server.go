package keeper

import (
	"context"
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

func (qs queryServer) Counter(goCtx context.Context, r *auctiontypes.QueryCounterRequest) (*auctiontypes.QueryCounterResponse, error) {
	return &auctiontypes.QueryCounterResponse{}, nil
}
