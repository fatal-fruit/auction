package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (k *Keeper) ExportGenesis(ctx context.Context) (data *types.GenesisState) {
	return &types.GenesisState{}
}

func (k *Keeper) InitGenesis(ctx context.Context, genState *types.GenesisState) {

}
