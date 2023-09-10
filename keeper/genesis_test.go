package keeper_test

import (
	auctiontypes "github.com/fatal-fruit/auction/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitGenesis(t *testing.T) {
	fixture := initFixture(t)

	data := &auctiontypes.GenesisState{}
	err := fixture.k.InitGenesis(fixture.ctx, data)
	require.NoError(t, err)
}

func TestExportGenesis(t *testing.T) {
	fixture := initFixture(t)

	_, err := fixture.k.ExportGenesis(fixture.ctx)
	require.NoError(t, err)
}
