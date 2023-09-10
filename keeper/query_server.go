package keeper

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k Keeper) {

}

type queryServer struct {
	k Keeper
}
