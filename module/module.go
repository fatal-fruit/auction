package module

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/fatal-fruit/auction/keeper"
)

const ConsensusVersion = 1

type AppModule struct {
	cdc    codec.Codec
	keeper keeper.Keeper
}

func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
	return AppModule{
		cdc:    cdc,
		keeper: keeper,
	}
}
