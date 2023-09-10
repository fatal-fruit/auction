package keeper

import (
	"context"
	auctiontypes "github.com/fatal-fruit/auction/types"
)

func (k *Keeper) InitGenesis(ctx context.Context, data *auctiontypes.GenesisState) error {
	// TODO: Implement
	return nil
}

func (k *Keeper) ExportGenesis(ctx context.Context) (*auctiontypes.GenesisState, error) {
	// TODO: Implement
	return &auctiontypes.GenesisState{}, nil
}
