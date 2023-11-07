package abci

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/fatal-fruit/auction/keeper"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) error {
	err := k.ProcessActiveAuctions(ctx)
	if err != nil {
		return err
	}

	err = k.ProcessExpiredAuctions(ctx)
	if err != nil {
		return err
	}
	return nil
}
