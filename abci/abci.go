package abci

import (
	"context"
	"cosmossdk.io/log"
	"github.com/fatal-fruit/auction/keeper"
)

func EndBlocker(ctx context.Context, k keeper.Keeper, log log.Logger) error {
	logger := log
	logger.Info("EndBlocker :: Processing Active Auctions")
	err := k.ProcessActiveAuctions(ctx)
	if err != nil {
		return err
	}

	logger.Info("EndBlocker :: Processing Expired Auctions")
	err = k.ProcessExpiredAuctions(ctx)
	if err != nil {
		return err
	}

	logger.Info("EndBlocker :: Calculating Pending Auctions")
	err = k.GetPending(ctx)
	if err != nil {
		return err
	}

	logger.Info("EndBlocker :: Done")

	return nil
}
