package keeper

import "context"

func (k *Keeper) ExportGenesis(ctx context.Context) (data *types.GenesisState) {
	return &types.GenesisState{}
}

func (k *Keeper) InitGenesis(ctx context.Context, genState *types.GenesisState) {

}
