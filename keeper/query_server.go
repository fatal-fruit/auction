package keeper

//var _ auctiontypes.QueryServer = queryServer{}
//
//// NewQueryServerImpl returns an implementation of the module QueryServer.
//func NewQueryServerImpl(k Keeper) auctiontypes.QueryServer {
//	return queryServer{k}
//}

type queryServer struct {
	k Keeper
}

//func (qs queryServer) Auctions(goCtx context.Context, r *auctiontypes.QueryNameRequest) (*auctiontypes.QueryAuction, error) {
//
//}
