package integration_test

// ExampleModule is a configurator.ModuleOption that add the auction module to the app config.
//var ExampleModule = func() configurator.ModuleOption {
//	return func(config *configurator.Config) {
//		config.ModuleConfigs[auction.ModuleName] = &appv1alpha1.ModuleConfig{
//			Name:   auction.ModuleName,
//			Config: appconfig.WrapAny(&auctionmodulev1.Module{}),
//		}
//	}
//}

//func TestIntegration(t *testing.T) {
//	t.Parallel()
//
//	logger := log.NewTestLogger(t)
//	appConfig := depinject.Configs(
//		configurator.NewAppConfig(
//			configurator.AuthModule(),
//			configurator.BankModule(),
//			configurator.StakingModule(),
//			configurator.TxModule(),
//			configurator.ConsensusModule(),
//			configurator.GenutilModule(),
//			configurator.MintModule(),
//			ExampleModule(),
//			configurator.WithCustomInitGenesisOrder(
//				"auth",
//				"bank",
//				"staking",
//				"mint",
//				"genutil",
//				"consensus",
//				auction.ModuleName,
//			),
//		),
//		depinject.Supply(logger))
//
//	var keeper keeper.Keeper
//	app, err := simtestutil.Setup(appConfig, &keeper)
//	require.NoError(t, err)
//	require.NotNil(t, app) // use the app or the keeper for running integration tests
//}
