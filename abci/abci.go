package abci

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/fatal-fruit/auction/keeper"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) error {
	k.ProcessActiveAuctions(ctx)
	return nil
}
