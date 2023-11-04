package integration_test

import (
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/testutil/configurator"
	// blank import for app wiring registration
	_ "github.com/cosmos/cosmos-sdk/x/auth"
	_ "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	_ "github.com/cosmos/cosmos-sdk/x/bank"
	_ "github.com/cosmos/cosmos-sdk/x/consensus"
	_ "github.com/cosmos/cosmos-sdk/x/genutil"
	_ "github.com/cosmos/cosmos-sdk/x/mint"
	_ "github.com/cosmos/cosmos-sdk/x/staking"
	_ "github.com/fatal-fruit/auction/module"

	cosmosapp "cosmossdk.io/api/cosmos/app/v1alpha1"
	"cosmossdk.io/core/appconfig"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	auctionmodule "github.com/fatal-fruit/auction/api/module/v1"
	"github.com/fatal-fruit/auction/keeper"
	auctiontypes "github.com/fatal-fruit/auction/types"
	"testing"
)

var AuctionModule = func() configurator.ModuleOption {
	return func(config *configurator.Config) {
		config.ModuleConfigs[auctiontypes.ModuleName] = &cosmosapp.ModuleConfig{
			Name:   auctiontypes.ModuleName,
			Config: appconfig.WrapAny(&auctionmodule.Module{}),
		}
	}
}

func TestIntegration(t *testing.T) {
	t.Parallel()
	logger := log.NewTestLogger(t)
	appConfig := depinject.Configs(
		configurator.NewAppConfig(
			configurator.AuthModule(),
			configurator.BankModule(),
			configurator.StakingModule(),
			configurator.TxModule(),
			configurator.ConsensusModule(),
			configurator.GenutilModule(),
			configurator.MintModule(),
			AuctionModule(),
			configurator.WithCustomInitGenesisOrder(
				"auth",
				"bank",
				"staking",
				"mint",
				"genutil",
				"consensus",
				auctiontypes.ModuleName,
			),
		),
		depinject.Supply(logger),
	)
	var kp keeper.Keeper
	app, err := simtestutil.Setup(appConfig, &kp)
	require.NoError(t, err)
	require.NotNil(t, app)
}
