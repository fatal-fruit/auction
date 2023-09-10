package module

import "cosmossdk.io/core/appmodule"

var _ appmodule.AppModule = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

func init() {

}

type ModuleInputs struct {
}

type ModuleOutputs struct {
}
